package mathx_test

import (
	"testing"

	"github.com/tmoeish/gokit/mathx"
)

func TestClamp(t *testing.T) {
	if mathx.Clamp(5, 1, 10) != 5 {
		t.Fatal("Clamp within range")
	}
	if mathx.Clamp(-1, 0, 10) != 0 {
		t.Fatal("Clamp below lo")
	}
	if mathx.Clamp(15, 0, 10) != 10 {
		t.Fatal("Clamp above hi")
	}
}

func TestAbs(t *testing.T) {
	if mathx.Abs(-5) != 5 {
		t.Fatal("Abs(-5)")
	}
	if mathx.Abs(5) != 5 {
		t.Fatal("Abs(5)")
	}
}

func TestSum(t *testing.T) {
	if mathx.Sum(1, 2, 3, 4, 5) != 15 {
		t.Fatal("Sum")
	}
}

func TestAverage(t *testing.T) {
	avg := mathx.Average(1, 2, 3, 4, 5)
	if avg != 3.0 {
		t.Fatalf("Average: got %f", avg)
	}
}

func TestGCD(t *testing.T) {
	if mathx.GCD(12, 8) != 4 {
		t.Fatal("GCD(12,8)")
	}
	if mathx.GCD(100, 75) != 25 {
		t.Fatal("GCD(100,75)")
	}
}

func TestLCM(t *testing.T) {
	if mathx.LCM(4, 6) != 12 {
		t.Fatal("LCM(4,6)")
	}
}

func TestPow(t *testing.T) {
	if mathx.Pow(2, 10) != 1024 {
		t.Fatal("Pow(2,10)")
	}
}

func TestIsPrime(t *testing.T) {
	primes := []int{2, 3, 5, 7, 11, 13, 17, 19}
	for _, p := range primes {
		if !mathx.IsPrime(p) {
			t.Fatalf("IsPrime(%d) should be true", p)
		}
	}
	notPrimes := []int{0, 1, 4, 6, 8, 9, 10}
	for _, n := range notPrimes {
		if mathx.IsPrime(n) {
			t.Fatalf("IsPrime(%d) should be false", n)
		}
	}
}

func TestFibonacci(t *testing.T) {
	expected := []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34}
	for i, want := range expected {
		if got := mathx.Fibonacci(i); got != want {
			t.Fatalf("Fibonacci(%d) = %d, want %d", i, got, want)
		}
	}
}

func TestPercent(t *testing.T) {
	if mathx.Percent(1, 4) != 25.0 {
		t.Fatal("Percent(1,4)")
	}
	if mathx.Percent(0, 0) != 0 {
		t.Fatal("Percent divide by zero")
	}
}

func TestRoundTo(t *testing.T) {
	if mathx.RoundTo(3.14159, 2) != 3.14 {
		t.Fatal("RoundTo")
	}
}

func TestInRange(t *testing.T) {
	if !mathx.InRange(5, 1, 10) {
		t.Fatal("InRange")
	}
	if mathx.InRange(0, 1, 10) {
		t.Fatal("InRange out of range")
	}
}
