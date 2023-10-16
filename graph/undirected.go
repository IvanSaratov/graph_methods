package graph

import (
	"errors"
)

// Имлементация интерфейса для ненапрвленного графа
type undirected[K comparable, T any] struct {
	hash   Hash[K, T]
	traits *Traits
	store  Store[K, T]
}

// Конструктор создания
func newUndirected[K comparable, T any](hash Hash[K, T], traits *Traits, store Store[K, T]) *undirected[K, T] {
	return &undirected[K, T]{
		hash:   hash,
		traits: traits,
		store:  store,
	}
}

// Гетер для получения черт графа
func (u *undirected[K, T]) Traits() *Traits {
	return u.traits
}

func (u *undirected[K, T]) AddVertex(value T, options ...func(*VertexProperties)) error {
	hash := u.hash(value)

	prop := VertexProperties{
		Weight:     0,
		Attributes: make(map[string]string),
	}

	for _, option := range options {
		option(&prop)
	}

	return u.store.AddVertex(hash, value, prop)
}

func (u *undirected[K, T]) Vertex(hash K) (T, error) {
	vertex, _, err := u.store.Vertex(hash)
	return vertex, err
}

func (u *undirected[K, T]) VertexWithProperties(hash K) (T, VertexProperties, error) {
	vertex, prop, err := u.store.Vertex(hash)
	if err != nil {
		return vertex, VertexProperties{}, err
	}

	return vertex, prop, nil
}

func (u *undirected[K, T]) RemoveVertex(hash K) error {
	return u.store.RemoveVertex(hash)
}

func (u *undirected[K, T]) AddEdge(source, target K, options ...func(*EdgeProperties)) error {
	if _, _, err := u.store.Vertex(source); err != nil {
		return err
	}

	if _, _, err := u.store.Vertex(target); err != nil {
		return err
	}

	//nolint:govet // False positive.
	if _, err := u.Edge(source, target); !errors.Is(err, ErrorEdgeNotFound) {
		return ErrorEdgeExists
	}

	edge := Edge[K]{
		Source: source,
		Target: target,
		Properties: EdgeProperties{
			Attributes: make(map[string]string),
		},
	}

	for _, option := range options {
		option(&edge.Properties)
	}

	if err := u.addEdge(source, target, edge); err != nil {
		return err
	}

	return nil
}

func (u *undirected[K, T]) AddEdgesFrom(g Graph[K, T]) error {
	edges, err := g.Edges()
	if err != nil {
		return err
	}

	for _, edge := range edges {
		if err := u.AddEdge(copyEdge(edge)); err != nil {
			return err
		}
	}

	return nil
}

func (u *undirected[K, T]) AddVerticesFrom(g Graph[K, T]) error {
	adjacencyMap, err := g.AdjacencyMap()
	if err != nil {
		return err
	}

	for hash := range adjacencyMap {
		vertex, properties, err := g.VertexWithProperties(hash)
		if err != nil {
			return err
		}

		if err = u.AddVertex(vertex, copyVertexProperties(properties)); err != nil {
			return err
		}
	}

	return nil
}

func (u *undirected[K, T]) Edge(source, target K) (Edge[T], error) {
	// In an undirected graph, since multigraphs aren't supported, the edge AB
	// is the same as BA. Therefore, if source[target] cannot be found, this
	// function also looks for target[source].
	edge, err := u.store.Edge(source, target)
	if errors.Is(err, ErrorEdgeNotFound) {
		edge, err = u.store.Edge(target, source)
	}

	if err != nil {
		return Edge[T]{}, err
	}

	sourceVertex, _, err := u.store.Vertex(source)
	if err != nil {
		return Edge[T]{}, err
	}

	targetVertex, _, err := u.store.Vertex(target)
	if err != nil {
		return Edge[T]{}, err
	}

	return Edge[T]{
		Source: sourceVertex,
		Target: targetVertex,
		Properties: EdgeProperties{
			Weight:     edge.Properties.Weight,
			Attributes: edge.Properties.Attributes,
			Data:       edge.Properties.Data,
		},
	}, nil
}

