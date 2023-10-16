package graph

type Traits struct {
	IsDirected bool
	IsWeighted bool
	IsRooted   bool
}

// Фнукция указания что граф направленный
func Directed() func(*Traits) {
	return func(t *Traits) {
		t.IsDirected = true
	}
}

// Функция указания что граф взвешанный
func Weighted() func(*Traits) {
	return func(t *Traits) {
		t.IsWeighted = true
	}
}

// Корневой граф
func Rooted() func(*Traits) {
	return func(t *Traits) {
		t.IsRooted = true
	}
}
