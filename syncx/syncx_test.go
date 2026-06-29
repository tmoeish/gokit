package syncx_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/tmoeish/gokit/syncx"
)

func TestSafeMap(t *testing.T) {
	sm := syncx.NewSafeMap[string, int]()
	sm.Set("a", 1)
	sm.Set("b", 2)

	v, ok := sm.Get("a")
	if !ok || v != 1 {
		t.Fatal("Get failed")
	}

	if sm.Len() != 2 {
		t.Fatal("Len failed")
	}

	sm.Delete("a")
	if sm.Has("a") {
		t.Fatal("Delete failed")
	}

	// concurrent access
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			sm.Set("key", i)
			sm.Get("key")
		}(i)
	}
	wg.Wait()
}

func TestSafeMapGetOrSet(t *testing.T) {
	sm := syncx.NewSafeMap[string, int]()
	v := sm.GetOrSet("x", 42)
	if v != 42 {
		t.Fatalf("GetOrSet: got %d", v)
	}
	v = sm.GetOrSet("x", 99)
	if v != 42 {
		t.Fatalf("GetOrSet should return existing: got %d", v)
	}
}

func TestOnce(t *testing.T) {
	var o syncx.Once[int]
	calls := 0
	for i := 0; i < 5; i++ {
		v, err := o.Do(func() (int, error) {
			calls++
			return 42, nil
		})
		if err != nil || v != 42 {
			t.Fatalf("Once: v=%d err=%v", v, err)
		}
	}
	if calls != 1 {
		t.Fatalf("Once called %d times, want 1", calls)
	}
}

func TestOnceError(t *testing.T) {
	var o syncx.Once[int]
	sentinel := errors.New("fail")
	v, err := o.Do(func() (int, error) { return 0, sentinel })
	if err != sentinel || v != 0 {
		t.Fatal("Once error case")
	}
	// second call returns same cached error
	v2, err2 := o.Do(func() (int, error) { return 99, nil })
	if err2 != sentinel || v2 != 0 {
		t.Fatal("Once should return cached error")
	}
}

func TestNotifier(t *testing.T) {
	var mu sync.Mutex
	called := 0

	n := syncx.NewNotifier[string](
		func(ctx context.Context, event string) error {
			mu.Lock()
			called++
			mu.Unlock()
			return nil
		},
		func(event string) string { return event },
		time.Second,
	)

	ctx := context.Background()
	if err := n.Notify(ctx, "key1"); err != nil {
		t.Fatalf("Notify: %v", err)
	}

	mu.Lock()
	c := called
	mu.Unlock()
	if c != 1 {
		t.Fatalf("expected 1 call, got %d", c)
	}
}

func TestWaitGroupCtx(t *testing.T) {
	var wg syncx.WaitGroupCtx
	wg.Add(1)
	go func() {
		time.Sleep(10 * time.Millisecond)
		wg.Done()
	}()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := wg.Wait(ctx); err != nil {
		t.Fatalf("WaitGroupCtx: %v", err)
	}
}

func TestWaitGroupCtxCancelled(t *testing.T) {
	var wg syncx.WaitGroupCtx
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := wg.Wait(ctx)
	wg.Done()
	if err == nil {
		t.Fatal("expected context error")
	}
}
