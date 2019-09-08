package interpreter

type StructObject struct {
	Slots map[string]Object
}

func (o StructObject) GetSlots() map[string]Object {
	return o.Slots
}
