// Package netx provides network utility functions.
package netx

import (
	"fmt"
	"net"
	"time"
)

// LocalIP returns the first non-loopback IPv4 address of the machine.
// Returns "" if none is found.
func LocalIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			if ip4 := ip.To4(); ip4 != nil {
				return ip4.String()
			}
		}
	}
	return ""
}

// LocalIPs returns all non-loopback IPv4 addresses of the machine.
func LocalIPs() []string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	var ips []string
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			if ip4 := ip.To4(); ip4 != nil {
				ips = append(ips, ip4.String())
			}
		}
	}
	return ips
}

// IsIPv4 reports whether s is a valid IPv4 address.
func IsIPv4(s string) bool {
	ip := net.ParseIP(s)
	return ip != nil && ip.To4() != nil
}

// IsIPv6 reports whether s is a valid IPv6 address.
func IsIPv6(s string) bool {
	ip := net.ParseIP(s)
	return ip != nil && ip.To4() == nil
}

// IsValidIP reports whether s is a valid IP address (IPv4 or IPv6).
func IsValidIP(s string) bool {
	return net.ParseIP(s) != nil
}

// IsValidCIDR reports whether s is a valid CIDR notation.
func IsValidCIDR(s string) bool {
	_, _, err := net.ParseCIDR(s)
	return err == nil
}

// IsIPInCIDR reports whether ip is within the CIDR range.
func IsIPInCIDR(ipStr, cidr string) (bool, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false, fmt.Errorf("invalid IP: %s", ipStr)
	}
	_, network, err := net.ParseCIDR(cidr)
	if err != nil {
		return false, err
	}
	return network.Contains(ip), nil
}

// GetFreePort asks the OS for a free TCP port.
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close() //nolint:errcheck
	return l.Addr().(*net.TCPAddr).Port, nil
}

// IsPortOpen reports whether a TCP connection can be made to host:port within timeout.
func IsPortOpen(host string, port int, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), timeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

// LookupHostname resolves hostname to a list of IP addresses.
func LookupHostname(hostname string) ([]string, error) {
	addrs, err := net.LookupHost(hostname)
	if err != nil {
		return nil, err
	}
	return addrs, nil
}

// ParseHostPort parses "host:port" into separate host and port.
func ParseHostPort(hostport string) (host string, port string, err error) {
	return net.SplitHostPort(hostport)
}