type tuple[K comparable] struct {
	source, target K
}

func (u *undirected[K, T]) Edges() ([]Edge[K], error) {
	storedEdges, err := u.store.ListEdges()
	if err != nil {
		return nil, err
	}

	edges := make([]Edge[K], 0, len(storedEdges)/2)
	added := make(map[tuple[K]]struct{})

	for _, storedEdge := range storedEdges {
		reversedEdge := tuple[K]{
			source: storedEdge.Target,
			target: storedEdge.Source,
		}
		if _, ok := added[reversedEdge]; ok {
			continue
		}

		edges = append(edges, storedEdge)

		addedEdge := tuple[K]{
			source: storedEdge.Source,
			target: storedEdge.Target,
		}

		added[addedEdge] = struct{}{}
	}

	return edges, nil
}

func (u *undirected[K, T]) EditEdge(source, target K, options ...func(properties *EdgeProperties)) error {
	existingEdge, err := u.store.Edge(source, target)
	if err != nil {
		return err
	}

	for _, option := range options {
		option(&existingEdge.Properties)
	}

	if err := u.store.EditEdge(source, target, existingEdge); err != nil {
		return err
	}

	reversedEdge := existingEdge
	reversedEdge.Source = existingEdge.Target
	reversedEdge.Target = existingEdge.Source

	return u.store.EditEdge(target, source, reversedEdge)
}

func (u *undirected[K, T]) RemoveEdge(source, target K) error {
	if _, err := u.Edge(source, target); err != nil {
		return err
	}

	if err := u.store.RemoveEdge(source, target); err != nil {
		return err
	}

	if err := u.store.RemoveEdge(target, source); err != nil {
		return err
	}

	return nil
}

func (u *undirected[K, T]) AdjacencyMap() (map[K]map[K]Edge[K], error) {
	vertices, err := u.store.ListVertices()
	if err != nil {
		return nil, err
	}

	edges, err := u.store.ListEdges()
	if err != nil {
		return nil, err
	}

	m := make(map[K]map[K]Edge[K], len(vertices))

	for _, vertex := range vertices {
		m[vertex] = make(map[K]Edge[K])
	}

	for _, edge := range edges {
		m[edge.Source][edge.Target] = edge
	}

	return m, nil
}

func (u *undirected[K, T]) Clone() (Graph[K, T], error) {
	traits := &Traits{
		IsDirected: u.traits.IsDirected,
		IsWeighted: u.traits.IsWeighted,
		IsRooted:   u.traits.IsRooted,
	}

	clone := &undirected[K, T]{
		hash:   u.hash,
		traits: traits,
		store:  newStore[K, T](),
	}

	if err := clone.AddVerticesFrom(u); err != nil {
		return nil, err
	}

	if err := clone.AddEdgesFrom(u); err != nil {
		return nil, err
	}

	return clone, nil
}

func (u *undirected[K, T]) Order() (int, error) {
	return u.store.VertexCount()
}

func (u *undirected[K, T]) Size() (int, error) {
	size := 0

	outEdges, err := u.AdjacencyMap()
	if err != nil {
		return 0, err
	}

	for _, outEdges := range outEdges {
		size += len(outEdges)
	}

	// Divide by 2 since every add edge operation on undirected graph is counted
	// twice.
	return size / 2, nil
}

func (u *undirected[K, T]) addEdge(source, target K, edge Edge[K]) error {
	err := u.store.AddEdge(source, target, edge)
	if err != nil {
		return err
	}

	rEdge := Edge[K]{
		Source: edge.Target,
		Target: edge.Source,
		Properties: EdgeProperties{
			Weight:     edge.Properties.Weight,
			Attributes: edge.Properties.Attributes,
			Data:       edge.Properties.Data,
		},
	}

	err = u.store.AddEdge(target, source, rEdge)
	if err != nil {
		return err
	}

	return nil
}
