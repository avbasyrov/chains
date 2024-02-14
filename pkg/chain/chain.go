package chain

import (
	"sort"
	"sync"
)

type Object int
type Relation [2]Object

type Chain struct {
	list []Object
	keys map[Object]bool

	mu sync.RWMutex
}

type Chains struct {
	chains map[Object]*Chain
}

func NewChains(relations []Relation) *Chains {
	chains := make([]*Chain, 0)

	sort.Slice(relations, func(i, j int) bool {
		if int(relations[i][0]) < int(relations[j][0]) {
			return true
		}

		if int(relations[i][0]) == int(relations[j][0]) && int(relations[i][1]) < int(relations[j][1]) {
			return true
		}

		return false
	})

	for _, relation := range relations {
		chainFound := false

		for _, chain := range chains {
			if chain.Exists(relation[0]) || chain.Exists(relation[1]) {
				chainFound = true
				chain.Add(relation[:]...)
				break
			}
		}

		if !chainFound {
			chain := New()
			chains = append(chains, chain)

			chain.Add(relation[:]...)
		}
	}

	// find intersected chains
	merged := make(map[int]bool)

	for idx1, chain1 := range chains {
		if _, ok := merged[idx1]; ok {
			continue
		}

		for idx2, chain2 := range chains {
			if idx1 == idx2 {
				continue // don't merge to itself
			}

			if _, ok := merged[idx2]; ok {
				continue
			}

			if chain1.Intersects(chain2) {
				chain1.Add(chain2.List()...)
				merged[idx2] = true
			}
		}
	}

	result := make(map[Object]*Chain)

	for idx, chain := range chains {
		if _, ok := merged[idx]; ok {
			continue
		}

		for _, obj := range chain.List() {
			result[obj] = chain
		}
	}

	return &Chains{chains: result}
}

func (c *Chains) Lookup(obj Object) []Object {
	if chain, ok := c.chains[obj]; ok {
		return chain.List()
	}

	return nil
}

func New() *Chain {
	return &Chain{keys: make(map[Object]bool)}
}

func (c *Chain) Add(objs ...Object) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, obj := range objs {
		// check that Object not exists yet
		if _, ok := c.keys[obj]; ok {
			continue
		}

		c.keys[obj] = true
		c.list = append(c.list, obj)
	}
}

func (c *Chain) Exists(obj Object) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, ok := c.keys[obj]

	return ok
}

func (c *Chain) List() []Object {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.list
}

func (c *Chain) Intersects(chain *Chain) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, obj := range chain.List() {
		if _, ok := c.keys[obj]; ok {
			return true
		}
	}

	return false
}
