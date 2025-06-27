package ecs

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type TestComponent struct {
	Value string
}

func (tc TestComponent) Type() ComponentType {
	return 1
}

func TestEntityCreation(t *testing.T) {
	u := uuid.New()
	entity := CreateEntity(u)
	assert.Equal(t, u, entity.AsUUID())

	newEntity := NewEntity()
	assert.NotEqual(t, uuid.Nil, newEntity.AsUUID())

	nilEntity := NilEntity()
	assert.Equal(t, uuid.Nil, nilEntity.AsUUID())
}

func TestComponentStorage(t *testing.T) {
	cs := CreateComponentStorage(1)
	entity := NewEntity()
	comp := TestComponent{"hello"}
	cs.Set(entity, comp)

	ret, ok := cs.Get(entity)
	assert.True(t, ok)
	assert.Equal(t, comp, ret)

	cs.Delete(entity)
	_, ok = cs.Get(entity)
	assert.False(t, ok)
}

func TestSystemExecution(t *testing.T) {
	called := false
	s := CreateSystem("test", 0, func(world *World, delta time.Duration) error {
		called = true
		return nil
	})

	world := CreateWorld()
	RegisterSystem(world, s)
	world.Tick(time.Millisecond)

	assert.True(t, called)
}

func TestWorldComponentLifecycle(t *testing.T) {
	world := CreateWorld()
	entity := NewEntity()
	comp := TestComponent{"value"}

	SetComponent(world, entity, comp)
	ret, ok := GetComponent[TestComponent](world, entity)
	assert.True(t, ok)
	assert.Equal(t, comp, ret)

	DeleteComponent[TestComponent](world, entity)
	_, ok = GetComponent[TestComponent](world, entity)
	assert.False(t, ok)
}

func TestWorldResourceManagement(t *testing.T) {
	world := CreateWorld()
	resourceKey := Resource("test_resource")
	value := "some data"

	SetResource(world, resourceKey, value)
	ret, err := GetResource[string](world, resourceKey)
	assert.NoError(t, err)
	assert.Equal(t, value, ret)

	RemoveResource(world, resourceKey)
	_, err = GetResource[string](world, resourceKey)
	assert.Error(t, err)
}

func TestEntityDeletion(t *testing.T) {
	world := CreateWorld()
	entity := NewEntity()
	comp := TestComponent{"remove me"}
	SetComponent(world, entity, comp)

	DeleteEntity(world, entity)
	_, ok := GetComponent[TestComponent](world, entity)
	assert.False(t, ok)
}

func TestIterateEntitiesWithComponent(t *testing.T) {
	world := CreateWorld()
	entities := []Entity{
		NewEntity(),
		NewEntity(),
		NewEntity(),
	}
	for _, e := range entities {
		SetComponent(world, e, TestComponent{"iter"})
	}

	count := 0
	for range IterateEntitiesWithComponent[TestComponent](world) {
		count++
	}

	assert.Equal(t, len(entities), count)
}

func TestQueryEntitiesWithComponent(t *testing.T) {
	world := CreateWorld()
	match := NewEntity()
	nonMatch := NewEntity()
	SetComponent(world, match, TestComponent{"match"})
	SetComponent(world, nonMatch, TestComponent{"skip"})

	results := []Entity{}

	for e := range QueryEntitiesWithComponent(world, func(c TestComponent) bool {
		return c.Value == "match"
	}) {
		results = append(results, e)
	}

	assert.Len(t, results, 1)
	assert.Equal(t, match, results[0])
}

type SecondComponent struct{ Flag bool }

func (SecondComponent) Type() ComponentType { return 2 }

func TestFindEntitiesWithComponents(t *testing.T) {
	world := CreateWorld()
	entity1 := NewEntity()
	entity2 := NewEntity()

	SetComponent(world, entity1, TestComponent{"a"})
	SetComponent(world, entity1, SecondComponent{true})
	SetComponent(world, entity2, TestComponent{"b"})

	results := FindEntitiesWithComponents(world, 1, 2)
	assert.Len(t, results, 1)
	assert.Equal(t, entity1, results[0])
}
