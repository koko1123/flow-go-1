package backdata

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/utils/unittest"
)

func TestArrayBackData_BelowLimit(t *testing.T) {
	size := 10

	bd := NewArrayBackData(uint32(size), 1, LRUEjection)

	entities := unittest.EntityListFixture(uint(size))

	// adding elements
	for i, e := range entities {
		// adding each element must be successful.
		require.True(t, bd.Add(e.ID(), e))

		// total of back data should be incremented by each addition.
		require.Equal(t, bd.Size(), uint(i+1))

		// entity should be placed at index i in back data
		id, entity, _ := bd.entities.get(uint32(i))
		require.Equal(t, e.ID(), id)
		require.Equal(t, e, entity)
	}

	// sanity checks
	for i := range entities {
		// since we are below limit, elements should be added sequentially at bucket 0.
		// first added element has a key index of 1, since 0 means unused key index in implementation.
		require.Equal(t, bd.buckets[0][i].keyIndex, uint32(i+1))
		// also, since we have not yet over-limited, entities are received valueIndex in the same order they
		// are added.
		require.Equal(t, bd.buckets[0][i].valueIndex, uint32(i))
		_, _, owner := bd.entities.get(uint32(i))
		require.Equal(t, owner, uint32(i))
	}

	// getting inserted elements
	for _, expected := range entities {
		actual, ok := bd.ByID(expected.ID())
		require.True(t, ok)
		require.Equal(t, expected, actual)
	}
}

func TestArrayBackDataStoreAndRetrievalWithoutEjection(t *testing.T) {
	for _, tc := range []struct {
		limit           uint32
		overLimitFactor uint32
		entityCount     uint32
		helpers         []func(*testing.T, *ArrayBackData, []*unittest.MockEntity)
	}{
		{ // two buckets, entities below limit.
			limit:           30,
			overLimitFactor: 2,
			entityCount:     10,
		},
		{ // two buckets, entities equal to limit.
			limit:           30,
			overLimitFactor: 2,
			entityCount:     30,
		},
		{ // multiple buckets, high limit, low entities.
			limit:           10000,
			overLimitFactor: 16,
			entityCount:     1000,
		},
		{ // multiple buckets, entities equal to limit.
			limit:           10000,
			overLimitFactor: 16,
			entityCount:     10000,
		},
	} {
		t.Run(fmt.Sprintf("%d-limit-%d-overlimit-%d-entities", tc.limit, tc.overLimitFactor, tc.entityCount), func(t *testing.T) {
			testArrayBackDataStoreAndRetrievalWithoutEjection(t, tc.limit, tc.overLimitFactor, tc.entityCount)
		})
	}
}

func testArrayBackDataStoreAndRetrievalWithoutEjection(t *testing.T, limit uint32, overLimitFactor uint32, entityCount uint32, helpers ...func(*testing.T, *ArrayBackData, []*unittest.MockEntity)) {
	h := []func(*testing.T, *ArrayBackData, []*unittest.MockEntity){
		func(t *testing.T, backData *ArrayBackData, entities []*unittest.MockEntity) {
			testInitialization(t, backData, entities)
		},
		func(t *testing.T, backData *ArrayBackData, entities []*unittest.MockEntity) {
			testAddingEntities(t, backData, entities)
		},
	}
	h = append(h, helpers...)

	withTestScenario(t, limit, overLimitFactor, entityCount,
		append(h, func(t *testing.T, backData *ArrayBackData, entities []*unittest.MockEntity) {
			testRetrievingSavedEntities(t, backData, entities)
		})...,
	)
}

func testArrayBackDataStoreAndRetrievalWitEjection(t *testing.T, limit uint32, overLimitFactor uint32, entityCount uint32, helpers ...func(*testing.T, *ArrayBackData, []*unittest.MockEntity)) {
	h := []func(*testing.T, *ArrayBackData, []*unittest.MockEntity){
		func(t *testing.T, backData *ArrayBackData, entities []*unittest.MockEntity) {
			testAddingEntities(t, backData, entities)
		},
	}
	h = append(h, helpers...)

	withTestScenario(t, limit, overLimitFactor, entityCount,
		append(h, func(t *testing.T, backData *ArrayBackData, entities []*unittest.MockEntity) {
			testRetrievingSavedEntities(t, backData, entities)
		})...,
	)
}

