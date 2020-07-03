package path

import (
	"testing"

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
		assert.Equal(t, "first terminator", value)
		assert.Nil(t, err)
	}

	{
		value, err := Get(d, "enemies.1")
		assert.Equal(t, "second terminator", value)
		assert.Nil(t, err)
	}

	{
		value, err := Get(d, "enemies.2")
		assert.Equal(t, "third terminator", value)
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
