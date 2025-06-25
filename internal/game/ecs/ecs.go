package ecs

import (
	"slices"
	"time"

	"code.haedhutner.dev/mvv/LastMUD/internal/logging"
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

func (s *System) DoWork(world *World, delta time.Duration) {
	err := s.work(world, delta)

	if err != nil {
		logging.Error("Error in system '", s.name, "': ", err.Error())
	}
}

type World struct {
	systems            []*System
	componentsByType   map[ComponentType]any
	componentsByEntity map[Entity]map[ComponentType]any
	resources          map[Resource]any
}

func CreateWorld() (world *World) {
	world = &World{
		systems:            []*System{},
		componentsByType:   map[ComponentType]any{},
		componentsByEntity: map[Entity]map[ComponentType]any{}, // TODO: Can't figure out use-case right now
		resources:          map[Resource]any{},
	}

	return
}

func (w *World) Tick(delta time.Duration) {
	for _, s := range w.systems {
		s.DoWork(w, delta)
	}
}

func DeleteEntity(world *World, entity Entity) {
	for _, s := range world.componentsByType {
		storage, ok := s.(*ComponentStorage[Component])

		if ok {
			storage.Delete(entity)
		}
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

func RegisterComponent[T Component](world *World, compType ComponentType) {
	world.componentsByType[compType] = CreateComponentStorage[T](compType)
}

func SetComponent[T Component](world *World, entity Entity, component T) {
	compStorage := world.componentsByType[component.Type()].(*ComponentStorage[T])
	compStorage.Set(entity, component)

	// if _, ok := world.componentsByEntity[entity]; !ok {
	// 	world.componentsByEntity[entity] = map[ComponentType]any{}
	// }

	// world.componentsByEntity[entity][component.Type()] = component
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

	return world.componentsByType[zero.Type()].(*ComponentStorage[T])
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
