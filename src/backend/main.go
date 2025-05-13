package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

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

var (
	graph        map[string]Node
	baseElements = map[string]struct{}{"Fire": {}, "Earth": {}, "Air": {}, "Water": {}}
)

func loadRecipes(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf(err.Error())
	}
	var raws []Raw
	err = json.Unmarshal(data, &raws)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	graph = make(map[string]Node, len(raws))
	for _, item := range raws {
		n := Node{Name: item.Name}
		for _, r := range item.Recipes {
			if len(r.Elements) == 2 {
				a, b := r.Elements[0], r.Elements[1]
				dup := false
				for _, p := range n.Parents {
					if (p[0] == a && p[1] == b) || (p[0] == b && p[1] == a) {
						dup = true
						break
					}
				}
				if !dup {
					n.Parents = append(n.Parents, [2]string{a, b})
				}
			}
		}
		graph[item.Name] = n
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	target := q.Get("target")
	method := q.Get("method")
	mode := q.Get("mode")
	if target == "" || method == "" || mode == "" {
		http.Error(w, "parameter tidak lengkap", http.StatusBadRequest)
		return
	}

	limit := 1
	if mode == "multiple" {
		if lv, err := strconv.Atoi(q.Get("limit")); err == nil && lv > 0 {
			limit = lv
		} else {
			http.Error(w, "batas tidak valid", http.StatusBadRequest)
			return
		}
	}

	var (
		resp    interface{}
		elapsed int64
		nodes   int
	)

	switch mode {
	case "shortest":
		if method == "dfs" {
			t, e, v := dfsSearch(target)
			resp, elapsed, nodes = t, e, v
		} else {
			t, e, v := bfsSearch(target)
			resp, elapsed, nodes = t, e, v
		}
	case "multiple":
		ts, e, v := multiSearch(target, method, limit)
		resp, elapsed, nodes = ts, e, v
	default:
		http.Error(w, "mode tidak didukung", http.StatusBadRequest)
		return
	}

	out := map[string]interface{}{
		"method":        method,
		"mode":          mode,
		"limit":         limit,
		"time_us":       elapsed,
		"nodes_visited": nodes,
		"recipes":       resp,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
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
