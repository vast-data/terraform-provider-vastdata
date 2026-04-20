// Copyright (c) HashiCorp, Inc.

package internalstate

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// expandIPSet
// ---------------------------------------------------------------------------

func TestExpandIPSet_SingleIPs(t *testing.T) {
	set, err := expandIPSet([]string{"10.0.0.1", "10.0.0.2", "10.0.0.3"})
	require.NoError(t, err)
	require.Equal(t, 3, len(set))
	require.Contains(t, set, "10.0.0.1")
	require.Contains(t, set, "10.0.0.2")
	require.Contains(t, set, "10.0.0.3")
}

func TestExpandIPSet_CIDR(t *testing.T) {
	// /30 = network + 2 hosts + broadcast = 4 addresses
	set, err := expandIPSet([]string{"10.0.0.0/30"})
	require.NoError(t, err)
	require.Equal(t, 4, len(set))
	require.Contains(t, set, "10.0.0.0")
	require.Contains(t, set, "10.0.0.1")
	require.Contains(t, set, "10.0.0.2")
	require.Contains(t, set, "10.0.0.3")
}

func TestExpandIPSet_Mixed(t *testing.T) {
	set, err := expandIPSet([]string{"10.0.0.1", "10.0.0.4/30"})
	require.NoError(t, err)
	// 1 individual + 4 from /30
	require.Equal(t, 5, len(set))
}

func TestExpandIPSet_InvalidIP(t *testing.T) {
	_, err := expandIPSet([]string{"not-an-ip"})
	require.Error(t, err)
}

func TestExpandIPSet_InvalidCIDR(t *testing.T) {
	_, err := expandIPSet([]string{"10.0.0.0/99"})
	require.Error(t, err)
}

func TestExpandIPSet_Empty(t *testing.T) {
	set, err := expandIPSet([]string{})
	require.NoError(t, err)
	require.Empty(t, set)
}

// ---------------------------------------------------------------------------
// IPSetsEquivalent — the core equivalence logic
// ---------------------------------------------------------------------------

func TestIPSetsEquivalent_IndividualVsAggregated(t *testing.T) {
	// Scenario from the bug ticket:
	// user writes individual IPs; VMS collapses them to CIDRs after ~20 min.
	individual := []string{
		"10.217.224.41",
		"10.217.224.42",
		"10.217.224.43",
		"10.217.224.44",
		"10.217.224.45",
		"10.217.224.46",
		"10.217.224.47",
	}
	// 10.217.224.42/31 = .42,.43   10.217.224.44/30 = .44,.45,.46,.47
	aggregated := []string{
		"10.217.224.41",
		"10.217.224.42/31",
		"10.217.224.44/30",
	}
	require.True(t, IPSetsEquivalent(individual, aggregated))
}

func TestIPSetsEquivalent_UserWritesCIDRsAPIReturnsSameCIDRs(t *testing.T) {
	// User already writes optimised CIDRs; API returns the same CIDRs.
	cidrs := []string{"10.217.224.42/31", "10.217.224.44/30"}
	require.True(t, IPSetsEquivalent(cidrs, cidrs))
}

func TestIPSetsEquivalent_MixedVsAggregated(t *testing.T) {
	// User writes mixed (some CIDRs, some individual IPs).
	mixed := []string{"10.217.224.41", "10.217.224.42/31"}
	// API collapses them all — .41,.42,.43
	aggregated := []string{"10.217.224.41", "10.217.224.42/31"}
	require.True(t, IPSetsEquivalent(mixed, aggregated))
}

func TestIPSetsEquivalent_DifferentSets_ExtraIP(t *testing.T) {
	// Someone adds 10.217.224.68 via the UI — drift should be detected.
	config := []string{"10.217.224.41", "10.217.224.42"}
	apiResp := []string{"10.217.224.41", "10.217.224.42", "10.217.224.68"}
	require.False(t, IPSetsEquivalent(config, apiResp))
}

func TestIPSetsEquivalent_DifferentSets_MissingIP(t *testing.T) {
	// Someone removes 10.217.224.42 via the UI — drift should be detected.
	config := []string{"10.217.224.41", "10.217.224.42"}
	apiResp := []string{"10.217.224.41"}
	require.False(t, IPSetsEquivalent(config, apiResp))
}

func TestIPSetsEquivalent_BothEmpty(t *testing.T) {
	require.True(t, IPSetsEquivalent([]string{}, []string{}))
}

func TestIPSetsEquivalent_InvalidEntryReturnsFalse(t *testing.T) {
	// Parse error → treated as non-equivalent so state gets updated.
	require.False(t, IPSetsEquivalent([]string{"bad"}, []string{"10.0.0.1"}))
}

func TestIPSetsEquivalent_FullTicketScenario(t *testing.T) {
	// Full nfs_no_squash list from the bug ticket vs the aggregated response.
	userList := []string{
		"10.220.131.0/27",
		"10.217.224.41",
		"10.217.224.42", "10.217.224.43",
		"10.217.224.44", "10.217.224.45", "10.217.224.46", "10.217.224.47",
		"10.217.224.48", "10.217.224.49", "10.217.224.50", "10.217.224.51",
		"10.217.224.52", "10.217.224.53", "10.217.224.54", "10.217.224.55",
		"10.217.224.56", "10.217.224.57", "10.217.224.58", "10.217.224.59",
		"10.217.224.60", "10.217.224.61", "10.217.224.62", "10.217.224.63",
		"10.217.224.64", "10.217.224.65", "10.217.224.66", "10.217.224.67",
	}
	// What VMS returns after aggregation
	apiList := []string{
		"10.217.224.41",
		"10.217.224.42/31",
		"10.217.224.44/30",
		"10.217.224.64/30",
		"10.217.224.48/28",
		"10.220.131.0/27",
	}
	require.True(t, IPSetsEquivalent(userList, apiList))
}

// ---------------------------------------------------------------------------
// IPSetsEquivalentAttrs — attr.Value wrapper
// ---------------------------------------------------------------------------

func listAttr(ips ...string) attr.Value {
	elems := make([]attr.Value, len(ips))
	for i, ip := range ips {
		elems[i] = types.StringValue(ip)
	}
	v, _ := types.ListValue(types.StringType, elems)
	return v
}

func TestIPSetsEquivalentAttrs_EqualSets(t *testing.T) {
	a := listAttr("10.0.0.1", "10.0.0.2/31")
	b := listAttr("10.0.0.1", "10.0.0.2", "10.0.0.3")
	require.True(t, IPSetsEquivalentAttrs(a, b))
}

func TestIPSetsEquivalentAttrs_DifferentSets(t *testing.T) {
	a := listAttr("10.0.0.1")
	b := listAttr("10.0.0.2")
	require.False(t, IPSetsEquivalentAttrs(a, b))
}

func TestIPSetsEquivalentAttrs_NullReturnsFalse(t *testing.T) {
	require.False(t, IPSetsEquivalentAttrs(types.ListNull(types.StringType), listAttr("10.0.0.1")))
	require.False(t, IPSetsEquivalentAttrs(listAttr("10.0.0.1"), types.ListNull(types.StringType)))
}
