package main

import (
	"os"

	"github.com/IvanSaratov/graph_methods/draw"
	"github.com/IvanSaratov/graph_methods/graph"
)

func main() {
	g := graph.New(graph.IntHash)

	// Создание обычного графа
	_ = g.AddVertex(1)
	_ = g.AddVertex(2)
	_ = g.AddVertex(3)
	_ = g.AddVertex(4)
	_ = g.AddVertex(5)

	_ = g.AddEdge(1, 2)
	_ = g.AddEdge(1, 3)
	_ = g.AddEdge(2, 4)
	_ = g.AddEdge(4, 5)

	file, _ := os.Create("./test.gv")
	_ = draw.DOT(g, file)

}
