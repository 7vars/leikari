package cluster

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/memberlist"
)

type Node interface {
	Name() string
	Addr() string
	Port() int
	Address() string
}

type node struct {
	id uuid.UUID
	name string
	addr string
	port int
}

func newNode(n *memberlist.Node) Node {
	var id uuid.UUID
	if len(n.Meta) > 0 {
		id, _ = uuid.FromBytes(n.Meta)
	}
	return &node{
		id: id,
		name: n.Name,
		addr: n.Addr.String(),
		port: int(n.Port),
	}
}

func (n *node) Name() string {
	return n.name
}

func (n *node) Addr() string {
	return n.addr
}

func (n *node) Port() int {
	return n.port
}

func (n *node) Address() string {
	return fmt.Sprintf("%s:%d", n.addr, n.port)
}