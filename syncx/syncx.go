// Package syncx provides concurrent data structures and utilities.
package syncx

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// SafeMap is a generic thread-safe map.
type SafeMap[K comparable, V any] struct {
	mu sync.RWMutex
	m  map[K]V
}

// NewSafeMap creates a new SafeMap.
func NewSafeMap[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{m: make(map[K]V)}
}

// Set sets the value for key.
func (sm *SafeMap[K, V]) Set(key K, val V) {
	sm.mu.Lock()
	sm.m[key] = val
	sm.mu.Unlock()
}

// Get returns the value for key and whether it was found.
func (sm *SafeMap[K, V]) Get(key K) (V, bool) {
	sm.mu.RLock()
	v, ok := sm.m[key]
	sm.mu.RUnlock()
	return v, ok
}

// GetOrSet returns the existing value for key, or sets and returns def if absent.
func (sm *SafeMap[K, V]) GetOrSet(key K, def V) V {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if v, ok := sm.m[key]; ok {
		return v
	}
	sm.m[key] = def
	return def
}

// Delete removes the key from the map.
func (sm *SafeMap[K, V]) Delete(key K) {
	sm.mu.Lock()
	delete(sm.m, key)
	sm.mu.Unlock()
}

// Has reports whether the key exists.
func (sm *SafeMap[K, V]) Has(key K) bool {
	sm.mu.RLock()
	_, ok := sm.m[key]
	sm.mu.RUnlock()
	return ok
}

// Len returns the number of entries.
func (sm *SafeMap[K, V]) Len() int {
	sm.mu.RLock()
	n := len(sm.m)
	sm.mu.RUnlock()
	return n
}

// Keys returns a snapshot of all keys.
func (sm *SafeMap[K, V]) Keys() []K {
	sm.mu.RLock()
	keys := make([]K, 0, len(sm.m))
	for k := range sm.m {
		keys = append(keys, k)
	}
	sm.mu.RUnlock()
	return keys
}

// Values returns a snapshot of all values.
func (sm *SafeMap[K, V]) Values() []V {
	sm.mu.RLock()
	vals := make([]V, 0, len(sm.m))
	for _, v := range sm.m {
		vals = append(vals, v)
	}
	sm.mu.RUnlock()
	return vals
}

// Range calls fn for each key-value pair. Iteration stops if fn returns false.
func (sm *SafeMap[K, V]) Range(fn func(K, V) bool) {
	sm.mu.RLock()
	snapshot := make(map[K]V, len(sm.m))
	for k, v := range sm.m {
		snapshot[k] = v
	}
	sm.mu.RUnlock()
	for k, v := range snapshot {
		if !fn(k, v) {
			return
		}
	}
}

// Clear removes all entries.
func (sm *SafeMap[K, V]) Clear() {
	sm.mu.Lock()
	sm.m = make(map[K]V)
	sm.mu.Unlock()
}

// ────────────────────────────────────────────────
// Once – generic lazy initializer
// ────────────────────────────────────────────────

// Once initializes a value exactly once.
type Once[T any] struct {
	once  sync.Once
	value T
	err   error
}

// Do calls fn exactly once and caches the result.
// Subsequent calls return the cached value and error without calling fn again.
func (o *Once[T]) Do(fn func() (T, error)) (T, error) {
	o.once.Do(func() {
		o.value, o.err = fn()
	})
	return o.value, o.err
}

// ────────────────────────────────────────────────
// WaitGroup with context support
// ────────────────────────────────────────────────

// WaitGroupCtx wraps sync.WaitGroup with context-aware Wait.
type WaitGroupCtx struct {
	wg sync.WaitGroup
}

// Add adds delta to the counter.
func (wg *WaitGroupCtx) Add(delta int) { wg.wg.Add(delta) }

// Done decrements the counter.
func (wg *WaitGroupCtx) Done() { wg.wg.Done() }

