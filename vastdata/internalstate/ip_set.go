// Copyright (c) HashiCorp, Inc.

package internalstate

import (
	"fmt"
	"net/netip"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ipRange is an inclusive [start, end] address range within one IP family.
type ipRange struct{ start, end netip.Addr }

// mergeIPRanges parses a list of IP addresses and CIDR notations and returns
// the sorted, minimal list of disjoint, non-adjacent ranges they cover, so two
// lists cover the same addresses iff their merged ranges are equal.
// Both IPv4 and IPv6 are supported; IPv4-mapped IPv6 addresses are unmapped.
// Returns an error if any entry cannot be parsed.
func mergeIPRanges(entries []string) ([]ipRange, error) {
	ranges := make([]ipRange, 0, len(entries))
	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		var prefix netip.Prefix
		if strings.Contains(entry, "/") {
			p, err := netip.ParsePrefix(entry)
			if err != nil {
				return nil, fmt.Errorf("invalid CIDR %q: %w", entry, err)
			}
			prefix = p
		} else {
			a, err := netip.ParseAddr(entry)
			if err != nil {
				return nil, fmt.Errorf("invalid IP address %q: %w", entry, err)
			}
			prefix = netip.PrefixFrom(a, a.BitLen())
		}
		// A mapped prefix shorter than /96 extends beyond ::ffff:0:0/96, so it stays IPv6.
		if a := prefix.Addr(); a.Is4In6() && prefix.Bits() >= 96 {
			prefix = netip.PrefixFrom(a.Unmap(), prefix.Bits()-96)
		}
		prefix = prefix.Masked()
		ranges = append(ranges, ipRange{prefix.Addr(), lastAddr(prefix)})
	}

	slices.SortFunc(ranges, func(x, y ipRange) int { return x.start.Compare(y.start) })
	merged := ranges[:0]
	for _, r := range ranges {
		if n := len(merged); n > 0 {
			cur := &merged[n-1]
			// Compare orders IPv4 before IPv6, and Next() of a family's last address is
			// the zero Addr, so ranges never merge across families.
			if r.start.Compare(cur.end) <= 0 || r.start == cur.end.Next() {
				if r.end.Compare(cur.end) > 0 {
					cur.end = r.end
				}
				continue
			}
		}
		merged = append(merged, r)
	}
	return merged, nil
}

// lastAddr returns the highest address inside a masked prefix.
func lastAddr(p netip.Prefix) netip.Addr {
	b := p.Addr().AsSlice()
	for i := p.Bits(); i < len(b)*8; i++ {
		b[i/8] |= 0x80 >> (i % 8)
	}
	a, _ := netip.AddrFromSlice(b)
	return a
}

// IPSetsEquivalent returns true when two lists of IPs / CIDRs cover exactly
// the same set of individual host addresses.
// Parsing errors on either side are treated as non-equivalent so the caller
// falls back to a full state update.
func IPSetsEquivalent(a, b []string) bool {
	rangesA, err := mergeIPRanges(a)
	if err != nil {
		return false
	}
	rangesB, err := mergeIPRanges(b)
	if err != nil {
		return false
	}
	return slices.Equal(rangesA, rangesB)
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
