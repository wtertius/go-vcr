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
		trRecorded := track.New().Call(emptyFn).With(emptyArgs...).ResultsIn(emptyResults...)
		trRecorded.Record()

		t.Run("Record track to a cassete", func(t *testing.T) {
			t.Run("Only finished track can be recorded to cassete", func(t *testing.T) {
				t.Run("Can't record not finished track to cassete", func(t *testing.T) {
					tr := track.New()

					cas := cassete.New()
					err := cas.Record(tr)
					assert.NotNil(t, err)
					assert.Equal(t, cassete.ErrTrackWasntRecorded, err)
				})
				t.Run("Can record finished track to cassete", func(t *testing.T) {
					cas := cassete.New()
					err := cas.Record(trRecorded)
					assert.Nil(t, err)
				})
			})
			t.Run("Length is correct", func(t *testing.T) {
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
		t.Run("Next track from cassete", func(t *testing.T) {
			t.Run("Nil on not existing track next try", func(t *testing.T) {
				t.Run("Empty cassete", func(t *testing.T) {
					cas := cassete.New()
					tr := track.New()

					trNext := cas.Next(tr.Key())
					assert.Nil(t, trNext)
				})
				t.Run("No track found", func(t *testing.T) {
					cas := cassete.New()
					cas.Record(trRecorded)

					key := track.New().Call(func(string) {}).With("hey girl").Key()

					trNext := cas.Next(key)
					assert.Nil(t, trNext)
				})
				t.Run("Can't play twice what was recorded once", func(t *testing.T) {
					cas := cassete.New()
					cas.Record(trRecorded)

					trNext := cas.Next(trRecorded.Key())
					assert.NotNil(t, trNext)

					trNext = cas.Next(trRecorded.Key())
					assert.Nil(t, trNext)
				})
			})
			t.Run("Next successfully", func(t *testing.T) {
				t.Run("One track", func(t *testing.T) {
					cas := cassete.New()
					cas.Record(trRecorded)

					trNext := cas.Next(trRecorded.Key())
					assert.Equal(t, trRecorded.Key(), trNext.Key())
				})
				t.Run("Two tracks", func(t *testing.T) {
					cas := cassete.New()

					names := []string{"Alice", "Mary"}
					for _, name := range names {
						fn := func() string { return "hey " + name }

						resultStr := ""
						tr := track.New().Call(fn).ResultsIn(&resultStr)
						tr.Record()

						err := cas.Record(tr)
						assert.Nil(t, err)
					}

					for _, name := range names {
						fn := func() string { return "hey girl" }
						tr := track.New().Call(fn)

						trNext := cas.Next(tr.Key())
						assert.NotNil(t, trNext)
						assert.Equal(t, tr.Key(), trNext.Key())

						resultStr := ""
						trNext.Playback(&resultStr)
						assert.Equal(t, "hey "+name, resultStr)
					}
				})
			})
		})
	})
}
