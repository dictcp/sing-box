package rule

import (
	"context"
	"net/netip"
	"testing"

	"github.com/sagernet/sing-box/adapter"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	slogger "github.com/sagernet/sing/common/logger"

	"github.com/stretchr/testify/require"
)

func TestNewRuleAction_ResolveRouteOnly(t *testing.T) {
	t.Parallel()
	options := option.RuleAction{
		Action: "resolve",
		ResolveOptions: option.RouteActionResolve{
			Server:    "local",
			Strategy:  option.DomainStrategy(C.DomainStrategyPreferIPv4),
			RouteOnly: true,
		},
	}
	action, err := NewRuleAction(context.Background(), slogger.NOP(), options)
	require.NoError(t, err)
	require.NotNil(t, action)

	resolveAction, ok := action.(*RuleActionResolve)
	require.True(t, ok)
	require.Equal(t, "local", resolveAction.Server)
	require.Equal(t, C.DomainStrategyPreferIPv4, resolveAction.Strategy)
	require.True(t, resolveAction.RouteOnly)
}

func TestRuleActionResolve_String(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		action   *RuleActionResolve
		expected string
	}{
		{
			name:     "empty",
			action:   &RuleActionResolve{},
			expected: "resolve",
		},
		{
			name: "server only",
			action: &RuleActionResolve{
				Server: "google",
			},
			expected: "resolve(google)",
		},
		{
			name: "route_only",
			action: &RuleActionResolve{
				RouteOnly: true,
			},
			expected: "resolve(route_only)",
		},
		{
			name: "server and route_only",
			action: &RuleActionResolve{
				Server:    "google",
				RouteOnly: true,
			},
			expected: "resolve(google,route_only)",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, testCase.expected, testCase.action.String())
		})
	}
}

func TestRouteOnlyAddresses_MatchIPRules(t *testing.T) {
	t.Parallel()
	cidr, err := NewIPCIDRItem(false, []string{"8.8.8.0/24"})
	require.NoError(t, err)
	privateItem := NewIPIsPrivateItem(false)
	acceptAnyItem := NewIPAcceptAnyItem()

	testCases := []struct {
		name     string
		metadata adapter.InboundContext
		matched  map[RuleItem]bool
	}{
		{
			name: "ip_cidr matches RouteDestinationAddresses",
			metadata: func() adapter.InboundContext {
				return testRouteOnlyMetadata("lookup.example", netip.MustParseAddr("8.8.8.8"))
			}(),
			matched: map[RuleItem]bool{
				cidr:          true,
				privateItem:   false,
				acceptAnyItem: true,
			},
		},
		{
			name: "ip_is_private matches RouteDestinationAddresses",
			metadata: func() adapter.InboundContext {
				return testRouteOnlyMetadata("lookup.example", netip.MustParseAddr("10.0.0.1"))
			}(),
			matched: map[RuleItem]bool{
				cidr:          false,
				privateItem:   true,
				acceptAnyItem: true,
			},
		},
		{
			name: "no resolved addresses, ip_accept_any false",
			metadata: func() adapter.InboundContext {
				return testRouteOnlyMetadata("lookup.example")
			}(),
			matched: map[RuleItem]bool{
				cidr:          false,
				privateItem:   false,
				acceptAnyItem: false,
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			for item, expected := range testCase.matched {
				require.Equal(t, expected, item.Match(&testCase.metadata), "item=%s", item.String())
			}
		})
	}
}

func testRouteOnlyMetadata(domain string, addresses ...netip.Addr) adapter.InboundContext {
	m := testMetadata(domain)
	m.RouteDestinationAddresses = addresses
	return m
}
