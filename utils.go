package go_requests

func updateMap(origin map[string]string, new map[string]string) {
	for key, value := range new {
		origin[key] = value
	}
}
