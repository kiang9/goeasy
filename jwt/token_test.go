package token

import (
	"testing"

	"github.com/kiang9/goeasy/assert"
)

func TestParse(t *testing.T) {
	Init(&Config{Secret: "secret"})

	payload := M{"name": "tom"}
	token, err := GenerateToken(payload)
	assert.NoError(t, err)

	p, err := Parse(token, conf.Secret)
	assert.NoError(t, err)

	assert.Equal(t, payload, p)
}
