package graph

// Структура нашего графа
type Graph[K comparable, T any] interface {
	// Для указания какой граф будет
	// Передается через интерфейс Traits при инициализации гарфа
	Traits() *Traits

	// Добавляет новую вершину
	AddVertex(value T, options ...func(*VertexProperties)) error
	// Возврашает нам вершину
	Vertex(hash K) (T, error)
	// Удаляет вершину
	RemoveVertex(hash K) error
	// Добавляет новую дугу
	AddEdge(sorce, target K, options ...func(*EdgeProperties)) error
	// Выводит дугу содененную двумя вершинами soruce и target
	Edge(source, target K) (Edge[T], error)
	// Возвращает все дуги графа
	Edges() ([]Edge[K], error)
	// Обновляет данные о дуге находящуюся между думая вершинами source и target
	EditEdge(source, target K, options ...func(properties *EdgeProperties)) error
	// Удаляет дугу между дух вершин source и target
	RemoveEdge(source, target K) error

	// Дополнительные функции для копирования существуюшего графа в новый
	AddVerticesFrom(g Graph[K, T]) error
	AddEdgesFrom(g Graph[K, T]) error

	// Карта смежности
	AdjacencyMap() (map[K]map[K]Edge[K], error)
	// Возврашает вершину с его доп переменными
	VertexWithProperties(hash K) (T, VertexProperties, error)
	// Дополнительные функции
	// Клонирование графа
	Clone() (Graph[K, T], error)

	// Функция возврата количества вершин в графе
	Order() (int, error)

	// Функци возврата количества дуг в графе
	Size() (int, error)
}

// Дополнительная структура данных, представляемые динамические параметры вершины
// Например можно указывать цвет, запах, любимые игры
// graph.VertexProps("color": "red")
type VertexProperties struct {
	Attributes map[string]string
	Weight     int
}

// Структура дуги
// Принимает любое значение.
// Указывает на две вершины которые соединяет source и target
// Так же имет структуру с дополнительными полями
type Edge[T any] struct {
	Source     T
	Target     T
	Properties EdgeProperties
}

// Похожая на VertexPoperties, но с добавленным полем Data
type EdgeProperties struct {
	Attributes map[string]string
	Weight     int
	Data       any
}

// Весы для графа
func EdgeWeight(weight int) func(*EdgeProperties) {
	return func(e *EdgeProperties) {
		e.Weight = weight
	}
}

// Обертка для преобразования K в T
type Hash[K comparable, T any] func(T) K

func New[K comparable, T any](hash Hash[K, T], options ...func(*Traits)) Graph[K, T] {
	var t Traits

	for _, option := range options {
		option(&t)
	}

	if t.IsDirected {
		return newDirected(hash, &t, newStore[K, T]())
	}

	return newUndirected(hash, &t, newStore[K, T]())
}

// Функция создания нового графа по входному другому графу
func NewLike[K comparable, T any](g Graph[K, T]) Graph[K, T] {
	copyTraits := func(t *Traits) {
		t.IsDirected = g.Traits().IsDirected
		t.IsRooted = g.Traits().IsRooted
		t.IsWeighted = g.Traits().IsWeighted
	}

	var hash Hash[K, T]

	if g.Traits().IsDirected {
		hash = g.(*directed[K, T]).hash
	} else {
		hash = g.(*undirected[K, T]).hash
	}

	return New(hash, copyTraits)
}

// Функции определения чем заполняется граф

func StringHash(s string) string {
	return s
}

func IntHash(i int) int {
	return i
}
