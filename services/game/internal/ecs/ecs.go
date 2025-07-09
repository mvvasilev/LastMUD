package ecs

import (
	"code.haedhutner.dev/mvv/LastMUD/services/game/internal/gameutils"
	"code.haedhutner.dev/mvv/LastMUD/shared/log"
	"iter"
	"maps"
	"slices"
	"time"

	"github.com/google/uuid"
)

type Entity uuid.UUID

func CreateEntity(uuid uuid.UUID) Entity {
	return Entity(uuid)
}

func NewEntity() Entity {
	return Entity(uuid.New())
}

func NilEntity() Entity {
	return Entity(uuid.Nil)
}

func (e Entity) AsUUID() uuid.UUID {
	return uuid.UUID(e)
}

type ComponentType int16

type Resource string

type Component interface {
	Type() ComponentType
}

type ComponentStorage struct {
	forType ComponentType
	storage map[Entity]Component
}

func CreateComponentStorage(forType ComponentType) *ComponentStorage {
	return &ComponentStorage{
		forType: forType,
		storage: map[Entity]Component{},
	}
}

func (cs *ComponentStorage) ComponentType() ComponentType {
	return cs.forType
}

func (cs *ComponentStorage) Entities() iter.Seq[Entity] {
	return maps.Keys(cs.storage)
}

func queryComponents(cs *ComponentStorage, query func(comp Component) bool) iter.Seq[Entity] {
	return func(yield func(Entity) bool) {
		for k, v := range cs.storage {
			if !query(v) {
				continue
			}

			if !yield(k) {
				return
			}
		}
	}
}

func (cs *ComponentStorage) Set(e Entity, component Component) {
	cs.storage[e] = component
}

func (cs *ComponentStorage) Get(e Entity) (component Component, ok bool) {
	component, ok = cs.storage[e]
	return
}

func (cs *ComponentStorage) Delete(e Entity) {
	delete(cs.storage, e)
}

func (cs *ComponentStorage) All() map[Entity]Component {
	return cs.storage
}

type SystemExecutor func(world *World, delta time.Duration) (err error)

type System struct {
	name     string
	priority int
	work     SystemExecutor
}

func CreateSystem(name string, priority int, work SystemExecutor) *System {
	return &System{
		name:     name,
		priority: priority,
		work:     work,
	}
}

func (s *System) Priority() int {
	return s.priority
}

func (s *System) Execute(world *World, delta time.Duration) {
	err := s.work(world, delta)

	if err != nil {
		log.Error("Error in system '", s.name, "': ", err.Error())
	}
}

type World struct {
	systems          []*System
	componentsByType map[ComponentType]*ComponentStorage
	resources        map[Resource]any
}

func CreateWorld() (world *World) {
	world = &World{
		systems:          []*System{},
		componentsByType: map[ComponentType]*ComponentStorage{},
		resources:        map[Resource]any{},
	}

	return
}

func (w *World) Tick(delta time.Duration) {
	for _, s := range w.systems {
		s.Execute(w, delta)
	}
}

func DeleteEntity(world *World, entity Entity) {
	for _, s := range world.componentsByType {
		s.Delete(entity)
	}
}

func DeleteEntities(world *World, entities ...Entity) {
	for _, e := range entities {
		DeleteEntity(world, e)
	}
}

func SetResource(world *World, r Resource, val any) {
	world.resources[r] = val
}

func GetResource[T any](world *World, r Resource) (res T, err error) {
	val, ok := world.resources[r]

	if !ok {
		err = newECSError("Resource '", r, "' not found.")
		return
	}

	res, ok = val.(T)

	if !ok {
		err = newECSError("Incompatible type for resource '", r, "'")
	}

	return
}

func RemoveResource(world *World, r Resource) {
	delete(world.resources, r)
}

func registerComponent(world *World, compType ComponentType) {
	if _, ok := world.componentsByType[compType]; ok {
		return
	}

	world.componentsByType[compType] = CreateComponentStorage(compType)
}

func SetComponent[T Component](world *World, entity Entity, component T) {
	registerComponent(world, component.Type())

	compStorage := world.componentsByType[component.Type()]

	compStorage.Set(entity, component)
}

func GetComponent[T Component](world *World, entity Entity) (component T, exists bool) {
	storage := GetComponentStorage[T](world)

	val, exists := storage.Get(entity)
	casted, castSuccess := val.(T)

	return casted, exists && castSuccess
}

func DeleteComponent[T Component](world *World, entity Entity) {
	storage := GetComponentStorage[T](world)

	storage.Delete(entity)
}

func GetComponentStorage[T Component](world *World) (compStorage *ComponentStorage) {
	var zero T

	// This is ok because the `Type` function is expected to return a hard-coded value and not depend on component state
	compType := zero.Type()

	registerComponent(world, compType)

	return world.componentsByType[compType]
}

func IterateEntitiesWithComponent[T Component](world *World) iter.Seq[Entity] {
	storage := GetComponentStorage[T](world)

	return storage.Entities()
}

func QueryEntitiesWithComponent[T Component](world *World, query func(comp T) bool) iter.Seq[Entity] {
	storage := GetComponentStorage[T](world)

	return queryComponents(storage, func(comp Component) bool {
		val, ok := comp.(T)

		// Cast unsuccessful, assume failure
		if !ok {
			return false
		}

		return query(val)
	})
}

func QueryFirstEntityWithComponent[T Component](world *World, query func(comp T) bool) Entity {
	for entity := range QueryEntitiesWithComponent(world, query) {
		return entity
	}

	return NilEntity()
}

func FindEntitiesWithComponents(world *World, componentTypes ...ComponentType) (entities []Entity) {
	entities = []Entity{}

	isFirst := true

	for _, compType := range componentTypes {
		// If we've gone through at least one component, and we have an empty result already, return it
		if !isFirst && len(entities) == 0 {
			return
		}

		storage, ok := world.componentsByType[compType]

		// If we can't find the storage for this component, then it hasn't been used yet.
		// Therefore, no entity could have all components requested. Return empty.
		if !ok {
			return []Entity{}
		}

		// For the first component, simply add all entities to the array
		if isFirst {
			for entity := range storage.Entities() {
				entities = append(entities, entity)
			}

			isFirst = false
			continue
		}

		// For later components, intersect
		entities = gameutils.IntersectSliceWithIterator(entities, storage.Entities())
	}

	return entities
}

func RegisterSystem(world *World, s *System) {
	world.systems = append(world.systems, s)
	slices.SortFunc(
		world.systems,
		func(a, b *System) int {
			return a.priority - b.priority
		},
	)
}

func RegisterSystems(world *World, systems ...*System) {
	world.systems = append(world.systems, systems...)
	slices.SortFunc(
		world.systems,
		func(a, b *System) int {
			return a.priority - b.priority
		},
	)
}
