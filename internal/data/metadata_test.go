package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateMetadata(t *testing.T) {
	t.Run("returns empty metadata when no records", func(t *testing.T) {
		result := calculateMetadata(0, 1, 10)

		assert.Equal(t, Metadata{}, result)
	})

	t.Run("calculates correct metadata for single page", func(t *testing.T) {
		result := calculateMetadata(5, 1, 10)

		assert.Equal(t, 1, result.CurrentPage)
		assert.Equal(t, 10, result.PageSize)
		assert.Equal(t, 1, result.FirstPage)
		assert.Equal(t, 1, result.LastPage)
		assert.Equal(t, 5, result.TotalRecords)
	})

	t.Run("calculates correct metadata for multiple pages", func(t *testing.T) {
		result := calculateMetadata(25, 2, 10)

		assert.Equal(t, 2, result.CurrentPage)
		assert.Equal(t, 10, result.PageSize)
		assert.Equal(t, 1, result.FirstPage)
		assert.Equal(t, 3, result.LastPage)
		assert.Equal(t, 25, result.TotalRecords)
	})

	t.Run("calculates last page correctly when records divide evenly", func(t *testing.T) {
		result := calculateMetadata(30, 1, 10)

		assert.Equal(t, 3, result.LastPage)
		assert.Equal(t, 30, result.TotalRecords)
	})

	t.Run("calculates last page correctly when records don't divide evenly", func(t *testing.T) {
		result := calculateMetadata(31, 1, 10)

		assert.Equal(t, 4, result.LastPage)
		assert.Equal(t, 31, result.TotalRecords)
	})

	t.Run("handles edge case with 1 record", func(t *testing.T) {
		result := calculateMetadata(1, 1, 10)

		assert.Equal(t, 1, result.LastPage)
		assert.Equal(t, 1, result.TotalRecords)
	})

	t.Run("handles large page sizes", func(t *testing.T) {
		result := calculateMetadata(100, 1, 50)

		assert.Equal(t, 2, result.LastPage)
		assert.Equal(t, 50, result.PageSize)
	})
}
