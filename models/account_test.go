package models

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAccount(t *testing.T) {
	acc, err := NewAccount("a", "b", "adasda@mail.com", "123")
	assert.Nil(t, err)

	fmt.Printf("%+v\n", acc)
}
