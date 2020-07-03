package path

import "github.com/benpate/derp"

type testStruct struct {
	Name      string
	Email     string
	Relatives testStructArray
}

func (d testStruct) GetPath(path Path) (interface{}, error) {

	if path.IsEmpty() {
		return d, nil
	}

	switch path.Head() {
	case "name":
		return d.Name, nil
	case "email":
		return d.Email, nil
	case "relatives":
		return d.Relatives.GetPath(path.Tail())
	}

	return nil, derp.New(500, "path.testData", "unsupported")
}

type testStructArray []testStruct

func (d testStructArray) GetPath(path Path) (interface{}, error) {

	if path.IsEmpty() {
		return d, nil
	}

	index, err := path.Index(len(d))

	if err != nil {
		return nil, derp.Wrap(err, "path.testDataArray", "Invalid array index", path)
	}

	return d[index].GetPath(path.Tail())
}

func getTestStruct() testStruct {

	return testStruct{
		Name:  "John Connor",
		Email: "john@connor.mil",
		Relatives: testStructArray{
			{
				Name:  "Sarah Connor",
				Email: "sarah@sky.net",
				Relatives: testStructArray{
					{Name: "John Connor"},
					{Name: "Kyle Reese"},
				},
			},
			{
				Name:  "Kyle Reese",
				Email: "kyle@resistance.mil",
				Relatives: testStructArray{
					{Name: "John Connor"},
					{Name: "Sarah Connor"},
				},
			},
		},
	}
}

func getTestData() map[string]interface{} {

	return map[string]interface{}{
		"name":  "John Connor",
		"email": "john@connor.mil",
		"relatives": map[string]interface{}{
			"mom": "Sarah Connor",
			"dad": "Kyle Reese",
		},
		"enemies": []interface{}{"first terminator", "second terminator", "third terminator"},
	}
}
