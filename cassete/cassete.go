package cassete

import (
	"errors"
	"sync"

	"go-vcr/track"
	"go-vcr/tracklist"
)

var ErrTrackWasntRecorded = errors.New("Not recorded track can't be recorded to cassete")

type TrackMap map[track.Key]*tracklist.TrackList

type Cassete struct {
	id         uint64
	tracks     TrackMap
	isRecorded bool

	mutex sync.RWMutex
}

func New() *Cassete {
	return &Cassete{
		tracks: make(TrackMap),
	}
}

func (c *Cassete) ID() uint64 {
	return c.id
}

func (c *Cassete) Record(tracks ...*track.Track) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, tr := range tracks {
		err := c.record(tr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Cassete) record(tr *track.Track) error {
	if !tr.IsRecorded() {
		return ErrTrackWasntRecorded
	}

	if _, ok := c.tracks[tr.Key()]; !ok {
		c.tracks[tr.Key()] = tracklist.New()
	}

	c.tracks[tr.Key()].Append(tr)

	return nil
}

func (c *Cassete) Length() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	length := 0
	for _, trackList := range c.tracks {
		length += trackList.Length()
	}

	return length
}

func (c *Cassete) GetTrack(key track.Key) *track.Track {
	if trackList, ok := c.tracks[key]; ok {
		return trackList.Next()
	}

	return nil
}

func (c *Cassete) Exec(tr *track.Track) error {
	if c.isRecorded {
		return c.GetTrack(tr.Key()).ResultsAs(tr).Playback()
	}

	err := tr.Record()
	if err != nil {
		return err
	}

	return c.Record(tr)
}
