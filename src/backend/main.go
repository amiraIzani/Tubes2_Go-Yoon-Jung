package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"

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
	data, _ := ioutil.ReadFile(path)
	var raws []Raw
	_ = json.Unmarshal(data, &raws)

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

func parallelExpand(current *Tree, parentPairs [][2]string, visited *int32) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, pair := range parentPairs {
		wg.Add(1)
		go func(p [2]string) {
			defer wg.Done()

			left := &Tree{Name: p[0]}
			right := &Tree{Name: p[1]}

			mu.Lock()
			current.Children = append(current.Children, left, right)
			mu.Unlock()

			atomic.AddInt32(visited, 2)
		}(pair)
	}
	wg.Wait()
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	target := q.Get("target")
	method := q.Get("method")
	mode := q.Get("mode")
	if target == "" || method == "" || mode == "" {
		http.Error(w, "Parameter kurang", http.StatusBadRequest)
		return
	}

	limit := 1
	if mode == "multiple" {
		lv, err := strconv.Atoi(q.Get("limit"))
		if err != nil || lv < 1 {
			http.Error(w, "Limit tidak valid", http.StatusBadRequest)
			return
		}
		limit = lv
	}

	var recipes []*Tree
	var elapsed int64
	var nodesVisited int

	switch mode {
	case "one":
		if method == "dfs" {
			t, e, v := dfsSearch(target)
			recipes = []*Tree{t}
			elapsed, nodesVisited = e, v
		} else {
			t, e, v := bfsSearch(target)
			recipes = []*Tree{t}
			elapsed, nodesVisited = e, v
		}
	case "multiple":
		recipes, elapsed, nodesVisited = multiSearch(target, method, limit)
	default:
		http.Error(w, "Mode tidak didukung", http.StatusBadRequest)
		return
	}

	resp := map[string]interface{}{ 
		"method":        method,
		"mode":          mode,
		"limit":         limit,
		"time_us":       elapsed,
		"nodes_visited": nodesVisited,
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
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./fe")))

	handler := handlers.CORS(handlers.AllowedOrigins([]string{"*"}))(r)
	fmt.Printf("Server berjalan di port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
