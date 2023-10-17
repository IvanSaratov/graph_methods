package main

import (
	"fmt"
	"os"

	"github.com/IvanSaratov/graph_methods/draw"
	"github.com/IvanSaratov/graph_methods/graph"
)

func main() {
	g := graph.New(graph.IntHash)

	//Создание обычного графа
	_ = g.AddVertex(1)
	_ = g.AddVertex(2)
	_ = g.AddVertex(3)
	_ = g.AddVertex(4)
	_ = g.AddVertex(5)
	// _ = g.AddVertex(6)

	_ = g.AddEdge(1, 2)
	_ = g.AddEdge(1, 3)
	_ = g.AddEdge(2, 4)
	_ = g.AddEdge(4, 5)
	_ = g.AddEdge(5, 1)

	tmp, _ := g.AdjacencyMap()
	size, _ := g.Order()
	for vertex := range tmp {
		count := 0
		_ = graph.DFS(g, vertex, func(value int) bool {
			count++
			return false
		})

		if count != size {
			fmt.Println("Граф не связанный")
			return
		}
	}

	fmt.Println("Граф связанный")

	file, _ := os.Create("./test.gv")
	_ = draw.DOT(g, file)

}
