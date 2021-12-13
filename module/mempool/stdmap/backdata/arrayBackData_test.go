package backdata

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/module/mempool/stdmap/backdata/arraylinkedlist"
	"github.com/onflow/flow-go/utils/unittest"
)

func TestArrayBackData_SingleBucket(t *testing.T) {
	limit := 10

	bd := NewArrayBackData(uint32(limit), 1, arraylinkedlist.LRUEjection)

	entities := unittest.EntityListFixture(uint(limit))

	// adds all entities to backdata
	testAddEntities(t, bd, entities)

	// sanity checks
	for i := range entities {
		// since we are below limit, elements should be added sequentially at bucket 0.
		// first added element has a key index of 1, since 0 means unused key index in implementation.
		require.Equal(t, bd.buckets[0][i].keyIndex, uint64(i+1))
		// also, since we have not yet over-limited, entities are received valueIndex in the same order they
		// are added.
		require.Equal(t, bd.buckets[0][i].valueIndex, uint32(i))
		_, _, owner := bd.entities.Get(uint32(i))
		require.Equal(t, owner, uint64(i))
	}

	// getting inserted elements
	testRetrievableFrom(t, bd, entities, 0)
}

func TestArrayBackData_Adjust(t *testing.T) {
	limit := 100_000

	bd := NewArrayBackData(uint32(limit), 8, arraylinkedlist.LRUEjection)

	entities := unittest.EntityListFixture(uint(limit))

	// adds all entities to backdata
	testAddEntities(t, bd, entities)

}

// TestArrayBackData_WriteHeavy evaluates correctness of backdata under the writing and retrieving
// a heavy load of entities up to its limit.
func TestArrayBackData_WriteHeavy(t *testing.T) {
	limit := 100_000

	bd := NewArrayBackData(uint32(limit), 8, arraylinkedlist.LRUEjection)

	entities := unittest.EntityListFixture(uint(limit))

	// adds all entities to backdata
	testAddEntities(t, bd, entities)

	// retrieves all entities from backdata
	testRetrievableFrom(t, bd, entities, 0)
}

func TestArrayBackData_LRU_Ejection(t *testing.T) {
	// mempool has the limit of 100K, but we put 1M
	// (10 time more than its capacity)
	limit := 100_000
	items := uint(1_000_000)

	bd := NewArrayBackData(uint32(limit), 8, arraylinkedlist.LRUEjection)

	entities := unittest.EntityListFixture(items)

	// adds all entities to backdata
	testAddEntities(t, bd, entities)

	// only last 100K items must be retrievable, and
	// the rest must be ejected.
	testRetrievableFrom(t, bd, entities, 900_000)
}

func TestArrayBackData_Random_Ejection(t *testing.T) {
	// mempool has the limit of 100K, but we put 1M
	// (10 time more than its capacity)
	limit := 100_000
	items := uint(1_000_000)

	bd := NewArrayBackData(uint32(limit), 8, arraylinkedlist.RandomEjection)

	entities := unittest.EntityListFixture(items)

	// adds all entities to backdata
	testAddEntities(t, bd, entities)

	// only 100K (random) items must be retrievable, as the rest
	// are randomly ejected to make room.
	testRetrievableCount(t, bd, entities, 100_000)
}

func TestArrayBackData_AddDuplicate(t *testing.T) {
	limit := 100

	bd := NewArrayBackData(uint32(limit), 8, arraylinkedlist.LRUEjection)

	entities := unittest.EntityListFixture(uint(limit))

	// adds all entities to backdata
	testAddEntities(t, bd, entities)

	// adding duplicate entity should fail
	for _, entity := range entities {
		require.False(t, bd.Add(entity.ID(), entity))
	}

	// still all mut be retrievable from backdata
	testRetrievableFrom(t, bd, entities, 0)
}

