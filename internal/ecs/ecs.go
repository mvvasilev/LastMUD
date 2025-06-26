package ecs

import (
	"iter"
	"maps"
	"slices"
	"time"

	"code.haedhutner.dev/mvv/LastMUD/internal/logging"
	"code.haedhutner.dev/mvv/LastMUD/internal/util"
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

type ComponentStorage[T Component] struct {
	forType ComponentType
	storage map[Entity]T
}

func CreateComponentStorage[T Component](forType ComponentType) *ComponentStorage[T] {
	return &ComponentStorage[T]{
		forType: forType,
		storage: map[Entity]T{},
	}
}

func (cs *ComponentStorage[T]) ComponentType() ComponentType {
	return cs.forType
}

func (cs *ComponentStorage[T]) Entities() iter.Seq[Entity] {
	return maps.Keys(cs.storage)
}

func (cs *ComponentStorage[T]) Query(query func(comp T) bool) iter.Seq[Entity] {
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

func (cs *ComponentStorage[T]) Set(e Entity, component T) {
	cs.storage[e] = component
}

func (cs *ComponentStorage[T]) Get(e Entity) (component T, ok bool) {
	component, ok = cs.storage[e]
	return
}

func (cs *ComponentStorage[T]) Delete(e Entity) {
	delete(cs.storage, e)
}

func (cs *ComponentStorage[T]) All() map[Entity]T {
	return cs.storage
}

type System struct {
	name     string
	priority int
	work     func(world *World, delta time.Duration) (err error)
}

func CreateSystem(name string, priority int, work func(world *World, delta time.Duration) (err error)) *System {
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
		logging.Error("Error in system '", s.name, "': ", err.Error())
	}
}

type World struct {
	systems          []*System
	componentsByType map[ComponentType]any
	resources        map[Resource]any
}

func CreateWorld() (world *World) {
	world = &World{
		systems:          []*System{},
		componentsByType: map[ComponentType]any{},
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
		storage := s.(*ComponentStorage[Component])

		storage.Delete(entity)
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

func registerComponent[T Component](world *World, compType ComponentType) {
	if _, ok := world.componentsByType[compType]; ok {
		return
	}

	world.componentsByType[compType] = CreateComponentStorage[T](compType)
}

func SetComponent[T Component](world *World, entity Entity, component T) {
	registerComponent[T](world, component.Type())

	compStorage := world.componentsByType[component.Type()].(*ComponentStorage[T])

	compStorage.Set(entity, component)
}

func GetComponent[T Component](world *World, entity Entity) (component T, exists bool) {
	storage := GetComponentStorage[T](world)

	return storage.Get(entity)
}

func DeleteComponent[T Component](world *World, entity Entity) {
	storage := GetComponentStorage[T](world)

	storage.Delete(entity)
}

func GetComponentStorage[T Component](world *World) (compStorage *ComponentStorage[T]) {
	var zero T

	compType := zero.Type()

	registerComponent[T](world, compType)

	return world.componentsByType[compType].(*ComponentStorage[T])
}

func IterateEntitiesWithComponent[T Component](world *World) iter.Seq[Entity] {
	storage := GetComponentStorage[T](world)

	return storage.Entities()
}

func QueryEntitiesWithComponent[T Component](world *World, query func(comp T) bool) iter.Seq[Entity] {
	storage := GetComponentStorage[T](world)

	return storage.Query(query)
}

func FindEntitiesWithComponents(world *World, componentTypes ...ComponentType) (entities []Entity) {
	entities = []Entity{}

	isFirst := true

	for _, compType := range componentTypes {
		// If we've gone through at least one component, and we have an empty result already, return it
		if !isFirst && len(entities) == 0 {
			return
		}

		storage, ok := world.componentsByType[compType].(*ComponentStorage[Component])

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
		entities = util.IntersectSliceWithIterator(entities, storage.Entities())
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
	for _, s := range systems {
		RegisterSystem(world, s)
	}
}
