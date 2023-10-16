package graph

// Фапйл с дополнительными функциями которые не подошли логически

// Копирует все доп. атрибуты
func copyVertexProperties(source VertexProperties) func(*VertexProperties) {
	return func(p *VertexProperties) {
		for k, v := range source.Attributes {
			p.Attributes[k] = v
		}
		p.Weight = source.Weight
	}
}
