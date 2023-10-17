package graph

// Функция обхода в глубину. Использует не рекурсивный метод через кучу
func DFS[K comparable, T any](g Graph[K, T], start K, visit func(K) bool) error {
	// Строит карту смежности
	adjacencyMap, err := g.AdjacencyMap()
	if err != nil {
		return err
	}

	// Находим в карте смежности заданную нами вершину
	if _, ok := adjacencyMap[start]; !ok {
		return err
	}

	// Инициируем кучу
	stack := newStack[K]()
	// Создаем переменные под посещенные вершины
	visited := make(map[K]bool)

	// Кладем нашу начальную вершину
	stack.push(start)
	// Пока куча не кончилась
	for !stack.isEmpty() {
		current, _ := stack.pop()

		if _, ok := visited[current]; !ok {
			// Заканчиваем обход если мы уже посещали данную вершину
			if stop := visit(current); stop {
				break
			}
			visited[current] = true

			// Кладем следущю вершину
			for adjacency := range adjacencyMap[current] {
				stack.push(adjacency)
			}
		}
	}

	return nil
}