func TestInvalidateEntity(t *testing.T) {
	for _, tc := range []struct {
		limit           uint32
		overLimitFactor uint32
		entityCount     uint32
		helpers         []func(*testing.T, *ArrayBackData, []*unittest.MockEntity)
	}{
		{ // edge-case empty list
			limit:           30,
			overLimitFactor: 2,
			entityCount:     0,
		},
		{ // edge-case single element-list
			limit:           30,
			overLimitFactor: 2,
			entityCount:     1,
		},
		{ // two buckets, entities below limit.
			limit:           30,
			overLimitFactor: 2,
			entityCount:     10,
		},
		{ // two buckets, entities equal to limit.
			limit:           30,
			overLimitFactor: 2,
			entityCount:     30,
		},
		{ // multiple buckets, high limit, low entities.
			limit:           100,
			overLimitFactor: 16,
			entityCount:     10,
		},
		{ // multiple buckets, entities equal to limit.
			limit:           100,
			overLimitFactor: 16,
			entityCount:     100,
		},
	} {
		// head invalidation test (LRU)
		t.Run(fmt.Sprintf("head-invalidation-%d-limit-%d-overlimit-%d-entities", tc.limit, tc.overLimitFactor, tc.entityCount), func(t *testing.T) {
			testInvalidateEntity(t, tc.limit, tc.overLimitFactor, tc.entityCount, func(t *testing.T, backData *ArrayBackData, entities []*unittest.MockEntity) {
				testInvalidatingHead(t, backData, entities)
			})
		})
		// tail invalidation test
		t.Run(fmt.Sprintf("tail-invalidation-%d-limit-%d-overlimit-%d-entities-", tc.limit, tc.overLimitFactor, tc.entityCount), func(t *testing.T) {
			testInvalidateEntity(t, tc.limit, tc.overLimitFactor, tc.entityCount, func(t *testing.T, backData *ArrayBackData, entities []*unittest.MockEntity) {
				testInvalidatingTail(t, backData, entities)
			})
		})
		// random invalidation test
		t.Run(fmt.Sprintf("random-invalidation-%d-limit-%d-overlimit-%d-entities-", tc.limit, tc.overLimitFactor, tc.entityCount),
			func(t *testing.T) {
				testInvalidateEntity(t, tc.limit, tc.overLimitFactor, tc.entityCount, func(t *testing.T, backData *ArrayBackData, entities []*unittest.MockEntity) {
					testInvalidateAtRandom(t, backData, entities)
				})
			})
	}
}

func testInvalidateEntity(t *testing.T, limit uint32, overLimitFactor uint32, entityCount uint32, helpers ...func(*testing.T, *ArrayBackData, []*unittest.MockEntity)) {
	h := append([]func(*testing.T, *ArrayBackData, []*unittest.MockEntity){
		func(t *testing.T, backData *ArrayBackData, entities []*unittest.MockEntity) {
			testAddingEntities(t, backData, entities)
		},
	}, helpers...)

	withTestScenario(t, limit, overLimitFactor, entityCount, h...)
}

func withTestScenario(t *testing.T,
	size uint32,
	overLimitFactor uint32,
	entityCount uint32,
	helpers ...func(*testing.T, *ArrayBackData, []*unittest.MockEntity)) {

	bd := NewArrayBackData(size, overLimitFactor, LRUEjection)
	// head on underlying linked list value should be uninitialized
	require.True(t, bd.entities.used.head.isUndefined())
	require.Equal(t, bd.Size(), uint(0))
	entities := unittest.EntityListFixture(uint(entityCount))

	for _, helper := range helpers {
		helper(t, bd, entities)
	}
}

