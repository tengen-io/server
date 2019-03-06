package providers

import (
	"github.com/stretchr/testify/assert"
	"github.com/tengen-io/server/test"
	"testing"
)

func TestIdentityProvider_GetIdentityById(t *testing.T) {
	db := test.MakeDb()
	p := NewIdentityProvider(db, 4)
	res, err := p.GetIdentityById(1)
	assert.NoError(t, err)

	assert.Equal(t, res.Email, "test1@tengen.io")
	assert.Equal(t, res.Name, "Test User 1")
}

func TestIdentityProvider_CreateIdentity(t *testing.T) {
	db := test.MakeDb()
	p := NewIdentityProvider(db, 4)
	res, err := p.CreateIdentity("test-createidentity@tengen.io", "hunter2", "Test User CreateIdentity")
	assert.NoError(t, err)

	assert.Equal(t, "test-createidentity@tengen.io", res.Email)
	assert.Equal(t, "Test User CreateIdentity", res.Name)
}
