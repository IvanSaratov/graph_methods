package draw

import (
	"io"
	"text/template"

	"github.com/IvanSaratov/graph_methods/graph"
)

const dotTemplate = `strict {{.GraphType}} {
{{range $k, $v := .Attributes}}
	{{$k}}="{{$v}}";
{{end}}
{{range $s := .Statements}}
	"{{.Source}}" {{if .Target}}{{$.EdgeOperator}} "{{.Target}}" [ {{range $k, $v := .EdgeAttributes}}{{$k}}="{{$v}}", {{end}} weight={{.EdgeWeight}} ]{{else}}[ {{range $k, $v := .SourceAttributes}}{{$k}}="{{$v}}", {{end}} weight={{.SourceWeight}} ]{{end}};
{{end}}
}
`

type description struct {
	GraphType    string
	Attributes   map[string]string
	EdgeOperator string
	Statements   []statement
}

type statement struct {
	Source           interface{}
	Target           interface{}
	SourceWeight     int
	SourceAttributes map[string]string
	EdgeWeight       int
	EdgeAttributes   map[string]string
}

func DOT[K comparable, T any](g graph.Graph[K, T], w io.Writer, options ...func(*description)) error {
	desc, err := generateDOT(g, options...)
	if err != nil {
		return err
	}

	return renderDOT(w, desc)
}

// GraphAttribute is a functional option for the [DOT] method.
func GraphAttribute(key, value string) func(*description) {
	return func(d *description) {
		d.Attributes[key] = value
	}
}

func generateDOT[K comparable, T any](g graph.Graph[K, T], options ...func(*description)) (description, error) {
	desc := description{
		GraphType:    "graph",
		Attributes:   make(map[string]string),
		EdgeOperator: "--",
		Statements:   make([]statement, 0),
	}

	for _, option := range options {
		option(&desc)
	}

	if g.Traits().IsDirected {
		desc.GraphType = "digraph"
		desc.EdgeOperator = "->"
	}

	adjacencyMap, err := g.AdjacencyMap()
	if err != nil {
		return desc, err
	}

	for vertex, adjacencies := range adjacencyMap {
		_, sourceProperties, err := g.VertexWithProperties(vertex)
		if err != nil {
			return desc, err
		}

		stmt := statement{
			Source:           vertex,
			SourceWeight:     sourceProperties.Weight,
			SourceAttributes: sourceProperties.Attributes,
		}
		desc.Statements = append(desc.Statements, stmt)

		for adjacency, edge := range adjacencies {
			stmt := statement{
				Source:         vertex,
				Target:         adjacency,
				EdgeWeight:     edge.Properties.Weight,
				EdgeAttributes: edge.Properties.Attributes,
			}
			desc.Statements = append(desc.Statements, stmt)
		}
	}

	return desc, nil
}

func renderDOT(w io.Writer, d description) error {
	tpl, err := template.New("dotTemplate").Parse(dotTemplate)
	if err != nil {
		return err
	}

	return tpl.Execute(w, d)
}
