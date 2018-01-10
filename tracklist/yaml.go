package tracklist

import (
	"go-vcr/track"
)

type trackListForYAML struct {
	Tracks []*track.Track
}

func (t *TrackList) MarshalYAML() (interface{}, error) {
	tl := trackListForYAML{
		Tracks: t.tracks,
	}

	return tl, nil
}

func (t *TrackList) UnmarshalYAML(unmarshal func(interface{}) error) error {
	tl := new(trackListForYAML)
	err := unmarshal(tl)
	if err != nil {
		return err
	}

	t.tracks = tl.Tracks

	return nil
}
