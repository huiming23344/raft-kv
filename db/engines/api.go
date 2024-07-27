package engines

type KvsEngine interface {

	// Set the value of a string key to a string.
	// If the key already exists, the previous value will be overwritten.
	Set(key, value string) error

	// Get the string value of a given string key.
	// Return `None` if the given key does not exits.
	Get(key string) (string, error)

	// Remove a given key.
	//  It returns `kvserror::KeyNotFound` if the given key not found.
	Remove(key string) error
}