// testInitialization evaluates the state of an initialized cachedEntity list before adding any element
// to it.
func testInitialization(t *testing.T, backData *ArrayBackData, _ []*unittest.MockEntity) {
	// head and tail of "used" linked-list must be undefined at initialization time.
	require.True(t, backData.entities.used.head.isUndefined())
	require.True(t, backData.entities.used.tail.isUndefined())

	for i := 0; i < len(backData.entities.entities); i++ {
		if i == 0 {
			// head of embedded "free" linked-list should point to index 0 of entities slice.
			require.Equal(t, uint32(i), backData.entities.free.head.sliceIndex())
			// previous element of tail must be undefined.
			require.True(t, backData.entities.entities[i].prev.isUndefined())
		}

		if i != 0 {
			// except head, any element should point back to its previous index in slice.
			require.Equal(t, uint32(i-1), backData.entities.entities[i].prev.sliceIndex())
		}

		if i != len(backData.entities.entities)-1 {
			// except tail, any element should point forward to its next index in slice.
			require.Equal(t, uint32(i+1), backData.entities.entities[i].next.sliceIndex())
		}

		if i == len(backData.entities.entities)-1 {
			// tail of embedded "free" linked-list should point to the last index in entities slice.
			require.Equal(t, uint32(i), backData.entities.free.tail.sliceIndex())
			// next element of tail must be undefined.
			require.True(t, backData.entities.entities[i].next.isUndefined())
		}
	}
}

func testAddingEntities(t *testing.T, backData *ArrayBackData, entities []*unittest.MockEntity) {
	// adding elements
	for i, e := range entities {
		// adding each element must be successful.
		require.True(t, backData.Add(e.ID(), e))

		// total of back data should be incremented by each addition.
		require.Equal(t, backData.Size(), uint(i+1))

		// entity should be placed at index i in back data
		_, entity, _ := backData.entities.get(uint32(i))
		require.Equal(t, e, entity)

		// linked-list sanity check
		// first insertion forward, head of backData should always point to
		// first entity in the list.
		require.Equal(t, entities[0], backData.entities.getHead().entity)
		require.True(t, backData.entities.getHead().prev.isUndefined())
		require.Equal(t, entities[i], backData.entities.getTail().entity)
		require.True(t, backData.entities.getTail().next.isUndefined())

		tailAccessibleFromHead(t,
			backData.entities.used.head.sliceIndex(),
			backData.entities.used.tail.sliceIndex(),
			backData,
			uint32(i+1))
		headAccessibleFromTail(t, backData, uint32(i+1))

		//tailAccessibleFromHead(t,
		//	backData.entities.free.head.sliceIndex(),
		//	backData.entities.free.tail.sliceIndex(),
		//	backData,
		//	uint32(len(entities)-i-1))
	}
}

func testRetrievingSavedEntities(t *testing.T, backData *ArrayBackData, entities []*unittest.MockEntity) {
	for _, expected := range entities {
		actual, ok := backData.ByID(expected.ID())
		require.True(t, ok)
		require.Equal(t, expected, actual)
	}
}

// testInvalidatingHead keeps invalidating elements at random and evaluates whether double-linked list remains
// connected on both head and tail.
func testInvalidateAtRandom(t *testing.T, backData *ArrayBackData, entities []*unittest.MockEntity) {
	size := len(entities)
	for i := 0; i < size; i++ {
		backData.entities.invalidateRandomEntity()

		// size of list should be shrunk after each invalidation.
		require.Equal(t, uint32(size-i-1), backData.entities.size())

		// except when the list is empty, head and tail must be accessible after each invalidation
		// i.e., the linked list remains connected despite invalidation.
		if i != size-1 {
			tailAccessibleFromHead(t,
				backData.entities.used.head.sliceIndex(),
				backData.entities.used.tail.sliceIndex(),
				backData,
				backData.entities.size())
			//tailAccessibleFromHead(t,
			//	backData.entities.free.head.sliceIndex(),
			//	backData.entities.free.tail.sliceIndex(),
			//	backData,
			//	backData.entities.size())

			headAccessibleFromTail(t, backData, backData.entities.size())
		}
	}
}

