package id

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateID(t *testing.T) {
	id, err := GenerateID()
	assert.NoError(t, err)
	assert.NotEmpty(t, id)
	assert.Len(t, id, 10)
}
