package main

import (
	"fmt"
	"sort"

	"github.com/IvanSaratov/graph_methods/graph"
)

func main() {
	g := graph.New(graph.IntHash, graph.Weighted())

	//Создание обычного графа
	_ = g.AddVertex(1)
	_ = g.AddVertex(2)
	_ = g.AddVertex(3)
	_ = g.AddVertex(4)
	_ = g.AddVertex(5)

	_ = g.AddEdge(1, 2, graph.EdgeWeight(2))
	_ = g.AddEdge(2, 3, graph.EdgeWeight(4))
	_ = g.AddEdge(3, 4, graph.EdgeWeight(1))
	_ = g.AddEdge(4, 1, graph.EdgeWeight(3))
	_ = g.AddEdge(5, 2, graph.EdgeWeight(6))
	_ = g.AddEdge(5, 4, graph.EdgeWeight(2))

	// Достаем список всех дуг
	edges, _ := g.Edges()
	// Сортируем веса
	sort.SliceStable(edges, func(i, j int) bool {
		return edges[i].Properties.Weight < edges[j].Properties.Weight
	})

	// Созадем tree_id - список номер деревье
	n, _ := g.Order()
	var tree_id = make([]int, n+1)
	for i := 1; i <= n; i++ {
		tree_id[i] = i
	}

	var result [][2]int

	cost := 0
	// Количество дуг в графе
	m, _ := g.Size()
	for i := 0; i < m; i++ {
		if tree_id[edges[i].Source] != tree_id[edges[i].Target] {
			cost++
			result = append(result, [2]int{edges[i].Source, edges[i].Target})
			old_id := tree_id[edges[i].Target]
			new_id := tree_id[edges[i].Source]

			for j := 0; j < n; j++ {
				if tree_id[j] == old_id {
					tree_id[j] = new_id
				}
			}
		}
	}

	fmt.Println(result)

	// file, _ := os.Create("./test.gv")
	// _ = draw.DOT(g, file)
}
