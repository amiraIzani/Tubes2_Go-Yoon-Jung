package main

import (
	"context"
	"runtime"
	"sync/atomic"
	"time"
)

func bfsSearch(target string) (*Tree, int64, int) {
	start := time.Now()
	root := &Tree{Name: target}

	type item struct {
		node    *Tree
		visited map[string]struct{}
	}
	queue := []item{{node: root, visited: map[string]struct{}{target: {}}}}
	count := 0

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		count++
		name := curr.node.Name

		if _, ok := baseElements[name]; ok {
			return root, time.Since(start).Microseconds(), count
		}

		for _, pair := range graph[name].Parents {
			a, b := pair[0], pair[1]
			if _, seen := curr.visited[a]; seen {
				continue
			}
			if _, seen := curr.visited[b]; seen {
				continue
			}

			left := &Tree{Name: a}
			right := &Tree{Name: b}
			curr.node.Children = append(curr.node.Children, left, right)

			newVisited := make(map[string]struct{}, len(curr.visited)+2)
			for k := range curr.visited {
				newVisited[k] = struct{}{}
			}
			newVisited[a], newVisited[b] = struct{}{}, struct{}{}

			queue = append(queue,
				item{node: left, visited: newVisited},
				item{node: right, visited: newVisited},
			)
			break
		}
	}
	return root, time.Since(start).Microseconds(), count
}

func dfsSearch(target string) (*Tree, int64, int) {
	start := time.Now()
	root := &Tree{Name: target}
	count := 0

	var dfs func(n *Tree, visited map[string]struct{})
	dfs = func(n *Tree, visited map[string]struct{}) {
		count++
		if _, ok := baseElements[n.Name]; ok {
			return
		}
		for _, pair := range graph[n.Name].Parents {
			if _, seen := visited[pair[0]]; seen {
				continue
			}
			if _, seen := visited[pair[1]]; seen {
				continue
			}
			left := &Tree{Name: pair[0]}
			right := &Tree{Name: pair[1]}
			n.Children = append(n.Children, left, right)
			newVisited := make(map[string]struct{}, len(visited)+1)
			for k := range visited {
				newVisited[k] = struct{}{}
			}
			newVisited[n.Name] = struct{}{}

			dfs(left, newVisited)
			dfs(right, newVisited)

			break
		}
	}

	dfs(root, map[string]struct{}{target: {}})
	return root, time.Since(start).Microseconds(), count
}

func multiSearch(target, method string, limit int) ([]*Tree, int64, int) {
	start := time.Now()
	var visited int32
	node := graph[target]
	jobs := make(chan [2]string, len(node.Parents))
	results := make(chan *Tree, len(node.Parents))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	nWorkers := runtime.GOMAXPROCS(0) - (1 + (runtime.GOMAXPROCS(0)-3)/(runtime.GOMAXPROCS(0)/2))

	worker := func() {
		for {
			select {
			case <-ctx.Done():
				return
			case p, ok := <-jobs:
				if !ok {
					return
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
			}
		}
	}

	for i := 0; i < nWorkers; i++ {
		go worker()
	}

	go func() {
		for _, p := range node.Parents {
			jobs <- p
		}
		close(jobs)
	}()

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
