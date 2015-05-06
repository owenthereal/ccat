package router

func MapToArray(m map[string]string) (a []string) {
	for k, v := range m {
		a = append(a, k, v)
	}
	return
}