// testInvalidatingHead keeps invalidating the head and evaluates the linked-list keeps updating its head
// and remains connected.
func testInvalidatingHead(t *testing.T, backData *ArrayBackData, entities []*unittest.MockEntity) {
	size := len(entities)
	for i := 0; i < size; i++ {
		headIndex := backData.entities.invalidateHead()
		require.Equal(t, uint32(i), headIndex)

		// size of list should be shrunk after each invalidation.
		require.Equal(t, uint32(size-i-1), backData.entities.size())
		// unclaimed head should be appended to free entities
		require.Equal(t, backData.entities.free.tail.sliceIndex(), headIndex)

		// except when the list is empty, head must be updated after invalidation,
		// except when the list is empty, head and tail must be accessible after each invalidation
		// i.e., the linked list remains connected despite invalidation.
		if i != size-1 {
			require.Equal(t, entities[i+1].ID(), backData.entities.getHead().id)
			tailAccessibleFromHead(t,
				backData.entities.used.head.sliceIndex(),
				backData.entities.used.tail.sliceIndex(),
				backData,
				backData.entities.size())
			headAccessibleFromTail(t, backData, backData.entities.size())
		}
	}
}

// testInvalidatingHead keeps invalidating the tail and evaluates the linked-list keeps updating its tail
// and remains connected.
func testInvalidatingTail(t *testing.T, backData *ArrayBackData, entities []*unittest.MockEntity) {
	size := len(entities)
	for i := 0; i < size; i++ {
		// invalidates tail index
		tail := backData.entities.used.tail.sliceIndex()
		backData.entities.invalidateEntityAtIndex(tail)
		// old head index must be invalidated
		require.True(t, backData.entities.isInvalidated(tail))

		// size of list should be shrunk after each invalidation.
		require.Equal(t, uint32(size-i-1), backData.entities.size())

		// except when the list is empty, tail must be updated after invalidation,
		// and also head and tail must be accessible after each invalidation
		// i.e., the linked list remains connected despite invalidation.
		if i != size-1 {
			require.Equal(t, entities[size-i-2].ID(), backData.entities.getTail().id)
			tailAccessibleFromHead(t,
				backData.entities.used.head.sliceIndex(),
				backData.entities.used.tail.sliceIndex(),
				backData,
				backData.entities.size())
			headAccessibleFromTail(t, backData, backData.entities.size())
		}
	}
}

func tailAccessibleFromHead(t *testing.T, headSliceIndex uint32, tailSliceIndex uint32, backData *ArrayBackData, total uint32) {
	seen := make(map[uint32]struct{})

	index := headSliceIndex
	for i := uint32(0); i < total; i++ {
		if i == total-1 {
			require.Equal(t, tailSliceIndex, index, "tail not reachable after total steps")
			return
		}

		require.NotEqual(t, tailSliceIndex, index, "tail visited in less expected steps (potential inconsistency)", i, total)
		_, ok := seen[index]
		require.False(t, ok, "duplicate identifiers found")

		index = backData.entities.entities[index].next.sliceIndex()
	}
}

func headAccessibleFromTail(t *testing.T, backData *ArrayBackData, steps uint32) {
	seen := make(map[flow.Identifier]struct{})
	headId := backData.entities.getHead().id

	n := backData.entities.getTail()
	for i := uint32(0); i < steps; i++ {
		if i == steps-1 {
			require.Equal(t, headId, n.id, "tail not reachable after total steps")
			return
		}

		require.NotEqual(t, headId, n.id, "tail visited in less expected steps (potential inconsistency)", i, steps)
		_, ok := seen[*n.id]
		require.False(t, ok, "duplicate identifiers found")

		n = &backData.entities.entities[n.prev.sliceIndex()]
	}
}
