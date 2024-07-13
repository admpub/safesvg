package safesvg

import "fmt"

type useRef struct {
	parent *useRef
	count  uint64
	id     string
}

const maxReferences uint64 = 500

func (u *useRef) Add(n uint64) error {
	u.count += n
	c := u.Count(u.count)
	//println(`~~~~~~~~~~~>`, c)
	if c > maxReferences {
		return fmt.Errorf(`%w: more than %d (>%d)`, ErrTooManyReferences, maxReferences, c)
	}
	return nil
}

func (u *useRef) Count(n uint64) uint64 {
	if u.parent != nil {
		if u.parent.count == 0 {
			return u.parent.Count(n)
		}
		return u.parent.Count(u.parent.count * n)
	}
	return n
}

func (u *useRef) String() string {
	return fmt.Sprintf(`{"id":%q,"count":%d,"parent":%s}`, u.id, u.count, u.parent)
}

type useRefs map[string]*useRef

func (u useRefs) Add(parent *useRef, id string, n uint64) error {
	ref, ok := u[id]
	if !ok {
		ref = &useRef{parent: parent, id: id}
		u[id] = ref
	}
	return ref.Add(n)
}

func (u useRefs) New(parent *useRef, id string) {
	_, ok := u[id]
	if !ok {
		u[id] = &useRef{parent: parent, id: id}
	}
}
