package reasoner

// Intern maps strings to compact uint32 IDs and back.
// ID 0 is reserved (never assigned), so zero-value uint32 can represent "empty".
type Intern struct {
	toID  map[string]uint32
	toStr []string
}

// NewIntern creates a new interning pool.
func NewIntern() *Intern {
	p := &Intern{
		toID:  make(map[string]uint32),
		toStr: make([]string, 1, 64), // index 0 is reserved ("")
	}
	p.toStr[0] = ""
	return p
}

// ID returns the uint32 ID for s, assigning a new one if not yet interned.
func (p *Intern) ID(s string) uint32 {
	if id, ok := p.toID[s]; ok {
		return id
	}
	id := uint32(len(p.toStr))
	p.toID[s] = id
	p.toStr = append(p.toStr, s)
	return id
}

// Str returns the string for the given ID.
// Panics if id is out of range.
func (p *Intern) Str(id uint32) string {
	return p.toStr[id]
}

// Size returns the number of interned strings (excluding the reserved slot).
func (p *Intern) Size() int {
	return len(p.toStr) - 1
}
