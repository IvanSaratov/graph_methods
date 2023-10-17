package main

import (
	"fmt"
	"os"

	"github.com/IvanSaratov/graph_methods/draw"
	"github.com/IvanSaratov/graph_methods/graph"
)

func main() {
	g := graph.New(graph.IntHash, graph.Directed())

	//Создание обычного графа
	_ = g.AddVertex(1)
	_ = g.AddVertex(2)
	_ = g.AddVertex(3)
	_ = g.AddVertex(4)
	_ = g.AddVertex(5)

	_ = g.AddEdge(1, 2)
	_ = g.AddEdge(1, 3)
	_ = g.AddEdge(2, 4)
	_ = g.AddEdge(4, 5)
	_ = g.AddEdge(5, 1)

	index := 4 // Заданная вершина
	vertexs, _ := g.AdjacencyMap()
	for vertex := range vertexs {
		if vertex == index {
			continue
		}

		_ = graph.BFS(g, vertex, func(value int) bool {
			if value == index {
				fmt.Printf("Из вершины %v можно попасть в вершину %v\n", vertex, index)
				return true
			}
			return false
		})

	}
	file, _ := os.Create("./test.gv")
	_ = draw.DOT(g, file)

}
