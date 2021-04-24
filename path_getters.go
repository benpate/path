package path

import (
	"github.com/benpate/derp"
)

/////////////////////////////
// Interface-based getters
/////////////////////////////

func getSliceOfString(path Path, value []string) (string, error) {

	index, err := path.Index(len(value))

	if err != nil {
		return "", err
	}

	if path.HasTail() {
		return "", derp.New(500, "path.Path.getSliceOfString", "Invalid path", path)
	}

	return value[index], nil
}

func getSliceOfInt(path Path, value []int) (int, error) {

	index, err := path.Index(len(value))

	if err != nil {
		return 0, err
	}

	if path.HasTail() {
		return 0, derp.New(500, "path.Path.getSliceOfString", "Invalid path", path)
	}

	return value[index], nil
}

func getSliceOfInterface(path Path, value []interface{}) (interface{}, error) {

	index, err := path.Index(len(value))

	if err != nil {
		return nil, err
	}

	return path.Tail().Get(value[index])
}

func getSliceOfGetter(path Path, value []Getter) (interface{}, error) {

	index, err := path.Index(len(value))

	if err != nil {
		return nil, err
	}

	return path.Tail().Get(value[index])
}

func getMapOfInterface(path Path, value map[string]interface{}) (interface{}, error) {

	key := path.Head()

	field, ok := value[key]

	if !ok {
		return nil, derp.New(500, "path.Path.Get", "Map entry does not exist", path, value)
	}

	return path.Tail().Get(field)
}
