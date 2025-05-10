package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Recipe struct {
	Elements []string `json:"elements"`
}

type Raw struct {
	Name    string   `json:"name"`
	Recipes []Recipe `json:"recipes"`
}

type Node struct {
	Name    string
	Parents [][2]string
}

type Tree struct {
	Name     string  `json:"name"`
	Children []*Tree `json:"children,omitempty"`
}

var graph map[string]Node

func loadRecipes(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read %s: %v", path, err)
	}
	var raws []Raw
	if err := json.Unmarshal(data, &raws); err != nil {
		log.Fatalf("invalid JSON format: %v", err)
	}
	graph = make(map[string]Node, len(raws))
	for _, item := range raws {
		n := Node{Name: item.Name}
		for _, r := range item.Recipes {
			if len(r.Elements) == 2 {
				n.Parents = append(n.Parents, [2]string{r.Elements[0], r.Elements[1]})
			}
		}
		graph[item.Name] = n
	}
}

func parallelExpand(node *Tree, pairs [][2]string, visited *int32) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, p := range pairs {
		wg.Add(1)
		go func(p [2]string) {
			defer wg.Done()
			left := &Tree{Name: p[0]}
			right := &Tree{Name: p[1]}
			mu.Lock()
			node.Children = append(node.Children, left, right)
			mu.Unlock()
			atomic.AddInt32(visited, 2)
		}(p)
	}
	wg.Wait()
}

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

func searchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	target := q.Get("target")
	method := q.Get("method")
	mode := q.Get("mode")
	if target == "" || method == "" || mode == "" {
		http.Error(w, "missing parameter", http.StatusBadRequest)
		return
	}
	limit := 1
	if mode == "multiple" {
		lv, err := strconv.Atoi(q.Get("limit"))
		if err != nil || lv < 1 {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}
		limit = lv
	}

	var recipes []*Tree
	var elapsed int64
	var nodes int

	switch mode {
	case "shortest":
		if method == "dfs" {
			t, e, v := dfsSearch(target)
			recipes = []*Tree{t}
			elapsed, nodes = e, v
		} else {
			t, e, v := bfsSearch(target)
			recipes = []*Tree{t}
			elapsed, nodes = e, v
		}
	case "multiple":
		recipes, elapsed, nodes = multiSearch(target, method, limit)
	default:
		http.Error(w, "unsupported mode", http.StatusBadRequest)
		return
	}

	resp := map[string]interface{}{
		"method":        method,
		"mode":          mode,
		"limit":         limit,
		"time_us":       elapsed,
		"nodes_visited": nodes,
		"recipes":       recipes,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	loadRecipes("recipes.json")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r := mux.NewRouter()
	r.HandleFunc("/search", searchHandler).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))//buat front end sementara
	handler := handlers.CORS(handlers.AllowedOrigins([]string{"*"}))(r)
	fmt.Printf("Server started on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
