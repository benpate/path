package path

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {

	d := getTestData()

	// Not implemented yet, so this should just error
	assert.NotNil(t, Set(d, "anywhere.doesnt.matter", "1"))
}
