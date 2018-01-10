package track

import (
	"encoding/json"
	"errors"
	"reflect"
	"time"
)

var ErrNotFunc = errors.New("The first argument must have the type 'func'")
var ErrWrongFuncSignature = errors.New("The args or results signagure doesn't match the function")
var ErrTrackWasntRecorded = errors.New("Can't playback track that wasn't recorded")
var ErrTrackRewritingProhibited = errors.New("Track rewriting is prohibited")

type Track struct {
	fn      typeF
	args    []interface{}
	results []interface{}

	out        []reflect.Value
	isRecorded bool
	duration   time.Duration
}

type Key string

const KeyEmpty = Key("")

type typeF interface{}

func New() *Track {
	return new(Track)
}

func (track *Track) Key() Key {
	key := Key("")
	if track.fn != nil {
		key += Key(reflect.TypeOf(track.fn).String())
	}
	if track.args != nil {
		key += Key(track.argsJSON())
	}

	return key
}

func (track *Track) argsJSON() string {
	argsJSON, _ := json.Marshal(track.args)
	return string(argsJSON)
}

func (track *Track) Call(fn typeF) *Track {
	track.fn = fn
	return track
}

func (track *Track) With(args ...interface{}) *Track {
	track.args = args
	return track
}

func (track *Track) ResultsIn(results ...interface{}) *Track {
	track.results = results
	return track
}

func (track *Track) ResultsAs(trackSource *Track) *Track {
	track.results = trackSource.results
	return track
}

func (track *Track) IsRecorded() bool {
	return track.isRecorded
}

func (track *Track) Playback(results ...interface{}) error {
	if !track.IsRecorded() {
		return ErrTrackWasntRecorded
	}

	if track.fn != nil {
		err := track.checkResults(track.results)
		if err != nil {
			return err
		}
	}

	track.setResults(results)

	return nil
}

func (track *Track) Record() error {
	if track.IsRecorded() {
		return ErrTrackRewritingProhibited
	}

	err := track.CheckFn()
	if err != nil {
		return err
	}

	track.do()

	return nil
}

func (track *Track) setDurationSince(startTime time.Time) {
	track.duration = time.Since(startTime)
}

func (track *Track) do() {
	defer track.setDurationSince(time.Now())

	in := track.getFnIn(track.args)
	track.out = reflect.ValueOf(track.fn).Call(in)
	track.isRecorded = true

	track.setResults(track.results)
}

func (track *Track) setResults(results []interface{}) {
	if results == nil {
		results = track.results
	}

	for i := range track.out {
		reflect.ValueOf(results[i]).Elem().Set(track.out[i])
	}
}

func (track *Track) CheckFn() error {
	v := reflect.ValueOf(track.fn)
	t := v.Type()

	if v.Kind() != reflect.Func {
		return ErrNotFunc
	}
	if t.NumIn() != len(track.args) || t.NumOut() != len(track.results) {
		return ErrWrongFuncSignature
	}

	for i := 0; i < t.NumIn(); i++ {
		if t.In(i).String() != reflect.TypeOf(track.args[i]).String() {
			return ErrWrongFuncSignature
		}
	}

	err := track.checkResults(track.results)
	if err != nil {
		return err
	}

	return nil
}

func (track *Track) checkResults(results []interface{}) error {
	t := reflect.TypeOf(track.fn)

	for i := 0; i < t.NumOut(); i++ {
		if t.Out(i).String() != reflect.TypeOf(results[i]).Elem().String() {
			return ErrWrongFuncSignature
		}
	}

	return nil
}

func (track *Track) getFnIn(args []interface{}) []reflect.Value {
	in := make([]reflect.Value, len(args))

	for i := range args {
		in[i] = reflect.ValueOf(args[i])
	}

	return in
}
