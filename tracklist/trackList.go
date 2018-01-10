package tracklist

import (
	"sync"

	"go-vcr/track"
)

type TrackList struct {
	tracks   []*track.Track
	iterator int

	mutex sync.RWMutex
}

func New() *TrackList {
	return &TrackList{
		tracks: make([]*track.Track, 0, 1),
	}
}

func (t *TrackList) Append(tr *track.Track) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.tracks = append(t.tracks, tr)
}

func (t *TrackList) Length() int {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.length()
}

func (t *TrackList) length() int {
	return len(t.tracks)
}

func (t *TrackList) Next() *track.Track {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.iterator == t.length() {
		return nil
	}

	tr := t.tracks[t.iterator]
	t.iterator++

	return tr
}

func (t *TrackList) ResetIterator() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.iterator = 0
}