// TestArrayBackData_All_BelowLimit checks correctness of All method when mempool is not full yet.
func TestArrayBackData_All(t *testing.T) {
	tt := []struct {
		limit        uint32
		items        uint32
		ejectionMode arraylinkedlist.EjectionMode
	}{
		{ // mempool has the limit of 100K, but we put 10K
			limit:        100_000,
			items:        10_000,
			ejectionMode: arraylinkedlist.LRUEjection,
		},
		{ // mempool has the limit of 100K, and we put exactly 100K items
			limit:        100_000,
			items:        100_000,
			ejectionMode: arraylinkedlist.LRUEjection,
		},
		{ // mempool has the limit of 100K, and we put 1M items with LRU ejection.
			limit:        100_000,
			items:        1_000_000,
			ejectionMode: arraylinkedlist.LRUEjection,
		},
		{ // mempool has the limit of 100K, and we put 1M items with random ejection.
			limit:        100_000,
			items:        1_000_000,
			ejectionMode: arraylinkedlist.RandomEjection,
		},
	}

	for _, tc := range tt {
		t.Run(fmt.Sprintf("%d-limit-%d-items-%s-ejection", tc.limit, tc.items, tc.ejectionMode), func(t *testing.T) {
			bd := NewArrayBackData(tc.limit, 8, tc.ejectionMode)
			entities := unittest.EntityListFixture(uint(tc.items))

			testAddEntities(t, bd, entities)

			if tc.ejectionMode == arraylinkedlist.RandomEjection {
				// in random ejection mode we count total number of matched entities
				// with All map.
				testMapMatchCount(t, bd.All(), entities, int(tc.limit))
			} else {
				// in LRU ejection mode we match All items based on a from index (i.e., last "from" items).
				from := int(tc.items) - int(tc.limit)
				if from < 0 {
					// we are below limit, hence we start matching from index 0
					from = 0
				}
				testMapMatchFrom(t, bd.All(), entities, from)
			}
		})
	}
}

// TestArrayBackData_Rem checks correctness of All method when mempool is not full yet.
func TestArrayBackData_Rem(t *testing.T) {
	tt := []struct {
		limit uint32
		items uint32
		from  int // index start to be removed (set -1 to remove randomly)
		count int // total elements to be removed
	}{
		{ // removing range from mempool with items below its limit
			limit: 100_000,
			items: 10_000,
			from:  188,
			count: 2012,
		},
		{ // removing range from full mempool
			limit: 100_000,
			items: 100_000,
			from:  50_333,
			count: 6667,
		},
		{ // removing random from mempool with items below its limit
			limit: 100_000,
			items: 10_000,
			from:  -1,
			count: 6888,
		},
		{ // removing random from full mempool
			limit: 100_000,
			items: 10_000,
			from:  -1,
			count: 7328,
		},
	}

	for _, tc := range tt {
		t.Run(fmt.Sprintf("%d-limit-%d-items-%dfrom-%dcount", tc.limit, tc.items, tc.from, tc.count), func(t *testing.T) {
			bd := NewArrayBackData(tc.limit, 8, arraylinkedlist.RandomEjection)
			entities := unittest.EntityListFixture(uint(tc.items))

			testAddEntities(t, bd, entities)

			if tc.from == -1 {
				// random removal
				testRemoveAtRandom(t, bd, entities, tc.count)
				// except removed ones, the rest must be retrievable
				testRetrievableCount(t, bd, entities, uint64(int(tc.items)-tc.count))
			} else {
				testRemoveRange(t, bd, entities, tc.from, tc.from+tc.count)
				testCheckRangeRemoved(t, bd, entities, tc.from, tc.from+tc.count)
			}
		})
	}
}

// testAddEntities is a test helper that checks entities are added successfully to the backdata.
// and each entity is retrievable right after it is written to backdata.
func testAddEntities(t *testing.T, bd *ArrayBackData, entities []*unittest.MockEntity) {
	// adding elements
	for i, e := range entities {
		// adding each element must be successful.
		require.True(t, bd.Add(e.ID(), e))

		if uint32(i) < bd.limit {
			// when we are below limit the total of
			// backdata should be incremented by each addition.
			require.Equal(t, bd.Size(), uint(i+1))
		} else {
			// when we cross the limit, the ejection kicks in, and
			// size must be steady at the limit.
			require.Equal(t, uint32(bd.Size()), bd.limit)
		}

		// entity should be immediately retrievable
		actual, ok := bd.ByID(e.ID())
		require.True(t, ok)
		require.Equal(t, e, actual)
	}
}

