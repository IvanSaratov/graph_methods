package graph

import "errors"

// Вынесенные заранее ошибки
var (
	ErrorVertexExists    = errors.New("Вершина уже существует")
	ErrorVertextNotFound = errors.New("Вершина не найден")
	ErrorEdgeExists      = errors.New("Дуг уже существует")
	ErrorEdgeNotFound    = errors.New("Дуга не найдена")

	ErrorVertexHashEdges = errors.New("У вершины ещё есть дуги")
)
