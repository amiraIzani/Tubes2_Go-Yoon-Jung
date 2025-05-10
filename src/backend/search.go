package main

import (
	"context"
	"sync/atomic"
	"time"
)

func bfsSearch(target string) (*Tree, int64, int) {
	start := time.Now()
	root := &Tree{Name: target}
	var visited int32
	atomic.AddInt32(&visited, 1)
	if node, ok := graph[target]; ok && len(node.Parents) > 0 {
		parallelExpand(root, node.Parents, &visited)
	}
	elapsed := time.Since(start).Microseconds()
	return root, elapsed, int(visited)
}

func dfsSearch(target string) (*Tree, int64, int) {
	start := time.Now()
	root := &Tree{Name: target}
	visitedNodes := make(map[string]struct{})
	var visitedCount int32

	var dfs func(n *Tree)
	dfs = func(n *Tree) {
		if _, seen := visitedNodes[n.Name]; seen {
			return
		}
		visitedNodes[n.Name] = struct{}{}

		atomic.AddInt32(&visitedCount, 1)
		if node, ok := graph[n.Name]; ok && len(node.Parents) > 0 {
			parallelExpand(n, node.Parents, &visitedCount)
			for _, c := range n.Children {
				dfs(c)
			}
		}
	}

	dfs(root)
	elapsed := time.Since(start).Microseconds()
	return root, elapsed, int(visitedCount)
}

func multiSearch(target, method string, limit int) ([]*Tree, int64, int) {
	start := time.Now()
	var visited int32
	results := make(chan *Tree, len(graph[target].Parents))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	node := graph[target]
	for _, p := range node.Parents {
		go func(p [2]string) {
			select {
			case <-ctx.Done():
				return
			default:
			}
			var left, right *Tree
			var v1, v2 int
			if method == "dfs" {
				left, _, v1 = dfsSearch(p[0])
				right, _, v2 = dfsSearch(p[1])
			} else {
				left, _, v1 = bfsSearch(p[0])
				right, _, v2 = bfsSearch(p[1])
			}
			atomic.AddInt32(&visited, int32(v1+v2+1))
			results <- &Tree{Name: target, Children: []*Tree{left, right}}
		}(p)
	}

	var out []*Tree
	timeout := time.After(5 * time.Second)
	for {
		select {
		case t := <-results:
			out = append(out, t)
			if len(out) >= limit {
				cancel()
				elapsed := time.Since(start).Microseconds()
				return out, elapsed, int(visited)
			}
		case <-timeout:
			elapsed := time.Since(start).Microseconds()
			return out, elapsed, int(visited)
		}
	}
}
