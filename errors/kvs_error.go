package errors

type KvsError struct {
	cause string
}

func (kvs KvsError) Error() string {
	return kvs.cause
}

var KeyNotFound = KvsError{"Key not found"}
