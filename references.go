package safesvg

import (
	"fmt"
	"os"
	"strconv"
)

type useRef struct {
	parent *useRef
	count  uint64
	id     string
}

var maxReferences uint64 = getMaxReferences()

func getMaxReferences() uint64 {
	v := os.Getenv(`SVG_ELEMENT_MAX_REFERENCES`)
	var n uint64
	if len(v) > 0 {
		n, _ = strconv.ParseUint(v, 10, 64)
	}
	if n == 0 {
		n = 50
	}
	return n
}

func (u *useRef) Add(n uint64) error {
	u.count += n
	c := u.calcCount(u.count)
	//println(`~~~~~~~~~~~>`, c)
	if c > maxReferences {
		return fmt.Errorf(`%w(id=%q): more than %d (>%d)`, ErrTooManyReferences, u.id, maxReferences, c)
	}
	return nil
}

func (u *useRef) calcCount(n uint64) uint64 {
	if u.parent != nil {
		if u.parent.count == 0 {
			return u.parent.calcCount(n)
		}
		return u.parent.calcCount(u.parent.count * n)
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
