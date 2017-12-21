package cassete_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go-vcr/cassete"
	"go-vcr/track"
)

var emptyFn = func() {}
var emptyArgs = []interface{}{}
var emptyResults = []interface{}{}

func TestCassete(t *testing.T) {
	t.Run("New cassete isn't nil", func(t *testing.T) {
		cas := cassete.New()
		assert.NotNil(t, cas)
	})
	t.Run("Work with tracks", func(t *testing.T) {
		t.Run("Record track to a cassete", func(t *testing.T) {
			trRecorded := track.New().Call(emptyFn).With(emptyArgs...).ResultsIn(emptyResults...)
			trRecorded.Record()

			t.Run("Can't record not finished track to cassete", func(t *testing.T) {
				tr := track.New()

				cas := cassete.New()
				err := cas.Record(tr)
				assert.NotNil(t, err)
				assert.Equal(t, cassete.ErrTrackWasntRecorded, err)
			})
			t.Run("Can record finished track to cassete", func(t *testing.T) {
				tr := trRecorded

				cas := cassete.New()
				err := cas.Record(tr)
				assert.Nil(t, err)
			})
			t.Run("One track", func(t *testing.T) {
				cas := cassete.New()
				cas.Record(trRecorded)
				assert.Equal(t, 1, cas.Length())
			})
			t.Run("Two tracks with the same key", func(t *testing.T) {
				cas := cassete.New()
				cas.Record(trRecorded)
				cas.Record(trRecorded)
				assert.Equal(t, 2, cas.Length())
			})
			t.Run("Two tracks with different keys", func(t *testing.T) {
				fn := func(string) {}
				tracks := []*track.Track{
					track.New().Call(fn).With("one").ResultsIn(emptyResults...),
					track.New().Call(fn).With("two").ResultsIn(emptyResults...),
				}
				for _, tr := range tracks {
					tr.Record()
				}

				cas := cassete.New()
				cas.Record(tracks...)
				assert.Equal(t, 2, cas.Length())
			})
		})
	})
}
