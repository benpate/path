package path

import (
	"reflect"

	"github.com/benpate/derp"
)

func (path Path) setMapOfInterface(object map[string]interface{}, value interface{}) error {

	if path.HasTail() == false {
		object[path.Head()] = value
		return nil
	}

	return derp.New(500, "path.Set", "Unimplemented")
}

func (path Path) setArrayOfInterface(object []interface{}, value interface{}) error {

	return derp.New(500, "path.Set", "Unimplemented")

	/*
		index, err := path.Index(-1)

		if err != nil {
			return derp.Wrap(err, "path.Set", "Invalid Array Index", path, object, value)
		}

		object[index] = value
		return nil
	*/
}

func (path Path) setReflect(object reflect.Value, value reflect.Value) error {

	return derp.New(500, "path.Set", "Unimplemented")

	/*
		if object.CanSet() == false {
			return derp.New(500, "path.SetReflect", "Cannot set value", object)
		}

		if value.Type() == object.Type() {
			object.Set(value)
			return nil
		}

		return derp.New(500, "path.SetReflect", "Inompatable types", object, value)
	*/
}
