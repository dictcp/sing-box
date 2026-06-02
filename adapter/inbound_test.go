package adapter

import (
	"net"
	"net/netip"
	"testing"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/require"
)

func TestDNSResponseAddressesUnmapsHTTPSIPv4Hints(t *testing.T) {
	t.Parallel()

	ipv4Hint := net.ParseIP("1.1.1.1")
	require.NotNil(t, ipv4Hint)

	response := &dns.Msg{
		MsgHdr: dns.MsgHdr{
			Response: true,
			Rcode:    dns.RcodeSuccess,
		},
		Answer: []dns.RR{
			&dns.HTTPS{
				SVCB: dns.SVCB{
					Hdr: dns.RR_Header{
						Name:   dns.Fqdn("example.com"),
						Rrtype: dns.TypeHTTPS,
						Class:  dns.ClassINET,
						Ttl:    60,
					},
					Priority: 1,
					Target:   ".",
					Value: []dns.SVCBKeyValue{
						&dns.SVCBIPv4Hint{Hint: []net.IP{ipv4Hint}},
					},
				},
			},
		},
	}

	addresses := DNSResponseAddresses(response)
	require.Equal(t, []netip.Addr{netip.MustParseAddr("1.1.1.1")}, addresses)
	require.True(t, addresses[0].Is4())
}

func TestInboundContextClearRouteOnlyDestinationAddresses(t *testing.T) {
	t.Parallel()

	metadata := InboundContext{
		DestinationAddresses:          []netip.Addr{netip.MustParseAddr("192.0.2.1")},
		DestinationAddressesRouteOnly: true,
	}
	metadata.ClearRouteOnlyDestinationAddresses()

	require.Empty(t, metadata.DestinationAddresses)
	require.False(t, metadata.DestinationAddressesRouteOnly)
}

func TestInboundContextClearRouteOnlyDestinationAddressesKeepsDialAddresses(t *testing.T) {
	t.Parallel()

	addresses := []netip.Addr{netip.MustParseAddr("192.0.2.1")}
	metadata := InboundContext{DestinationAddresses: addresses}
	metadata.ClearRouteOnlyDestinationAddresses()

	require.Equal(t, addresses, metadata.DestinationAddresses)
}
