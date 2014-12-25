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
	// making sure that the order of the tags is the same and that this does not rely
	// only on a map - this is a bad and non-deterministic way to do this but fix later
	for i := 0; i < 100; i++ {
		require.Equal(t, "inject: foo tags{tagOne:1 tagTwo:two}", injectError.Error())
	}
	require.Equal(t, "foo", injectError.Type())
	tagOne, ok := injectError.GetTag("tagOne")
	require.True(t, ok)
	require.Equal(t, 1, tagOne)
	tagTwo, ok := injectError.GetTag("tagTwo")
	require.True(t, ok)
	require.Equal(t, "two", tagTwo)
}
