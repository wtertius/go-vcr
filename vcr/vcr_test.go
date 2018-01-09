package vcr_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go-vcr/cassete"
	"go-vcr/vcr"
)

func TestTrackList(t *testing.T) {
	t.Run("New VCR isn't nil", func(t *testing.T) {
		v := vcr.New()
		assert.NotNil(t, v)
	})
	t.Run("New VCR don't have any cassetes", func(t *testing.T) {
		v := vcr.New()
		assert.Equal(t, 0, v.Length())
	})
	t.Run("Can add and get back the added cassete", func(t *testing.T) {
		v := vcr.New()
		cas := cassete.New()

		v.Add(cas)
		casGot := v.Get(cas.ID())

		assert.Equal(t, casGot, cas)
		assert.Equal(t, 1, v.Length())
	})
	t.Run("Can delete the added cassete", func(t *testing.T) {
		v := vcr.New()
		cas := cassete.New()

		v.Add(cas)
		v.Delete(cas.ID())
		casGot := v.Get(cas.ID())

		assert.Nil(t, casGot)
		assert.Equal(t, 0, v.Length())
	})
}
