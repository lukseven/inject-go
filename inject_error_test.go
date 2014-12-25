package inject

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInjectErrorWithoutTags(t *testing.T) {
	errorBuilder := newErrorBuilder("foo")
	injectError := errorBuilder.build()
	require.Equal(t, "inject: foo", injectError.Error())
	require.Equal(t, "foo", injectError.Type())
	_, ok := injectError.GetTag("tagOne")
	require.False(t, ok)
}

func TestInjectErrorWithTags(t *testing.T) {
	errorBuilder := newErrorBuilder("foo")
	errorBuilder = errorBuilder.addTag("tagOne", 1)
	errorBuilder = errorBuilder.addTag("tagTwo", "two")
	injectError := errorBuilder.build()
	require.Contains(t, injectError.Error(), "inject: foo")
	require.Contains(t, injectError.Error(), "tagOne:1")
	require.Contains(t, injectError.Error(), "tagTwo:two")
	require.Equal(t, "foo", injectError.Type())
	tagOne, ok := injectError.GetTag("tagOne")
	require.True(t, ok)
	require.Equal(t, 1, tagOne)
	tagTwo, ok := injectError.GetTag("tagTwo")
	require.True(t, ok)
	require.Equal(t, "two", tagTwo)
}
