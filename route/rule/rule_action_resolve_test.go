package rule

import (
	"context"
	"testing"

	"github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	"github.com/stretchr/testify/require"
)

func TestNewRuleActionResolveRouteOnly(t *testing.T) {
	t.Parallel()

	action, err := NewRuleAction(context.Background(), nil, option.RuleAction{
		Action: constant.RuleActionTypeResolve,
		ResolveOptions: option.RouteActionResolve{
			RouteOnly: true,
		},
	})
	require.NoError(t, err)

	resolveAction, isResolveAction := action.(*RuleActionResolve)
	require.True(t, isResolveAction)
	require.True(t, resolveAction.RouteOnly)
	require.Equal(t, "resolve(route_only)", resolveAction.String())
}
