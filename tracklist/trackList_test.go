package tracklist_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go-vcr/track"
	"go-vcr/tracklist"
)

var emptyFn = func() {}
var emptyArgs = []interface{}{}
var emptyResults = []interface{}{}

func TestTrackList(t *testing.T) {
	t.Run("New track list isn't nil", func(t *testing.T) {
		tl := tracklist.New()
		assert.NotNil(t, tl)
	})
	t.Run("Work with tracks", func(t *testing.T) {
		trRecorded := track.New().Call(emptyFn).With(emptyArgs...).ResultsIn(emptyResults...)
		trRecorded.Record()

		t.Run("Append track to a tracklist", func(t *testing.T) {
			t.Run("Length is correct", func(t *testing.T) {
				t.Run("One track", func(t *testing.T) {
					tl := tracklist.New()
					tl.Append(trRecorded)
					assert.Equal(t, 1, tl.Length())
				})
				t.Run("Two tracks", func(t *testing.T) {
					tl := tracklist.New()
					tl.Append(trRecorded)
					tl.Append(trRecorded)
					assert.Equal(t, 2, tl.Length())
				})
			})
		})
		t.Run("Next track from track list", func(t *testing.T) {
			t.Run("Nil on not existing track next try", func(t *testing.T) {
				t.Run("Empty tracklist", func(t *testing.T) {
					tl := tracklist.New()

					trNext := tl.Next()
					assert.Nil(t, trNext)
				})
				t.Run("Can't get twice what was appended once", func(t *testing.T) {
					tl := tracklist.New()
					tl.Append(trRecorded)

					trNext := tl.Next()
					assert.NotNil(t, trNext)

					trNext = tl.Next()
					assert.Nil(t, trNext)
				})
			})
			t.Run("Next successfully", func(t *testing.T) {
				t.Run("One track", func(t *testing.T) {
					tl := tracklist.New()
					tl.Append(trRecorded)

					trNext := tl.Next()
					assert.Equal(t, trRecorded.Key(), trNext.Key())
				})
				t.Run("Two tracks", func(t *testing.T) {
					tl := tracklist.New()

					names := []string{"Alice", "Mary"}
					for _, name := range names {
						fn := func() string { return "hey " + name }

						resultStr := ""
						tr := track.New().Call(fn).ResultsIn(&resultStr)
						tr.Record()

						tl.Append(tr)
					}

					for _, name := range names {
						fn := func() string { return "hey girl" }
						tr := track.New().Call(fn)

						trNext := tl.Next()
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
