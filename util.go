package vision

func preAllocMap(m *map[int](map[int]int), key int) {
	if _, ok := (*m)[key]; !ok {
		(*m)[key] = map[int]int{}
	}
}

func cantorPairing(x, y int) int {
	return int(0.5 * float64((x+y)*(x+y+1)+y))
}
