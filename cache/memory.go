package cache

type InMemoryCache struct {
	// maps are passed by reference by default
	records map[string]int64
}

func NewInMemoryCache() InMemoryCache {
	records := make(map[string]int64)
	return InMemoryCache{records: records}
}

func (r *InMemoryCache) FlushData() error {
	r.records = make(map[string]int64)

	return nil
}

func (r *InMemoryCache) GetKey(k string) (int64, error) {
	v, ok := r.records[k]
	if !ok {
		return 0, ErrKeyNotFound
	}

	return v, nil
}

func (r *InMemoryCache) SetKey(k string, v int64) error {
	r.records[k] = v

	return nil
}
