package center

type APINode struct {
	Name string `json:"name"`
}

func GetNodes() (nodes []APINode) {
	endCh := make(chan struct{})
	center.taskPool.Submit(func() {
		nodes = make([]APINode, 0, len(center.nodes))
		for _, n := range center.nodes {
			nodes = append(nodes, APINode{Name: n.Name})
		}
		//nodes = append(nodes, APINode{Name: "test"})
		endCh <- struct{}{}
	})
	<-endCh
	return nodes
}
