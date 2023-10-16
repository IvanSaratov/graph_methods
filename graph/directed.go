package graph

import (
	"errors"
)

type directed[K comparable, T any] struct {
	hash   Hash[K, T]
	traits *Traits
	store  Store[K, T]
}

func newDirected[K comparable, T any](hash Hash[K, T], traits *Traits, store Store[K, T]) *directed[K, T] {
	return &directed[K, T]{
		hash:   hash,
		traits: traits,
		store:  store,
	}
}

func (d *directed[K, T]) Traits() *Traits {
	return d.traits
}

func (d *directed[K, T]) AddVertex(value T, options ...func(*VertexProperties)) error {
	hash := d.hash(value)
	properties := VertexProperties{
		Weight:     0,
		Attributes: make(map[string]string),
	}

	for _, option := range options {
		option(&properties)
	}

	return d.store.AddVertex(hash, value, properties)
}

func (d *directed[K, T]) AddVerticesFrom(g Graph[K, T]) error {
	adjacencyMap, err := g.AdjacencyMap()
	if err != nil {
		return err
	}

	for hash := range adjacencyMap {
		vertex, properties, err := g.VertexWithProperties(hash)
		if err != nil {
			return err
		}

		if err = d.AddVertex(vertex, copyVertexProperties(properties)); err != nil {
			return err
		}
	}

	return nil
}

// Возврашает вершину
func (d *directed[K, T]) Vertex(hash K) (T, error) {
	vertex, _, err := d.store.Vertex(hash)
	return vertex, err
}

func (d *directed[K, T]) VertexWithProperties(hash K) (T, VertexProperties, error) {
	vertex, properties, err := d.store.Vertex(hash)
	if err != nil {
		return vertex, VertexProperties{}, err
	}

	return vertex, properties, nil
}

// Удаляем вершину
func (d *directed[K, T]) RemoveVertex(hash K) error {
	return d.store.RemoveVertex(hash)
}

// Добавлеяем дугу
func (d *directed[K, T]) AddEdge(source, target K, options ...func(*EdgeProperties)) error {
	// Две проверки на существование вершины
	_, _, err := d.store.Vertex(source)
	if err != nil {
		return err
	}

	_, _, err = d.store.Vertex(target)
	if err != nil {
		return err
	}

	if _, err := d.Edge(source, target); !errors.Is(err, ErrorEdgeNotFound) {
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

	return d.addEdge(source, target, edge)
}

func (d *directed[K, T]) AddEdgesFrom(g Graph[K, T]) error {
	edges, err := g.Edges()
	if err != nil {
		return err
	}

	for _, edge := range edges {
		if err := d.AddEdge(copyEdge(edge)); err != nil {
			return err
		}
	}

	return nil
}

func (d *directed[K, T]) Edge(source, target K) (Edge[T], error) {
	edge, err := d.store.Edge(source, target)
	if err != nil {
		return Edge[T]{}, err
	}

	sourceVertex, _, err := d.store.Vertex(source)
	if err != nil {
		return Edge[T]{}, err
	}

	targetVertex, _, err := d.store.Vertex(target)
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

func (d *directed[K, T]) Edges() ([]Edge[K], error) {
	return d.store.ListEdges()
}

func (d *directed[K, T]) EditEdge(source, target K, options ...func(properties *EdgeProperties)) error {
	existingEdge, err := d.store.Edge(source, target)
	if err != nil {
		return err
	}

	for _, option := range options {
		option(&existingEdge.Properties)
	}

	return d.store.EditEdge(source, target, existingEdge)
}

func (d *directed[K, T]) RemoveEdge(source, target K) error {
	if _, err := d.Edge(source, target); err != nil {
		return err
	}

	if err := d.store.RemoveEdge(source, target); err != nil {
		return err
	}

	return nil
}

func (d *directed[K, T]) AdjacencyMap() (map[K]map[K]Edge[K], error) {
	vertices, err := d.store.ListVertices()
	if err != nil {
		return nil, err
	}

	edges, err := d.store.ListEdges()
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

func (d *directed[K, T]) PredecessorMap() (map[K]map[K]Edge[K], error) {
	vertices, err := d.store.ListVertices()
	if err != nil {
		return nil, err
	}

	edges, err := d.store.ListEdges()
	if err != nil {
		return nil, err
	}

	m := make(map[K]map[K]Edge[K], len(vertices))

	for _, vertex := range vertices {
		m[vertex] = make(map[K]Edge[K])
	}

	for _, edge := range edges {
		if _, ok := m[edge.Target]; !ok {
			m[edge.Target] = make(map[K]Edge[K])
		}
		m[edge.Target][edge.Source] = edge
	}

	return m, nil
}

func (d *directed[K, T]) addEdge(source, target K, edge Edge[K]) error {
	return d.store.AddEdge(source, target, edge)
}

func (d *directed[K, T]) Clone() (Graph[K, T], error) {
	traits := &Traits{
		IsDirected: d.traits.IsDirected,
		IsWeighted: d.traits.IsWeighted,
		IsRooted:   d.traits.IsRooted,
	}

	clone := &directed[K, T]{
		hash:   d.hash,
		traits: traits,
		store:  newStore[K, T](),
	}

	if err := clone.AddVerticesFrom(d); err != nil {
		return nil, err
	}

	if err := clone.AddEdgesFrom(d); err != nil {
		return nil, err
	}

	return clone, nil
}

func (d *directed[K, T]) Order() (int, error) {
	return d.store.VertexCount()
}

func (d *directed[K, T]) Size() (int, error) {
	size := 0
	outEdges, err := d.AdjacencyMap()
	if err != nil {
		return 0, err
	}

	for _, outEdges := range outEdges {
		size += len(outEdges)
	}

	return size, nil
}

func copyEdge[K comparable](edge Edge[K]) (K, K, func(properties *EdgeProperties)) {
	copyProperties := func(p *EdgeProperties) {
		for k, v := range edge.Properties.Attributes {
			p.Attributes[k] = v
		}
		p.Weight = edge.Properties.Weight
		p.Data = edge.Properties.Data
	}

	return edge.Source, edge.Target, copyProperties
}
