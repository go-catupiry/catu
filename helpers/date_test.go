package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractYearFromText(t *testing.T) {
	t.Run("Should extract year from fileName", func(t *testing.T) {
		mockText := "dfp_cia_aberta_2011.zip"
		result := ExtractYearFromText(mockText)
		assert.EqualValues(t, result, "2011")
	})

	t.Run("Should return a empty string without any year on text", func(t *testing.T) {
		mockText := "dfp_cia_aberta.zip"
		result := ExtractYearFromText(mockText)
		assert.EqualValues(t, result, "")
	})

	t.Run("Should return a empty string without any valid year on text", func(t *testing.T) {
		mockText := "dfp_cia_aberta_22_22.zip"
		result := ExtractYearFromText(mockText)
		assert.EqualValues(t, result, "")
	})
}
