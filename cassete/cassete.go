package cassete

import (
	"errors"
	"sync"

	"go-vcr/track"
)

var ErrTrackWasntRecorded = errors.New("Not recorded track can't be recorded to cassete")

type Tracks map[track.Key][]*track.Track

type Cassete struct {
	tracks Tracks
	mutex  sync.RWMutex
}

func New() *Cassete {
	return &Cassete{
		tracks: make(Tracks),
	}
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
		c.tracks[tr.Key()] = make([]*track.Track, 0, 1)
	}

	c.tracks[tr.Key()] = append(c.tracks[tr.Key()], tr)

	return nil
}

func (c *Cassete) Length() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	length := 0
	for _, trackList := range c.tracks {
		length += len(trackList)
	}

	return length
}
