package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceRemove(t *testing.T) {
	assert := assert.New(t)

	t.Run("Remove one item from array", func(t *testing.T) {
		list := []string{"apple", "avocado", "cake", "smartphone"}
		stringToRemove := "smartphone"
		expectedList := []string{"apple", "avocado", "cake"}

		newList, removed := SliceRemove(list, stringToRemove)

		assert.EqualValues(true, removed)
		assert.EqualValues(3, len(newList))
		assert.EqualValues(expectedList, newList)
	})

	t.Run("Remove one item in the middle of a array", func(t *testing.T) {
		list := []string{"apple", "avocado", "cake", "smartphone"}
		stringToRemove := "avocado"
		expectedList := []string{"apple", "cake", "smartphone"}

		newList, removed := SliceRemove(list, stringToRemove)

		assert.EqualValues(true, removed)
		assert.EqualValues(3, len(newList))
		assert.EqualValues(expectedList, newList)
	})

	t.Run("return the same array if not find the string to remove", func(t *testing.T) {
		list := []string{"apple", "avocado", "cake", "smartphone"}
		stringToRemove := "yaiuooo"

		newList, removed := SliceRemove(list, stringToRemove)

		assert.EqualValues(false, removed)
		assert.EqualValues(4, len(newList))
		assert.EqualValues(list, newList)
	})
}
