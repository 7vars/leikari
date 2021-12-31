package cluster

import (
	"fmt"
	"os"
	"time"

	"github.com/7vars/leikari"
	"github.com/google/uuid"
	"github.com/hashicorp/memberlist"
)

type ClusterSettings struct {
	*memberlist.Config
	id uuid.UUID
	Settings leikari.Settings
}

func newClusterSettings(s leikari.Settings) *ClusterSettings {
	result := &ClusterSettings{
		memberlist.DefaultLANConfig(),
		uuid.New(),
		s,
	}

	result.BindAddr = s.GetDefaultString("bindAddr", "0.0.0.0")
	result.BindPort = s.GetDefaultInt("bindPort", 7946)
	result.AdvertiseAddr = s.GetDefaultString("advertiseAddr", "")
	result.AdvertisePort = s.GetDefaultInt("advertisePort", 7946)
	result.TCPTimeout = s.GetDefaultDuration("tcpTimeout", 10 * time.Second)
	result.IndirectChecks = s.GetDefaultInt("indirectChecks", 3)
	result.RetransmitMult = s.GetDefaultInt("retransmitMult", 4)
	result.SuspicionMaxTimeoutMult = s.GetDefaultInt("suspicionMaxTimeoutMult", 6)
	result.PushPullInterval = s.GetDefaultDuration("pushPullInterval", 30 * time.Second)
	result.ProbeTimeout = s.GetDefaultDuration("probeTimeout", 500 * time.Millisecond)
	result.ProbeInterval = s.GetDefaultDuration("probeInterval", 1 * time.Second)
	result.DisableTcpPings = s.GetBool("disableTcpPings")
	result.AwarenessMaxMultiplier = s.GetDefaultInt("awarenessMaxMult", 8)

	result.GossipNodes = s.GetDefaultInt("gossipNodes", 3)
	result.GossipInterval = s.GetDefaultDuration("gossipInterval", 200 * time.Millisecond)
	result.GossipToTheDeadTime = s.GetDefaultDuration("gossipToTheDeadTime", 30 * time.Second)
	result.GossipVerifyIncoming = s.GetDefaultBool("gossipVerifyIncoming", true)
	result.GossipVerifyOutgoing = s.GetDefaultBool("gossipVerifyOutgoing", true)

	result.EnableCompression = s.GetDefaultBool("enableCompression", true)

	result.DNSConfigPath = s.GetDefaultString("dnsConfigPath", "/etc/resolve.conf")

	result.HandoffQueueDepth = s.GetDefaultInt("handoffQueueDepth", 1024)
	result.UDPBufferSize = s.GetDefaultInt("udpBufferSize", 1400)


	name, _ := os.Hostname()
	result.Name = s.GetDefaultString("name", fmt.Sprintf("%s-%s", name, result.id.String()))

	return result
}

func (cs *ClusterSettings) Seeds() []string {
	return cs.Settings.GetStringSlice("seeds")
}

func (cs *ClusterSettings) CloseTimeout() time.Duration {
	return cs.Settings.GetDefaultDuration("closeTimeout", 1 * time.Second)
}