// Wait blocks until the counter is zero or ctx is done, whichever comes first.
// Returns ctx.Err() if the context was canceled before all goroutines finished.
func (wg *WaitGroupCtx) Wait(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		wg.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// ────────────────────────────────────────────────
// Notifier — debounce/coalesce pattern
// ────────────────────────────────────────────────

// Clock abstracts time for testing.
type Clock interface {
	Now() time.Time
	After(d time.Duration) <-chan time.Time
}

type realClock struct{}

func (realClock) Now() time.Time                         { return time.Now() }
func (realClock) After(d time.Duration) <-chan time.Time { return time.After(d) }

// Notifier coalesces rapid repeated events for the same key:
// if a new event arrives while the handler is running, the handler is
// canceled and restarted with the latest event.
type Notifier[T any] struct {
	handler func(ctx context.Context, event T) error
	keyer   func(event T) string
	ttl     time.Duration
	clock   Clock

	keyStates sync.Map
}

// NotifierOption configures a Notifier.
type NotifierOption[T any] func(*Notifier[T])

// WithClock sets a custom clock (useful for testing).
func WithClock[T any](clock Clock) NotifierOption[T] {
	return func(n *Notifier[T]) { n.clock = clock }
}

// NewNotifier creates a Notifier.
// handler is called with each event. keyer extracts a deduplication key.
// ttl controls how long idle per-key state is retained before cleanup.
func NewNotifier[T any](
	handler func(ctx context.Context, event T) error,
	keyer func(event T) string,
	ttl time.Duration,
	opts ...NotifierOption[T],
) *Notifier[T] {
	n := &Notifier[T]{
		handler: handler,
		keyer:   keyer,
		ttl:     ttl,
		clock:   realClock{},
	}
	for _, opt := range opts {
		opt(n)
	}
	go n.cleanupLoop()
	return n
}

type keyExecutor[T any] struct {
	queue     chan *pendingEvent[T]
	cancel    chan struct{}
	slot      chan struct{}
	expiresAt atomic.Value
}

type pendingEvent[T any] struct {
	event T
	done  chan error
}

// Notify sends event to the handler, coalescing concurrent events for the same key.
func (n *Notifier[T]) Notify(ctx context.Context, event T) error {
	key := n.keyer(event)
	value, loaded := n.keyStates.LoadOrStore(key, &keyExecutor[T]{
		queue:  make(chan *pendingEvent[T], 1),
		cancel: make(chan struct{}, 1),
		slot:   make(chan struct{}, 1),
	})

	exec := value.(*keyExecutor[T])
	exec.expiresAt.Store(n.clock.Now().Add(n.ttl))
	_ = loaded

	resultChan := make(chan error, 1)
	pe := &pendingEvent[T]{event: event, done: resultChan}

	select {
	case exec.queue <- pe:
		select {
		case exec.slot <- struct{}{}:
			go n.runEventLoop(ctx, exec)
		default:
			select {
			case exec.cancel <- struct{}{}:
			default:
			}
		}
		return <-resultChan
	default:
		return nil
	}
}

func (n *Notifier[T]) runEventLoop(ctx context.Context, exec *keyExecutor[T]) {
	defer func() { <-exec.slot }()
	pe := <-exec.queue
	n.executeEvent(ctx, exec, pe)
	for {
		select {
		case pe := <-exec.queue:
			n.executeEvent(ctx, exec, pe)
		default:
			return
		}
	}
}

func (n *Notifier[T]) executeEvent(ctx context.Context, exec *keyExecutor[T], pe *pendingEvent[T]) {
	handlerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	done := make(chan error, 1)
	go func() { done <- n.handler(handlerCtx, pe.event) }()

	var err error
	select {
	case err = <-done:
	case <-exec.cancel:
		cancel()
		<-done
		err = nil
		for {
			select {
			case <-exec.cancel:
			default:
				pe.done <- err
				close(pe.done)
				return
			}
		}
	}
	pe.done <- err
	close(pe.done)
}

func (n *Notifier[T]) cleanupLoop() {
	for {
		<-n.clock.After(n.ttl)
		now := n.clock.Now()
		n.keyStates.Range(func(key, value any) bool {
			exec := value.(*keyExecutor[T])
			if len(exec.queue) > 0 {
				return true
			}
			if expVal := exec.expiresAt.Load(); expVal != nil {
				if now.After(expVal.(time.Time)) {
					n.keyStates.Delete(key)
				}
			}
			return true
		})
	}
}
