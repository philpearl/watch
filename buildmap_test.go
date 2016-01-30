package rebuilder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildMap(t *testing.T) {
	bm := buildMap{}

	bm.add("a", "a.a")
	bm.add("a", "a.b")
	bm.add("a.b", "a.b.c")
	bm.add("b", "a.a")
	bm.add("c", "a.b.c")

	assert.Equal(t, map[string]struct{}{"a.b": struct{}{}, "c": struct{}{}}, bm.rebuild("a.b.c"))
	assert.Equal(t, map[string]struct{}{"a": struct{}{}}, bm.rebuild("a.b"))
	assert.Equal(t, map[string]struct{}(nil), bm.rebuild("a"))
	assert.Equal(t, map[string]struct{}{"a": struct{}{}, "b": struct{}{}}, bm.rebuild("a.a"))
}
