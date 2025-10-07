package graph

type Graph struct {
	adjacency map[string][]string
	directed  bool
}

func (g *Graph) AddEdge(u, v string) {
	g.adjacency[u] = append(g.adjacency[u], v)
	if !g.directed {
		g.adjacency[v] = append(g.adjacency[v], u)
	}
}

func (g *Graph) Neighbors(n string) []string {
	return g.adjacency[n]
}

func (g *Graph) Exist(n string) bool {
	_, exist := g.adjacency[n]
	return exist
}
