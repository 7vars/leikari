package cluster

import (
	"github.com/7vars/leikari"
)

type clusterNode struct {
	
}

func (node *clusterNode) PreStart(ctx leikari.ActorContext) error {

	return nil
}

func (node *clusterNode) PostStop(ctx leikari.ActorContext) error {

	return nil
}