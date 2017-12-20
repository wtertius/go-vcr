package track_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go-vcr/track"
)

var emptyArgs = []interface{}{}
var emptyResults = []interface{}{}

func TestFnIsCalled(t *testing.T) {
	t.Run("New track isn't nil", func(t *testing.T) {
		tr := track.New()
		assert.NotNil(t, tr)
	})
	t.Run("Wrong func arg", func(t *testing.T) {
		t.Run("Not func", func(t *testing.T) {
			fn := "func must be here"

			tr := track.New().Call(fn).With(emptyArgs...).ResultsIn(emptyResults...)

			err := tr.Record()
			assert.NotNil(t, err)
			assert.Equal(t, track.ErrNotFunc, err)
		})
		t.Run("Wrong number of args", func(t *testing.T) {
			fn := func(string) {}

			tr := track.New().Call(fn).With(emptyArgs...).ResultsIn(emptyResults...)

			err := tr.Record()
			assert.NotNil(t, err)
			assert.Equal(t, track.ErrWrongFuncSignature, err)
		})
		t.Run("Wrong number of results", func(t *testing.T) {
			fn := func() string { return "" }

			tr := track.New().Call(fn).With(emptyArgs...).ResultsIn(emptyResults...)

			err := tr.Record()
			assert.NotNil(t, err)
			assert.Equal(t, track.ErrWrongFuncSignature, err)
		})
		t.Run("Wrong arg type", func(t *testing.T) {
			fn := func(string) {}

			tr := track.New().Call(fn).With(0).ResultsIn(emptyResults...)

			err := tr.Record()
			assert.NotNil(t, err)
			assert.Equal(t, track.ErrWrongFuncSignature, err)
		})
		t.Run("Wrong result type", func(t *testing.T) {
			fn := func(string) {}

			tr := track.New().Call(fn).With(emptyArgs...).ResultsIn(0)

			err := tr.Record()
			assert.NotNil(t, err)
			assert.Equal(t, track.ErrWrongFuncSignature, err)
		})
	})
	t.Run("Call succeeds", func(t *testing.T) {
		t.Run("Empty func", func(t *testing.T) {
			fn := func() {}

			tr := track.New().Call(fn)

			err := tr.Record()
			assert.Nil(t, err)
		})
		t.Run("One argument - one result value", func(t *testing.T) {
			argStr := "Hello world"
			resultStr := ""

			fn := func(argStr string) string { return argStr }

			tr := track.New().Call(fn).With(argStr).ResultsIn(&resultStr)

			err := tr.Record()
			assert.Nil(t, err)
			assert.Equal(t, argStr, resultStr)
		})
		t.Run("Two arguments - two result values", func(t *testing.T) {
			argStr := "Hello world"
			argNumber := 5

			resultStr := ""
			resultNumber := 0

			fn := func(argStr string, argNumber int) (string, int) { return argStr, argNumber }

			tr := track.New().Call(fn).With(argStr, argNumber).ResultsIn(&resultStr, &resultNumber)

			err := tr.Record()
			assert.Nil(t, err)
			assert.Equal(t, argStr, resultStr)
			assert.Equal(t, argNumber, resultNumber)
		})
	})
}

func TestTrackPlayback(t *testing.T) {
	t.Run("Record succeeds", func(t *testing.T) {
		t.Run("Playback not recorded track", func(t *testing.T) {
			fn := func() {}

			tr := track.New().Call(fn)

			err := tr.Playback()

			assert.NotNil(t, err)
			assert.Equal(t, track.ErrTrackWasntRecorded, err)
		})
		t.Run("One argument - one result value", func(t *testing.T) {
			argStr := "Hello world"
			resultStr := ""

			fn := func(argStr string) string { return argStr }

			tr := track.New().Call(fn).With(argStr).ResultsIn(&resultStr)

			tr.Record()

			resultStr = ""
			err := tr.Playback()

			assert.Nil(t, err)
			assert.Equal(t, argStr, resultStr)
		})
	})
}

func TestTrackRecord(t *testing.T) {
	t.Run("Rewriting is prohibited", func(t *testing.T) {
		fn := func() {}

		tr := track.New().Call(fn)

		err := tr.Record()
		assert.Nil(t, err)

		err = tr.Record()
		assert.NotNil(t, err)
		assert.Equal(t, track.ErrTrackRewritingProhibited, err)
	})
}
