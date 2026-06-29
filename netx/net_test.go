package netx_test

import (
	"testing"

	"github.com/tmoeish/gokit/netx"
)

func TestIsIPv4(t *testing.T) {
	if !netx.IsIPv4("192.168.1.1") {
		t.Fatal("IsIPv4 failed")
	}
	if netx.IsIPv4("2001:db8::1") {
		t.Fatal("IPv6 should not be IsIPv4")
	}
	if netx.IsIPv4("not-an-ip") {
		t.Fatal("invalid should not be IsIPv4")
	}
}

func TestIsIPv6(t *testing.T) {
	if !netx.IsIPv6("2001:db8::1") {
		t.Fatal("IsIPv6 failed")
	}
	if netx.IsIPv6("192.168.1.1") {
		t.Fatal("IPv4 should not be IsIPv6")
	}
}

func TestIsValidIP(t *testing.T) {
	if !netx.IsValidIP("10.0.0.1") {
		t.Fatal("IsValidIP")
	}
	if netx.IsValidIP("999.999.999.999") {
		t.Fatal("invalid IP should not be valid")
	}
}

func TestIsIPInCIDR(t *testing.T) {
	ok, err := netx.IsIPInCIDR("192.168.1.5", "192.168.1.0/24")
	if err != nil || !ok {
		t.Fatalf("IsIPInCIDR: %v %v", ok, err)
	}
	ok, err = netx.IsIPInCIDR("10.0.0.1", "192.168.1.0/24")
	if err != nil || ok {
		t.Fatalf("IsIPInCIDR outside range: %v %v", ok, err)
	}
}

func TestGetFreePort(t *testing.T) {
	port, err := netx.GetFreePort()
	if err != nil || port <= 0 {
		t.Fatalf("GetFreePort: %d %v", port, err)
	}
}

func TestLocalIP(t *testing.T) {
	// Just ensure it doesn't panic; result depends on the machine.
	_ = netx.LocalIP()
	_ = netx.LocalIPs()
}
