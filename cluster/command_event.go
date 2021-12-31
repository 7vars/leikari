package cluster

type NodeJoined struct {
	Node
}

type NodeLeft struct {
	Node
}

type NodeUpdated struct {
	Node
}