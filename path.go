package path

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/benpate/derp"
)

// Path is a reference to a value within another data object.
type Path []string

// New creates a new Path object
func New(value string) Path {

	if value == "" {
		return Path([]string{})
	}

	return Path(strings.Split(value, "."))
}

// Get tries to return the value of the object at the provided path.
func Get(object interface{}, path string) (interface{}, error) {
	return New(path).Get(object)
}

// Set tries to set the value of the object at the provided path.
func Set(object interface{}, path string, value interface{}) error {
	return New(path).Set(object, value)
}

// Get tries to return the value of the object at this path.
func (path Path) Get(object interface{}) (interface{}, error) {

	// If the path is empty, then we have reached our goal.  Return the value of this object
	if path.IsEmpty() {
		return object, nil
	}

	// Next steps depend on the type of object we're working with.
	switch obj := object.(type) {

	case Getter:
		return obj.GetPath(path)

	case map[string]interface{}:
		return path.getMapOfInterface(obj)

	case []interface{}:
		return path.getArrayOfInterface(obj)
	}

	// Fall through to here means that we're working with an object we don't immediately recognize.
	// Let's use some reflection :(

	value := reflect.ValueOf(object)

	return path.getReflect(value)
}

// Set tries to return the value of the object at this path.
func (path Path) Set(object interface{}, value interface{}) error {

	switch obj := object.(type) {

	case Setter:
		return obj.SetPath(path, value)

	case map[string]interface{}:
		return path.setMapOfInterface(obj, value)

	case []interface{}:
		return path.setArrayOfInterface(obj, value)

	}

	return path.setReflect(reflect.ValueOf(object), reflect.ValueOf(value))
}

// IsEmpty returns TRUE if this path does not contain any tokens
func (path Path) IsEmpty() bool {
	return len(path) == 0
}

// HasTail returns TRUE if this path has one or more items in its tail.
func (path Path) HasTail() bool {
	return len(path) > 1
}

// Head returns the first token in the path.
func (path Path) Head() string {
	return path[0]
}

// Tail returns a slice of all tokens *after the first token*
func (path Path) Tail() Path {
	return path[1:]
}

// String implements the Stringer interface, and converts the path into a readable string
func (path Path) String() string {
	return strings.Join(path, ".")
}

// Index is useful for vetting array indices.  It attempts to convert the Head() token int
// an integer, and then check that the integer is within the designated array bounds (is greater than zero,
// and less than the maximum value provided to the function).
//
// It returns the array index and an error
func (path Path) Index(maximum int) (int, error) {

	result, err := strconv.Atoi(path.Head())

	if err != nil {
		return 0, derp.Wrap(err, "path.Index", "Index must be an integer", path, maximum)
	}

	if result < 0 {
		return 0, derp.New(500, "path.Index", "Index cannot be negative", path, maximum)
	}

	if (maximum != -1) && (result >= maximum) {
		return 0, derp.New(500, "path.Index", "Index out of bounds", path, maximum)
	}

	// Fall through means that this is a valid array index
	return result, nil
}
