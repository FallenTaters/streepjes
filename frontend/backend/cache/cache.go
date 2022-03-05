package cache

import "time"

type value struct {
	data    interface{}
	expires time.Time
}

var cache = map[string]value{}

func getOrAdd(key string, duration time.Duration, addFunc func() (interface{}, error)) (interface{}, error) {
	val, exists := cache[key]
	if exists && val.expires.After(time.Now()) {
		return val.data, nil
	}

	data, err := addFunc()
	if err != nil {
		return nil, err
	}

	add(key, value{
		data:    data,
		expires: time.Now().Add(duration),
	})

	return data, nil
}

func add(key string, v value) {
	cache[key] = v
}
