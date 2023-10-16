package graph

import (
	"sync"
)

type Store[K comparable, T any] interface {
	// Добавление вершины
	AddVertex(hash K, value T, properties VertexProperties) error
	// Возврашает информацию о вешине
	Vertex(hash K) (T, VertexProperties, error)
	// Удаление вершины по ключу
	RemoveVertex(hash K) error
	// Возврашает список вершин
	ListVertices() ([]K, error)
	// Козвршает количество вершин в грфае
	VertexCount() (int, error)
	// Создает Дугу
	AddEdge(sourceHash, targetHash K, edge Edge[K]) error
	// Редактирует параметры дуги
	EditEdge(sourceHash, targetHash K, edge Edge[K]) error
	// Удаляет дугу
	RemoveEdge(sourceHash, targetHash K) error
	// Возврашает данные о дуге
	Edge(sourceHash, targetHash K) (Edge[K], error)
	// Возврашает список всег дуг в графе
	ListEdges() ([]Edge[K], error)
}

// Структура хранения графа
type store[K comparable, T any] struct {
	lock             sync.RWMutex
	vertices         map[K]T
	VertexProperties map[K]VertexProperties

	inEdges  map[K]map[K]Edge[K] // target -> source
	outEdges map[K]map[K]Edge[K] //source -> target
}

// Конструкт инициализации графа
func newStore[K comparable, T any]() Store[K, T] {
	return &store[K, T]{
		vertices:         make(map[K]T),
		VertexProperties: make(map[K]VertexProperties),
		inEdges:          make(map[K]map[K]Edge[K]),
		outEdges:         make(map[K]map[K]Edge[K]),
	}
}

func (s *store[K, T]) AddVertex(key K, value T, props VertexProperties) error {
	// Против гонки за ресурсами
	s.lock.Lock()
	defer s.lock.Unlock()

	// Проверка на существование вершины
	if _, ok := s.vertices[key]; ok {
		return ErrorVertexExists
	}

	s.vertices[key] = value
	s.VertexProperties[key] = props

	return nil
}

func (s *store[K, T]) ListVertices() ([]K, error) {
	// Блокируем чтение
	s.lock.RLock()
	defer s.lock.RUnlock()

	hashes := make([]K, 0, len(s.vertices))
	for key := range s.vertices {
		hashes = append(hashes, key)
	}

	return hashes, nil
}

func (s *store[K, T]) VertexCount() (int, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	// Возврашаем количество вершин
	return len(s.vertices), nil
}

func (s *store[K, T]) Vertex(key K) (T, VertexProperties, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	value, ok := s.vertices[key]
	if !ok {
		return value, VertexProperties{}, ErrorVertextNotFound
	}

	props := s.VertexProperties[key]

	return value, props, nil
}

func (s *store[K, T]) RemoveVertex(key K) error {
	s.lock.RLock()
	defer s.lock.RUnlock()

	// Проверка на существование вершины
	if _, ok := s.vertices[key]; !ok {
		return ErrorVertextNotFound
	}

	// Проверка на дуги
	if edges, ok := s.inEdges[key]; ok {
		// Если дуг больше нуля, сообщаем об этом и ничего не делаем
		if len(edges) > 0 {
			return ErrorVertexHashEdges
		}
		// Удаляем наши дуги
		delete(s.inEdges, key)
	}

	// Повторяем для направленных наружу дуг
	if edges, ok := s.outEdges[key]; ok {
		if len(edges) > 0 {
			return ErrorVertexHashEdges
		}

		delete(s.outEdges, key)
	}

	delete(s.vertices, key)
	delete(s.VertexProperties, key)

	return nil
}

func (s *store[K, T]) AddEdge(source, target K, edge Edge[K]) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	// Если дуг нету, инициализируем новую дугу
	if _, ok := s.outEdges[source]; !ok {
		s.outEdges[source] = make(map[K]Edge[K])
	}

	// Иначе перезапишем
	s.outEdges[source][target] = edge

	// Для направленых внутрб тоже самое
	if _, ok := s.inEdges[target]; !ok {
		s.inEdges[target] = make(map[K]Edge[K])
	}

	s.inEdges[target][source] = edge

	return nil
}

func (s *store[K, T]) EditEdge(soruce, target K, edge Edge[K]) error {
	// Проверяем что такая дуг существует и не соеденяет нужные вершины
	if _, err := s.Edge(soruce, target); err != nil {
		return err
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	// Вводим новые параметры
	s.outEdges[soruce][target] = edge
	s.inEdges[target][soruce] = edge

	return nil
}

func (s *store[K, T]) RemoveEdge(source, target K) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	// Удаляем дугу с обоих концов
	delete(s.inEdges[target], source)
	delete(s.outEdges[source], target)
	return nil
}

func (s *store[K, T]) Edge(source, target K) (Edge[K], error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	// Проверка что дуга вообще существует
	sourceEdges, ok := s.outEdges[source]
	if !ok {
		return Edge[K]{}, ErrorEdgeNotFound
	}

	// Проверяем что она связана
	edge, ok := sourceEdges[target]
	if !ok {
		return Edge[K]{}, ErrorEdgeNotFound
	}

	return edge, nil
}

func (s *store[K, T]) ListEdges() ([]Edge[K], error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	// Создаем карту и проходимся по всему графу
	res := make([]Edge[K], 0)
	for _, edges := range s.outEdges {
		for _, edge := range edges {
			res = append(res, edge)
		}
	}
	return res, nil
}
