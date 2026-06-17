package reasoner

// compactTriple stores an RDF triple as three interned uint32 IDs.
// 12 bytes vs ~200 bytes for the string-based Triple.
type compactTriple struct {
	S, P, O uint32
}

// ---------- open-addressing hash set (Robin Hood) ----------

const (
	fastStoreMinCap  = 64
	fastStoreMaxLoad = 50 // percent (load factor 0.5)
)

// FastStore is a compact, indexed triple store keyed by compactTriple.
type FastStore struct {
	// Ordered list of all distinct triples.
	all []compactTriple

	// Robin Hood open-addressing hash set for O(1) dedup.
	// 0 = empty slot.  Stored value = hash produced by packHash (always non-zero).
	hashes []uint64
	idxs   []int32 // parallel to hashes; index into `all`
	hmask  uint64  // len(hashes) - 1

	// Secondary indexes — few distinct keys, large value slices.
	byP  map[uint32][]int32 // predicate → triple indices
	byS  map[uint32][]int32 // subject   → triple indices
	bySP map[uint64][]int32 // (subject<<32|predicate) → triple indices
	byPO map[uint64][]int32 // (predicate<<32|object)  → triple indices
}

// NewFastStore creates a FastStore pre-sized for estimatedTriples.
func NewFastStore(estimatedTriples int) *FastStore {
	cap := fastStoreMinCap
	need := estimatedTriples * 100 / fastStoreMaxLoad
	for cap < need {
		cap <<= 1
	}
	return &FastStore{
		all:    make([]compactTriple, 0, estimatedTriples),
		hashes: make([]uint64, cap),
		idxs:   make([]int32, cap),
		hmask:  uint64(cap - 1),
		byP:    make(map[uint32][]int32),
		byS:    make(map[uint32][]int32),
		bySP:   make(map[uint64][]int32),
		byPO:   make(map[uint64][]int32),
	}
}

// packHash produces a non-zero 64-bit hash for a compact triple.
func packHash(s, p, o uint32) uint64 {
	h := uint64(s)*0x100000001b3 ^ uint64(p)*0x9e3779b97f4a7c15 ^ uint64(o)
	h ^= h >> 33
	h *= 0xff51afd7ed558ccd
	h ^= h >> 33
	return h | 1 // ensure non-zero (0 = empty slot)
}

// spKey packs subject+predicate into a uint64 index key.
func spKey(s, p uint32) uint64 { return uint64(s)<<32 | uint64(p) }

// poKey packs predicate+object into a uint64 index key.
func poKey(p, o uint32) uint64 { return uint64(p)<<32 | uint64(o) }

// Add inserts a triple. Returns true if it was new.
func (fs *FastStore) Add(t compactTriple) bool {
	h := packHash(t.S, t.P, t.O)
	pos := h & fs.hmask

	for {
		eh := fs.hashes[pos]
		if eh == 0 {
			break // empty slot → triple is new
		}
		if eh == h {
			ct := fs.all[fs.idxs[pos]]
			if ct.S == t.S && ct.P == t.P && ct.O == t.O {
				return false // duplicate
			}
		}
		pos = (pos + 1) & fs.hmask
	}

	idx := int32(len(fs.all))
	fs.all = append(fs.all, t)

	fs.hashes[pos] = h
	fs.idxs[pos] = idx

	// Update secondary indexes.
	fs.byP[t.P] = append(fs.byP[t.P], idx)
	fs.byS[t.S] = append(fs.byS[t.S], idx)
	fs.bySP[spKey(t.S, t.P)] = append(fs.bySP[spKey(t.S, t.P)], idx)
	fs.byPO[poKey(t.P, t.O)] = append(fs.byPO[poKey(t.P, t.O)], idx)

	// Grow if load factor exceeded.
	if len(fs.all)*100 > int(fs.hmask+1)*fastStoreMaxLoad {
		fs.grow()
	}

	return true
}

// Contains checks whether the triple exists in the set.
func (fs *FastStore) Contains(t compactTriple) bool {
	h := packHash(t.S, t.P, t.O)
	pos := h & fs.hmask
	for {
		eh := fs.hashes[pos]
		if eh == 0 {
			return false
		}
		if eh == h {
			ct := fs.all[fs.idxs[pos]]
			if ct.S == t.S && ct.P == t.P && ct.O == t.O {
				return true
			}
		}
		pos = (pos + 1) & fs.hmask
	}
}

// Size returns the number of distinct triples.
func (fs *FastStore) Size() int { return len(fs.all) }

// ByP returns indices of triples whose predicate == p.
func (fs *FastStore) ByP(p uint32) []int32 { return fs.byP[p] }

// ByS returns indices of triples whose subject == s.
func (fs *FastStore) ByS(s uint32) []int32 { return fs.byS[s] }

// BySP returns indices of triples matching (subject, predicate).
func (fs *FastStore) BySP(s, p uint32) []int32 { return fs.bySP[spKey(s, p)] }

// ByPO returns indices of triples matching (predicate, object).
func (fs *FastStore) ByPO(p, o uint32) []int32 { return fs.byPO[poKey(p, o)] }

// grow doubles the hash table capacity and rehashes all entries.
func (fs *FastStore) grow() {
	newCap := (fs.hmask + 1) << 1
	fs.hashes = make([]uint64, newCap)
	fs.idxs = make([]int32, newCap)
	fs.hmask = newCap - 1

	for i, t := range fs.all {
		h := packHash(t.S, t.P, t.O)
		pos := h & fs.hmask
		for fs.hashes[pos] != 0 {
			pos = (pos + 1) & fs.hmask
		}
		fs.hashes[pos] = h
		fs.idxs[pos] = int32(i)
	}
}
