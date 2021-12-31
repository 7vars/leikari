package cluster

import (
	"github.com/7vars/leikari"
	"github.com/hashicorp/memberlist"
)

func NodeName(name string) leikari.Option {
	return leikari.Option{
		Name: "name",
		Value: name,
	}
}

func BindAddress(addr string) leikari.Option {
	return leikari.Option{
		Name: "bindAddr",
		Value: addr,
	}
}

func BindPort(port int) leikari.Option {
	return leikari.Option{
		Name: "bindPort",
		Value: port,
	}
}

func Seeds(seed ...string) leikari.Option {
	return leikari.Option{
		Name: "seeds",
		Value: seed,
	}
}
// TODO define options

type cluster struct {
	system leikari.System
	settings ClusterSettings
	connection *memberlist.Memberlist
	broadcasts *memberlist.TransmitLimitedQueue
}

func newCluster(system leikari.System, opts ...leikari.Option) *cluster {
	settings := newClusterSettings(system.Settings().GetSub("cluster", opts...))
	return &cluster{
		system: system,
		settings: *settings,
	}
}

func (c *cluster) Receive(ctx leikari.ActorContext, msg leikari.Message) {
	
}

func (c *cluster) PreStart(ctx leikari.ActorContext) error {
	cfg := c.settings.Config
	
	cfg.LogOutput = newLogWrapper(ctx.Log())
	cfg.Events = c
	cfg.Delegate = c

	conn, err := memberlist.Create(cfg)
	if err != nil {
		return err
	}
	c.connection = conn

	c.broadcasts = &memberlist.TransmitLimitedQueue{
		NumNodes: func() int {
			return conn.NumMembers()
		},
		RetransmitMult: cfg.RetransmitMult,
	}

	n := conn.LocalNode()
	ctx.Log().Infof("cluster-node %s started at %s", n.Name, n.Address())

	return nil
}

func (c *cluster) PostStop(ctx leikari.ActorContext) error {
	c.connection.Leave(c.settings.CloseTimeout())
	return c.connection.Shutdown()
}

func (c *cluster) NotifyJoin(node *memberlist.Node) {
	c.system.Publish(NodeJoined{newNode(node)})
}

func (c *cluster) NotifyLeave(node *memberlist.Node) {
	c.system.Publish(NodeLeft{newNode(node)})
} 

func (c *cluster) NotifyUpdate(node *memberlist.Node) {
	c.system.Publish(NodeUpdated{newNode(node)})
}

func (c *cluster) NodeMeta(limit int) []byte {
	return c.settings.id[:]
}

func (c *cluster) NotifyMsg(b []byte) {
	// TODO gob
}

func (c *cluster) GetBroadcasts(overhead, limit int) [][]byte {
	if c.broadcasts != nil {
		return c.broadcasts.GetBroadcasts(overhead, limit)
	}
	return [][]byte{}
}

func (c *cluster) LocalState(join bool) []byte {
	// TODO implement
	return []byte{}
}

func (c *cluster) MergeRemoteState(buf []byte, join bool) {
	// TODO implement
}

func Cluster(system leikari.System, opts ...leikari.Option) (leikari.ActorHandler, error) {
	return system.ExecuteService(newCluster(system, opts...), "cluster", opts...)
}