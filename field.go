package errs

import (
	"fmt"
	"time"
)

// Field defines a key-value pair.
type Field interface {
	Key() string
	Value() interface{}
}

// F creates a new Field with the given key and value.
func F(key string, val interface{}) Field {
	return field{key: key, val: val}
}

type field struct {
	key string
	val interface{}
}

// Key gives the key of field.
func (f field) Key() string {
	return f.key
}

// Value gives the value of field.
func (f field) Value() interface{} {
	return f.val
}

func valueString(f Field) string {
	natural := fmt.Sprintf("%v", f.Value())

	switch val := f.Value().(type) {
	case nil:
		return natural
	case time.Time:
		return val.Format(time.RFC3339)
	case *time.Time:
		return val.Format(time.RFC3339)
	default:
		return natural
	}
}
