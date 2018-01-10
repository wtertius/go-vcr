package track

import (
	"reflect"
	"time"
)

type trackForYAML struct {
	Args    []interface{}
	Results []interface{}

	IsRecorded bool
	Duration   time.Duration
}

func (track *Track) MarshalYAML() (interface{}, error) {
	results := make([]interface{}, len(track.out))
	for i := range track.out {
		results[i] = track.out[i].Interface()
	}

	tr := trackForYAML{
		Args:    track.args,
		Results: results,

		IsRecorded: track.IsRecorded(),
		Duration:   track.duration,
	}

	return tr, nil
}

func (track *Track) UnmarshalYAML(unmarshal func(interface{}) error) error {
	tr := new(trackForYAML)
	err := unmarshal(tr)
	if err != nil {
		return err
	}

	track.args = tr.Args
	track.out = make([]reflect.Value, 0, len(tr.Results))
	for _, result := range tr.Results {
		track.out = append(track.out, reflect.ValueOf(result))
	}
	track.isRecorded = tr.IsRecorded
	track.duration = tr.Duration

	return nil
}
