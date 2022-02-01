package cache

var cache map[string]interface{}

func getOrAdd(key string, add func() interface{}) interface{} {
	exists := cache[key]
}

func add(key string)

// TODO: make nice. For example ability to call cache.GetOrFetch(), cache.ForceFetch() from outside this package?
// nevermind the above, empty interface crap should stay encapsulated within this package.
// so do something like cache.Catalog() and cache.ForceFetchCatalog() ? just see
