package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddition(t *testing.T) {
	result := 2 + 2
	assert.Equal(t, 4, result, "2+2 = 4")
}