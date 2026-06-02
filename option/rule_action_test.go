package option

import (
	"context"
	"testing"

	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing/common/json"

	"github.com/stretchr/testify/require"
)

func TestDNSRuleActionRespondUnmarshalJSON(t *testing.T) {
	t.Parallel()

	var action DNSRuleAction
	err := json.UnmarshalContext(context.Background(), []byte(`{"action":"respond"}`), &action)
	require.NoError(t, err)
	require.Equal(t, C.RuleActionTypeRespond, action.Action)
	require.Equal(t, DNSRouteActionOptions{}, action.RouteOptions)
}

func TestDNSRuleActionRespondRejectsUnknownFields(t *testing.T) {
	t.Parallel()

	var action DNSRuleAction
	err := json.UnmarshalContext(context.Background(), []byte(`{"action":"respond","disable_cache":true}`), &action)
	require.ErrorContains(t, err, "unknown field")
}

func TestResolveActionRouteOnlyUnmarshalJSON(t *testing.T) {
	t.Parallel()

	var action RuleAction
	err := json.Unmarshal([]byte(`{"action":"resolve","route_only":true}`), &action)
	require.NoError(t, err)
	require.Equal(t, C.RuleActionTypeResolve, action.Action)
	require.True(t, action.ResolveOptions.RouteOnly)
}

func TestResolveActionRouteOnlyMarshalJSON(t *testing.T) {
	t.Parallel()

	action := RuleAction{
		Action: C.RuleActionTypeResolve,
		ResolveOptions: RouteActionResolve{
			RouteOnly: true,
		},
	}
	data, err := json.Marshal(action)
	require.NoError(t, err)
	require.Contains(t, string(data), `"route_only":true`)
}

func TestResolveActionRouteOnlyFalseUnmarshalJSON(t *testing.T) {
	t.Parallel()

	var action RuleAction
	err := json.Unmarshal([]byte(`{"action":"resolve","route_only":false}`), &action)
	require.NoError(t, err)
	require.Equal(t, C.RuleActionTypeResolve, action.Action)
	require.False(t, action.ResolveOptions.RouteOnly)
}

func TestResolveActionRouteOnlyDefaultUnmarshalJSON(t *testing.T) {
	t.Parallel()

	var action RuleAction
	err := json.Unmarshal([]byte(`{"action":"resolve"}`), &action)
	require.NoError(t, err)
	require.Equal(t, C.RuleActionTypeResolve, action.Action)
	require.False(t, action.ResolveOptions.RouteOnly)
}
