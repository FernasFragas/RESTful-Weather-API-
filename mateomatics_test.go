package weatherservice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewMateoMaticAPI(t *testing.T) {
	WeatherServiceKeys := LoadEnvKey()

	mateoMaticAPI := NewMateoMaticAPI(WeatherServiceKeys.MateoMaticsAuths.Username, WeatherServiceKeys.MateoMaticsAuths.Password)

	assert.NotEmpty(t, mateoMaticAPI)
	assert.NotEmpty(t, mateoMaticAPI.token)

}
