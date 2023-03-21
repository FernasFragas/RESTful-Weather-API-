package chat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetLocationFromChat(t *testing.T) {
	loc := GetLocationFromChat("What is the weather in Lisbon")

	assert.Equal(t, "Lisbon", loc)
}
