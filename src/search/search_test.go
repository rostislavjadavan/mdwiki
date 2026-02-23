package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHighlightFromPositions(t *testing.T) {
	assert.Equal(t, `<b class="highlight">he</b>llo`, highlightFromPositions("hello", []int{0, 1}))
	assert.Equal(t, `h<b class="highlight">e</b>l<b class="highlight">l</b>o`, highlightFromPositions("hello", []int{1, 3}))
	assert.Equal(t, "hello", highlightFromPositions("hello", []int{}))
	assert.Equal(t, `<b class="highlight">hello</b>`, highlightFromPositions("hello", []int{0, 1, 2, 3, 4}))
}
