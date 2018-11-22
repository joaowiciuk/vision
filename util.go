package vision

func preAllocMap(m *map[int](map[int]int), key int) {
	if _, ok := (*m)[key]; !ok {
		(*m)[key] = map[int]int{}
	}
}
