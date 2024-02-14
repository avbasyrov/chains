package chain_test

import (
	"math/rand"
	"sort"
	"testing"

	"chains/pkg/chain"

	"github.com/stretchr/testify/assert"
)

func TestChains(t *testing.T) {
	/*
		2
		|
		3 — 4 — 5
		|
		1

		==========

		6 — 7 — 8
	*/

	relations := []chain.Relation{
		{2, 3},
		{3, 4},
		{5, 4},
		{6, 7},
		{1, 3},
		{8, 7},
	}

	chains := chain.NewChains(relations)

	t.Run("1-2-3-4-5", func(t *testing.T) {
		for i := 1; i <= 5; i++ {
			result := chains.Lookup(chain.Object(i))
			sort.Slice(result, func(i, j int) bool {
				return int(result[i]) < int(result[j])
			})

			assert.Equal(t, []chain.Object{1, 2, 3, 4, 5}, result)
		}
	})

	t.Run("6-7-8", func(t *testing.T) {
		for i := 6; i <= 8; i++ {
			result := chains.Lookup(chain.Object(i))
			sort.Slice(result, func(i, j int) bool {
				return int(result[i]) < int(result[j])
			})

			assert.Equal(t, []chain.Object{6, 7, 8}, result)
		}
	})

	t.Run("not exists", func(t *testing.T) {
		assert.Equal(t, []chain.Object(nil), chains.Lookup(0))
		assert.Equal(t, []chain.Object(nil), chains.Lookup(9))
		assert.Equal(t, []chain.Object(nil), chains.Lookup(10))
		assert.Equal(t, []chain.Object(nil), chains.Lookup(-1))
	})
}

const count = 50000
const averageChainLen = 1000

func BenchmarkNewChains(b *testing.B) {
	relations := prepareRelations(b)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		chain.NewChains(relations)
	}
}

func BenchmarkLookup(b *testing.B) {
	chains := chain.NewChains(prepareRelations(b))
	lookup := make([]chain.Object, 0, 100)

	for i := 0; i < b.N; i++ {
		lookup = append(lookup, chain.Object(rand.Intn(count)))
	}

	b.ResetTimer()

	for _, obj := range lookup {
		chains.Lookup(obj)
	}
}

func prepareRelations(b *testing.B) []chain.Relation {
	b.Helper()

	relations := make([]chain.Relation, 0, count)

	offset := 0
	for i := 0; i < count; i++ {
		a := i + offset
		b := a + 1

		relations = append(relations, chain.Relation{chain.Object(a), chain.Object(b)})

		if rand.Intn(averageChainLen) == 0 {
			offset++
		}
	}

	rand.Shuffle(len(relations), func(i, j int) {
		relations[i], relations[j] = relations[j], relations[i]
	})

	return relations
}
