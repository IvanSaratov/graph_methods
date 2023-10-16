package main

import (
	"fmt"
	"os"

	"github.com/IvanSaratov/graph_methods/draw"
	"github.com/IvanSaratov/graph_methods/graph"
)

func main() {
	g := graph.New(graph.IntHash, graph.Directed())

	_ = g.AddVertex(1)
	_ = g.AddVertex(2)
	_ = g.AddVertex(3)
	_ = g.AddVertex(4)
	_ = g.AddVertex(5)
	_ = g.AddVertex(6)

	_ = g.AddEdge(1, 2)
	_ = g.AddEdge(1, 3)
	_ = g.AddEdge(2, 3)
	_ = g.AddEdge(3, 4)
	_ = g.AddEdge(3, 2)
	_ = g.AddEdge(3, 5)
	_ = g.AddEdge(4, 5)
	_ = g.AddEdge(5, 6)

	index := 4 // Наша заданная вершина
	// Достаем все дуги
	edges, _ := g.Edges()

	sum := make(map[int]int)
	// Инициализиурем нашу карту, во избежаннии пустых заходов
	for _, edge := range edges {
		sum[edge.Source] = 0
	}
	for _, edge := range edges {
		sum[edge.Target]++
	}

	for edge, count := range sum {
		if edge != index && count < sum[index] {
			fmt.Printf("Вершина у которой полустепени заходи меньше чем у %d: %v\n", index, edge)
		}
	}

	file, _ := os.Create("./test.gv")
	_ = draw.DOT(g, file)
}
