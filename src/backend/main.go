package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

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
	Children []*Tree `json:"children"`
}

var graph map[string]Node

func loadRecipes(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read recipes.json: %v", err)
	}
	var raws []Raw
	if err := json.Unmarshal(data, &raws); err != nil {
		log.Fatalf("invalid JSON format: %v", err)
	}
	graph = make(map[string]Node)
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

func bfsSearch(target string) (*Tree, int, int) {
	start := time.Now()
	root := &Tree{Name: target}
	queue := []*Tree{root}
	visitedNodes := 0

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		visitedNodes++
		node, ok := graph[current.Name]
		if !ok || len(node.Parents) == 0 {
			continue
		}

		for _, pair := range node.Parents {
			left := &Tree{Name: pair[0]}
			right := &Tree{Name: pair[1]}
			current.Children = append(current.Children, left, right)
			queue = append(queue, left, right)
		}
		break
	}
	timeMs := int(time.Since(start).Milliseconds())
	return root, timeMs, visitedNodes
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	target := q.Get("target")
	if target == "" {
		http.Error(w, "missing target parameter", http.StatusBadRequest)
		return
	}

	tree, elapsed, nodes := bfsSearch(target)
	resp := map[string]interface{}{
		"tree":          tree,
		"time_ms":       elapsed,
		"nodes_visited": nodes,
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
	addr := ":" + port
	r := mux.NewRouter()
	r.HandleFunc("/search", searchHandler).Methods("GET")
	log.Printf("Server started at %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
