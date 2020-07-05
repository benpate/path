package path

import (
	"testing"

	"github.com/benpate/derp"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {

	{
		p := New("one.two.three")

		assert.Equal(t, 3, len(p))
		assert.Equal(t, "one", p[0])
		assert.Equal(t, "two", p[1])
		assert.Equal(t, "three", p[2])
	}

	{
		p := New("")
		assert.True(t, p.IsEmpty())
	}
}

func TestProperties(t *testing.T) {

	d := getTestData()

	{
		value, err := Get(d, "name")
		assert.Equal(t, "John Connor", value)
		assert.Nil(t, err)
	}

	{
		value, err := Get(d, "email")
		assert.Equal(t, "john@connor.mil", value)
		assert.Nil(t, err)
	}

	{
		value, err := Get(d, "missing property")
		assert.Nil(t, value)
		assert.NotNil(t, err)
	}
}

func TestSubProperties(t *testing.T) {

	d := getTestData()

	{
		value, err := Get(d, "relatives.mom")
		assert.Equal(t, "Sarah Connor", value)
		assert.Nil(t, err)
	}

	{
		value, err := Get(d, "relatives.dad")
		assert.Equal(t, "Kyle Reese", value)
		assert.Nil(t, err)
	}

	{
		value, err := Get(d, "relatives.sister")
		assert.Nil(t, value)
		assert.NotNil(t, err)
	}
}

func TestArrays(t *testing.T) {

	d := getTestData()

	{
		value, err := Get(d, "enemies.0")
		assert.Equal(t, "T-1000", value)
		assert.Nil(t, err)
	}

	{
		value, err := Get(d, "enemies.1")
		assert.Equal(t, "T-3000", value)
		assert.Nil(t, err)
	}

	{
		value, err := Get(d, "enemies.2")
		assert.Equal(t, "T-5000", value)
		assert.Nil(t, err)
	}

	{
		value, err := Get(d, "enemies.-1")
		assert.Nil(t, value)
		assert.NotNil(t, err)
	}

	{
		value, err := Get(d, "enemies.3")
		assert.Nil(t, value)
		assert.NotNil(t, err)
	}

	{
		value, err := Get(d, "enemies.100000")
		assert.Nil(t, value)
		assert.NotNil(t, err)
	}

	{
		value, err := Get(d, "enemies.fred")
		assert.Nil(t, value)
		assert.NotNil(t, err)
	}
}

func TestError(t *testing.T) {

	{
		value, err := Get("unsupported data", "property")
		assert.Nil(t, value)
		assert.NotNil(t, err)
	}

	{
		value, err := Get("string at the end of a path", "")
		assert.Equal(t, "string at the end of a path", value)
		assert.Nil(t, err)
	}
}

func TestGetter(t *testing.T) {

	d := getTestStruct()

	{
		value, err := Get(d, "name")
		assert.Equal(t, "John Connor", value)
		assert.Nil(t, err)
	}

	{
		value, err := Get(d, "email")
		assert.Equal(t, "john@connor.mil", value)
		assert.Nil(t, err)
	}

	{
		value, err := Get(d, "relatives.0.name")
		assert.Equal(t, "Sarah Connor", value)
		assert.Nil(t, err)
	}

	{
		value, err := Get(d, "relatives.1.relatives.1.name")
		assert.Equal(t, "Sarah Connor", value)
		assert.Nil(t, err)
	}

	{
		value, err := Get(d, "missing-property")
		assert.Nil(t, value)
		assert.NotNil(t, err)
	}
}

func TestSet(t *testing.T) {

	d := getTestData()

	// Not implemented yet, so this should just error
	assert.NotNil(t, Set(d, "anywhere.doesnt.matter", "1"))
}

func TestString(t *testing.T) {

	p := New("this.is.a.path.7")

	assert.Equal(t, "this.is.a.path.7", p.String())
}

func TestIndex(t *testing.T) {

	// Test valid index value
	{
		p := New("7")
		index, error := p.Index(10)
		assert.Equal(t, 7, index)
		assert.Nil(t, error)
	}

	// Test no maximum value
	{
		p := New("7")
		index, error := p.Index(-1)
		assert.Equal(t, 7, index)
		assert.Nil(t, error)
	}

	// Test negative index
	{
		p := New("-1")
		index, error := p.Index(5)
		assert.Equal(t, 0, index)
		assert.NotNil(t, error)
	}

	// Test overflow
	{
		p := New("7")
		index, error := p.Index(5)
		assert.Equal(t, 0, index)
		assert.NotNil(t, error)
	}

	// Test non-integer
	{
		p := New("non-integer")
		index, error := p.Index(5)
		assert.Equal(t, 0, index)
		assert.NotNil(t, error)
	}

}

/////////////////////////////////
// SUPPORT FUNCS
/////////////////////////////////

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
		"enemies": []interface{}{"T-1000", "T-3000", "T-5000"},
	}
}
