package ecs

type Entity struct {
	Id         int32
	Components map[string]*Component
}

func NewEntity(id int32) *Entity {
	return &Entity{
		Id:         id,
		Components: make(map[string]*Component),
	}
}

func (e *Entity) AddComponent(component *Component) {
	e.Components[(*component).Name()] = component
}

func (e *Entity) GetComponent(name string) *Component {
	return e.Components[name]
}
