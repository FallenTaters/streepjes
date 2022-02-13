package cache

var cache = map[string]interface{}{}

func getOrAdd(key string, addFunc func() (interface{}, error)) (interface{}, error) {
	data, exists := cache[key]
	if exists {
		return data, nil
	}

	data, err := addFunc()
	if err != nil {
		return nil, err
	}

	add(key, data)

	return data, nil
}

func add(key string, v interface{}) {
	cache[key] = v
}
