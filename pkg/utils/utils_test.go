package utils_test

import (
	"main/pkg/utils"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func StringOfRandomLength(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func TestMap(t *testing.T) {
	t.Parallel()

	source := []int{1, 2, 3}

	destination := utils.Map(source, strconv.Itoa)

	require.Equal(t, []string{"1", "2", "3"}, destination)
}

func TestContains(t *testing.T) {
	t.Parallel()

	source := []int{1, 2, 3}

	require.True(t, utils.Contains(source, 2))
	require.False(t, utils.Contains(source, 4))
}

func TestRemoveFirstSlash(t *testing.T) {
	t.Parallel()

	require.Equal(t, "source", utils.RemoveFirstSlash("/source"))
	require.Equal(t, "source", utils.RemoveFirstSlash("source"))
	require.Equal(t, "", utils.RemoveFirstSlash(""))
}

func TestRemoveStripTrailingDigits(t *testing.T) {
	t.Parallel()

	require.Equal(t, "123.456", utils.StripTrailingDigits("123.456789", 3))
	require.Equal(t, "123.4", utils.StripTrailingDigits("123.4", 3))
	require.Equal(t, "123", utils.StripTrailingDigits("123", 3))
	require.Equal(t, "123", utils.StripTrailingDigits("123.456", 0))
}

func TestBoolToFloat64(t *testing.T) {
	t.Parallel()

	require.InDelta(t, float64(1), utils.BoolToFloat64(true), 0.01)
	require.InDelta(t, float64(0), utils.BoolToFloat64(false), 0.01)
}

func TestSplitStringIntoChunksLessThanOneChunk(t *testing.T) {
	t.Parallel()

	str := StringOfRandomLength(10)
	chunks := utils.SplitStringIntoChunks(str, 20)
	assert.Len(t, chunks, 1, "There should be 1 chunk!")
}

func TestSplitStringIntoChunksExactlyOneChunk(t *testing.T) {
	t.Parallel()

	str := StringOfRandomLength(10)
	chunks := utils.SplitStringIntoChunks(str, 10)

	assert.Len(t, chunks, 1, "There should be 1 chunk!")
}

func TestSplitStringIntoChunksMoreChunks(t *testing.T) {
	t.Parallel()

	str := "aaaa\nbbbb\ncccc\ndddd\neeeee\n"
	chunks := utils.SplitStringIntoChunks(str, 10)
	assert.Len(t, chunks, 3, "There should be 3 chunks!")
}
