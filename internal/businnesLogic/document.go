package businnesLogic

import (
	"gorm.io/gorm"
	"strings"
	"time"
)

type Liter struct {
	Position float32
	Value    byte
}

type Document struct {
	gorm.Model
	ID        uint      `json:"id" gorm:"id"`
	Text      string    `json:"text" gorm:"text"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	data      []Liter
}

func CreateDocument() *Document {
	return &Document{data: []Liter{{0, '$'}, {1, '$'}}}
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

func (d *Document) ParseText() {
	output := strings.Builder{}
	for _, b := range d.data {
		if b.Value == '$' {
			continue
		}
		output.WriteByte(b.Value)
	}
	d.Text = output.String()
}

func (d *Document) PullData() {
	newData := []Liter{{0, '$'}}
	for i := range d.Text {
		newData = append(newData, Liter{(float32(i) + 1.0) * 1.0 / float32(len(newData)+1), d.Text[i]})
	}
	newData = append(newData, Liter{1.0, '$'})
	d.data = newData
}
