package cache

type Cache interface {

	// Add a kv set to cache.
	// If the key already exists, the previous value will be overwritten.
	Set(key, value string)

	// Get the string value of a given string key.
	// Return `None` if the given key does not exit.
	Get(key string) (string, bool)

	// Remove a given key.
	Remove(key string)
}
