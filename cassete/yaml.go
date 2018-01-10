package cassete

type casseteForYAML struct {
	ID     uint64
	Tracks TrackMap
}

func (c *Cassete) MarshalYAML() (interface{}, error) {
	cas := casseteForYAML{
		ID:     c.id,
		Tracks: c.tracks,
	}

	return cas, nil
}

func (c *Cassete) UnmarshalYAML(unmarshal func(interface{}) error) error {
	cas := new(casseteForYAML)
	err := unmarshal(cas)
	if err != nil {
		return err
	}

	c.id = cas.ID
	c.tracks = cas.Tracks
	c.isRecorded = true

	return nil
}
