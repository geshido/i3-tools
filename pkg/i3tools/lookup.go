package i3tools

import (
	"go.i3wm.org/i3/v4"
	"sync"
)

func FindAll(start *i3.Node, predicate func(*i3.Node) bool) []*i3.Node {
	var result []*i3.Node

	if predicate(start) {
		result = append(result, start)
	}

	ch := make(chan *i3.Node)
	wg := sync.WaitGroup{}
	wg.Add(2)
	finder := func(nodes []*i3.Node) {
		defer wg.Done()
		for _, n := range nodes {
			for _, nn := range FindAll(n, predicate) {
				ch <- nn
			}
		}
	}
	go finder(start.Nodes)
	go finder(start.FloatingNodes)
	go func() { wg.Wait(); close(ch) }()

	for n := range ch {
		result = append(result, n)
	}

	return result
}