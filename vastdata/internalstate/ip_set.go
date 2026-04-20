// Copyright (c) HashiCorp, Inc.

package internalstate

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// expandIPSet parses a list of IP addresses and CIDR notations and returns
// the full set of individual host IP strings they collectively cover.
// Both IPv4 and IPv6 are supported.
// Returns an error if any entry cannot be parsed.
func expandIPSet(entries []string) (map[string]struct{}, error) {
	result := make(map[string]struct{})
	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		if strings.Contains(entry, "/") {
			_, ipNet, err := net.ParseCIDR(entry)
			if err != nil {
				return nil, fmt.Errorf("invalid CIDR %q: %w", entry, err)
			}
			for ip := cloneIP(ipNet.IP); ipNet.Contains(ip); incrementIP(ip) {
				result[ip.String()] = struct{}{}
			}
		} else {
			ip := net.ParseIP(entry)
			if ip == nil {
				return nil, fmt.Errorf("invalid IP address %q", entry)
			}
			// Normalise to canonical string form so IPv4-mapped IPv6 addresses
			// don't cause false mismatches.
			if v4 := ip.To4(); v4 != nil {
				result[v4.String()] = struct{}{}
			} else {
				result[ip.String()] = struct{}{}
			}
		}
	}
	return result, nil
}

// IPSetsEquivalent returns true when two lists of IPs / CIDRs cover exactly
// the same set of individual host addresses.
// Parsing errors on either side are treated as non-equivalent so the caller
// falls back to a full state update.
func IPSetsEquivalent(a, b []string) bool {
	setA, err := expandIPSet(a)
	if err != nil {
		return false
	}
	setB, err := expandIPSet(b)
	if err != nil {
		return false
	}
	if len(setA) != len(setB) {
		return false
	}
	for ip := range setA {
		if _, ok := setB[ip]; !ok {
			return false
		}
	}
	return true
}

// attrValueToStringSlice extracts a []string from a Terraform list or set of
// strings.  Returns nil if the value is null, unknown, or not a
// list/set-of-string.
func attrValueToStringSlice(v attr.Value) []string {
	if v == nil || v.IsNull() || v.IsUnknown() {
		return nil
	}
	var elems []attr.Value
	switch tv := v.(type) {
	case types.List:
		elems = tv.Elements()
	case types.Set:
		elems = tv.Elements()
	default:
		return nil
	}
	out := make([]string, 0, len(elems))
	for _, e := range elems {
		s, ok := e.(types.String)
		if !ok {
			return nil
		}
		out = append(out, s.ValueString())
	}
	return out
}

// IPSetsEquivalentAttrs is the attr.Value level wrapper used in fillFromRecordInternal.
func IPSetsEquivalentAttrs(existing, incoming attr.Value) bool {
	a := attrValueToStringSlice(existing)
	b := attrValueToStringSlice(incoming)
	if a == nil || b == nil {
		return false
	}
	return IPSetsEquivalent(a, b)
}

// cloneIP returns a copy of the IP so we can safely mutate it while iterating.
func cloneIP(ip net.IP) net.IP {
	dup := make(net.IP, len(ip))
	copy(dup, ip)
	return dup
}

// incrementIP adds 1 to the lowest-order byte of ip (big-endian).
func incrementIP(ip net.IP) {
	if len(ip) == 4 {
		n := binary.BigEndian.Uint32(ip)
		binary.BigEndian.PutUint32(ip, n+1)
	} else {
		// IPv6: treat the 16-byte slice as a 128-bit big-endian integer.
		for i := len(ip) - 1; i >= 0; i-- {
			ip[i]++
			if ip[i] != 0 {
				break
			}
		}
	}
}
