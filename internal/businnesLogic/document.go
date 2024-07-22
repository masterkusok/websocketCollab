package businnesLogic

import (
	"gorm.io/gorm"
	"time"
)

type liter struct {
	position      float32
	value         byte
	isSpecialChar bool
}

type Document struct {
	gorm.Model
	ID        uint      `json:"id" gorm:"id"`
	Text      string    `json:"text" gorm:"text"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	data      []liter
}

func CreateDocument() *Document {
	return &Document{data: []liter{{0, '$', true}, {1, '$', true}}}
}

func (d *Document) Insert(position float32, value byte) {
	for i := 1; i < len(d.data); i++ {
		if position < d.data[i].position && position > d.data[i-1].position {
			newData := append([]liter{}, d.data[:i]...)
			newData = append(newData, liter{position, value, false})
			newData = append(newData, d.data[i:]...)
			d.data = newData
			d.Text = d.Text[:i-1] + string(value) + d.Text[i-1:]
			return
		}
	}
}

func (d *Document) Delete(position float32) {
	if position == 0 || position == 1 {
		return
	}
	for i := 0; i < len(d.data); i++ {
		if d.data[i].position == position {
			d.Text = d.Text[:i-1] + d.Text[i:]
			d.data = append(d.data[:i], d.data[i+1:]...)
		}
	}
}

func (d *Document) GetText() string {
	return d.Text
}

func (d *Document) PullData() {
	newData := []liter{{0, '$', true}}
	for i := range d.Text {
		newData = append(newData, liter{(float32(i) + 1.0) * 1.0 / float32(len(newData)+1), d.Text[i], false})
	}
	newData = append(newData, liter{1.0, '$', true})
	d.data = newData
}
