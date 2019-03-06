package providers

import (
	"github.com/stretchr/testify/assert"
	"github.com/tengen-io/server/test"
	"testing"
)

func TestUserProvider_GetUserById(t *testing.T) {
	db := test.MakeDb()
	p := NewUserProvider(db)

	res, err := p.GetUserById("1")
	assert.NoError(t, err)

	assert.Equal(t, "1", res.Id)
	assert.Equal(t, "Test User 1", res.Name)
}