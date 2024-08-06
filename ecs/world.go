package ecs

type World struct {
	Name           string
	nextSystemId   int32
	Systems        map[int32]*System
	nextEntityId   int32
	EntityRegistry map[int32]*Entity
}

func NewWorld(name string) *World {
	return &World{
		Name:           name,
		nextSystemId:   0,
		Systems:        make(map[int32]*System),
		nextEntityId:   0,
		EntityRegistry: make(map[int32]*Entity),
	}
}

func (w *World) RegisterNewSystem(system *System) int32 {
	w.Systems[w.nextSystemId] = system
	w.nextSystemId++
	return w.nextSystemId - 1
}

func (w *World) MakeNewEntity() int32 {
	entity := NewEntity(w.nextEntityId)
	w.EntityRegistry[w.nextEntityId] = entity
	w.nextEntityId++
	return entity.Id
}

func (w *World) GetEntity(id int32) *Entity {
	return w.EntityRegistry[id]
}

func (w *World) GetSystem(id int32) *System {
	return w.Systems[id]
}

func (w *World) AddComponentToEntity(entityId int32, component *Component) {
	entity := w.GetEntity(entityId)
	entity.AddComponent(component)
}

func (w *World) Update(deltaTime float32) {
	for _, system := range w.Systems {
		(*system).Update(w, deltaTime)
	}
}
