package documents

import "strings"

type Liter struct {
	Position float32
	Value    byte
}

type Document struct {
	Id   int
	Name string
	data []Liter
}

func CreateDocument(name string) *Document {
	return &Document{Id: 1, Name: name, data: []Liter{{0, '$'}, {1, '$'}}}
}

func (d *Document) Insert(position float32, value byte) {
	for i := 1; i < len(d.data); i++ {
		if position < d.data[i].Position && position > d.data[i-1].Position {
			newData := append([]Liter{}, d.data[:i]...)
			newData = append(newData, Liter{position, value})
			newData = append(newData, d.data[i:]...)
			d.data = newData
			return
		}
	}
}

func (d *Document) Delete(position float32) {
	if position == 0 || position == 1 {
		return
	}
	for i := 0; i < len(d.data); i++ {
		if d.data[i].Position == position {
			d.data = append(d.data[:i], d.data[i+1:]...)
		}
	}
}

func (d *Document) Text() string {
	output := strings.Builder{}
	for _, b := range d.data {
		if b.Value == '$' {
			continue
		}
		output.WriteByte(b.Value)
	}
	return output.String()
}