// testGettingEntities is a test helper that checks entities are retrievable from backdata.
func testRetrievableFrom(t *testing.T, bd *ArrayBackData, entities []*unittest.MockEntity, from int) {
	for i := range entities {
		expected := entities[i]
		actual, ok := bd.ByID(expected.ID())
		if i < from {
			require.False(t, ok, i)
			require.Nil(t, actual)
		} else {
			require.True(t, ok)
			require.Equal(t, expected, actual)
		}
	}
}

func testRemoveAtRandom(t *testing.T, bd *ArrayBackData, entities []*unittest.MockEntity, count int) {
	for removedCount := 0; removedCount < count; {
		unittest.RequireReturnsBefore(t, func() {
			index := rand.Int() % len(entities)
			expected, removed := bd.Rem(entities[index].ID())
			if !removed {
				return
			}
			require.Equal(t, entities[index], expected)
			removedCount++
			// size sanity check after removal
			require.Equal(t, bd.Size(), uint(len(entities)-removedCount))
		}, 100*time.Millisecond, "could not find element to remove")
	}
}

func testRemoveRange(t *testing.T, bd *ArrayBackData, entities []*unittest.MockEntity, from int, to int) {
	for i := from; i < to; i++ {
		expected, removed := bd.Rem(entities[i].ID())
		require.True(t, removed)
		require.Equal(t, entities[i], expected)
		// size sanity check after removal
		require.Equal(t, bd.Size(), uint(len(entities)-(i-from)-1))
	}
}

func testCheckRangeRemoved(t *testing.T, bd *ArrayBackData, entities []*unittest.MockEntity, from int, to int) {
	for i := from; i < to; i++ {
		// bot removal and retrieval must fail
		expected, removed := bd.Rem(entities[i].ID())
		require.False(t, removed)
		require.Nil(t, expected)

		expected, exists := bd.ByID(entities[i].ID())
		require.False(t, exists)
		require.Nil(t, expected)
	}
}

// testMapMatchFrom is a test helper that checks entities are retrievable from entitiesMap.
func testMapMatchFrom(t *testing.T, entitiesMap map[flow.Identifier]flow.Entity, entities []*unittest.MockEntity, from int) {
	require.Len(t, entitiesMap, len(entities)-from)

	for i := range entities {
		expected := entities[i]
		actual, ok := entitiesMap[expected.ID()]
		if i < from {
			require.False(t, ok, i)
			require.Nil(t, actual)
		} else {
			require.True(t, ok)
			require.Equal(t, expected, actual)
		}
	}
}

// testMapMatchFrom is a test helper that checks entities are retrievable from entitiesMap.
func testMapMatchCount(t *testing.T, entitiesMap map[flow.Identifier]flow.Entity, entities []*unittest.MockEntity, count int) {
	require.Len(t, entitiesMap, count)
	actualCount := 0
	for i := range entities {
		expected := entities[i]
		actual, ok := entitiesMap[expected.ID()]
		if ok {
			require.Equal(t, expected, actual)
			actualCount++
		}
	}
}

// testRetrievableCount is a test helper that checks the number of retrievable entities from backdata exactly matches
// the expectedCount.
func testRetrievableCount(t *testing.T, bd *ArrayBackData, entities []*unittest.MockEntity, expectedCount uint64) {
	actualCount := 0

	for i := range entities {
		expected := entities[i]
		actual, ok := bd.ByID(expected.ID())
		if !ok {
			continue
		}
		require.Equal(t, expected, actual)
		actualCount++
	}

	require.Equal(t, int(expectedCount), actualCount)
}
