package path

import (
	"reflect"

	"github.com/benpate/derp"
)

/////////////////////////////
// Interface-based getters
/////////////////////////////

// getMapOfInterface retrieves a value from a map[string]interface{}
func (path Path) getMapOfInterface(value map[string]interface{}) (interface{}, error) {

	key := path.Head()

	field, ok := value[key]

	if !ok {
		return nil, derp.New(500, "path.Path.Get", "Map entry does not exist", path, value)
	}

	return path.Tail().Get(field)
}

// getArrayOfInterface retrieves a value from a []interface{}
func (path Path) getArrayOfInterface(value []interface{}) (interface{}, error) {

	index, err := path.Index(len(value))

	if err != nil {
		return nil, derp.Wrap(err, "path.Path.Get", "Invalid index", path, value)
	}

	return path.Tail().Get(value[index])
}

/////////////////////////////
// Reflect-based getters
/////////////////////////////

// getReflect retrieves a value from a reflect.Value
func (path Path) getReflect(value reflect.Value) (interface{}, error) {

	switch value.Kind() {

	case reflect.Array:
		return path.getReflectArray(value)

	case reflect.Slice:
		return path.getReflectArray(value)

	case reflect.Ptr:
		return path.getReflectPointer(value)

	case reflect.Map:
		return path.getReflectMap(value)

	case reflect.Struct:
		return path.getReflectStruct(value)
	}

	return nil, derp.New(500, "path.Path.Get", "Unsupported data type", path, value)
}

// getReflect retrieves a value from a reflect.Value where Kind = Array or Slice
func (path Path) getReflectArray(value reflect.Value) (interface{}, error) {

	if (value.Kind() != reflect.Array) && (value.Kind() != reflect.Slice) {
		return nil, derp.New(500, "path.Get", "Expected array or slice type", path, value)
	}

	length := value.Len()

	index, err := path.Index(length)

	if err != nil {
		return nil, derp.Wrap(err, "path.Get", "Invalid array index", path, value)
	}

	// Get an interface to the array element
	child := value.Index(index).Interface()

	// Continue traversing the path of the child item.
	return path.Tail().Get(child)
}

// getReflect retrieves a value from a reflect.Value where Kind = Map
func (path Path) getReflectMap(value reflect.Value) (interface{}, error) {
	return nil, derp.New(500, "path.Get", "Unimplemented", path, value)
}

// getReflect retrieves a value from a reflect.Value where Kind = Pointer
func (path Path) getReflectPointer(value reflect.Value) (interface{}, error) {

	if value.Kind() != reflect.Ptr {
		return nil, derp.New(500, "path.Get", "Expected pointer type", path, value)
	}

	// indirect the pointer and continue digging
	return path.getReflect(value.Elem())
}

// getReflect retrieves a value from a reflect.Value where Kind = Struct
func (path Path) getReflectStruct(value reflect.Value) (interface{}, error) {

	// Fail if the data is not a map
	if value.Kind() != reflect.Map {
		return nil, derp.New(500, "path.Get", "Data does not match Schema.  Expected object type", path, value)
	}

	key := path.Head()

	field := value.FieldByName(key)

	// TODO: also look up struct tags?

	if field.IsZero() {
		return nil, derp.New(500, "path.Get", "Struct field does not exist", path, value)
	}

	// traverse the path to the child items
	return path.Tail().Get(field.Interface())
}
