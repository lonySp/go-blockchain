package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVM(t *testing.T) {
	data := []byte{0x02, 0x0a, 0x02, 0x0a, 0x0b}
	vm := NewVM(data)
	assert.Nil(t, vm.Run())

	assert.Equal(t, byte(4), vm.stack[vm.sp])

}
