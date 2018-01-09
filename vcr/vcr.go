package vcr

import (
	"go-vcr/cassete"
	"sync"
)

type CasseteMap map[uint64]*cassete.Cassete

type VCR struct {
	cassetes CasseteMap

	mutex sync.RWMutex
}

func New() *VCR {
	return &VCR{
		cassetes: make(CasseteMap),
	}
}

func (v *VCR) Length() int {
	v.mutex.RLock()
	defer v.mutex.RUnlock()

	return len(v.cassetes)
}

func (v *VCR) Get(id uint64) *cassete.Cassete {
	v.mutex.RLock()
	defer v.mutex.RUnlock()

	return v.cassetes[id]
}

func (v *VCR) Add(cas *cassete.Cassete) {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	v.cassetes[cas.ID()] = cas
}

func (v *VCR) Delete(id uint64) {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	delete(v.cassetes, id)
}
