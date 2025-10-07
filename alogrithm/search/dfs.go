package search

// Recursive
func DFS(graph map[string][]string, start string, visited map[string]bool) {
	if visited[start] {
		return
	}
	visited[start] = true
	for _, neighbor := range graph[start] {
		if !visited[neighbor] {
			DFS(graph, neighbor, visited)
		}
	}
}

// Stack
func DFS_(graph map[string][]string, start string) {
	visited := make(map[string]bool)
	stack := []string{start}

	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if visited[node] {
			continue
		}

		visited[node] = true
		neighbors := graph[node]
		for i := len(neighbors) - 1; i >= 0; i-- {
			n := neighbors[i]
			if !visited[n] {
				stack = append(stack, n)
			}
		}
	}
}
