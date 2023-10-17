package graph

func BFS[K comparable, T any](g Graph[K, T], start K, visit func(K) bool) error {
	ignoreDepth := func(vertex K, _ int) bool {
		return visit(vertex)
	}
	return BFSWithDepth(g, start, ignoreDepth)
}

// Функция обхода в ширину. Не рекурсивная. Использует очередь
// Принимает в себя функцию как аргумент с ограничием глубины. Если аргумент пропустить
func BFSWithDepth[K comparable, T any](g Graph[K, T], start K, visit func(K, int) bool) error {
	adjacencyMap, err := g.AdjacencyMap()
	if err != nil {
		return err
	}

	if _, ok := adjacencyMap[start]; !ok {
		return err
	}

	queue := make([]K, 0)
	visited := make(map[K]bool)

	visited[start] = true
	queue = append(queue, start)
	depth := 0

	for len(queue) > 0 {
		current := queue[0]

		queue = queue[1:]
		depth++

		// Останавливаем поиск если достигли пройденного
		if stop := visit(current, depth); stop {
			break
		}

		for adjacency := range adjacencyMap[current] {
			if _, ok := visited[adjacency]; !ok {
				visited[adjacency] = true
				queue = append(queue, adjacency)
			}
		}

	}

	return nil
}
