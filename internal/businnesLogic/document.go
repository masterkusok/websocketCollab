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
	if len(d.data) == 2 {
		d.insertAt(liter{position: position, value: value, isSpecialChar: false}, 1)
		return
	}

	l, r, m := 0, len(d.data)-1, 0

	for l <= r {
		m = l + (r-l)/2
		if d.data[m-1].position < position && d.data[m].position > position {
			d.insertAt(liter{position: position, value: value, isSpecialChar: false}, m)
		}

		if d.data[m-1].position > position {
			r = m - 1
		} else {
			l = m + 1
		}
	}
}

func (d *Document) insertAt(l liter, index int) {
	newData := append([]liter{}, d.data[:index]...)
	newData = append(newData, l)
	newData = append(newData, d.data[index:]...)
	d.data = newData
	d.Text = d.Text[:index-1] + string(l.value) + d.Text[index-1:]
}

func (d *Document) Delete(position float32) {
	if position == 0 || position == 1 {
		return
	}
	l, r, m := 0, len(d.data)-1, 0

	for l <= r {
		m = l + (r-l)/2
		if d.data[m].position == position {
			d.Text = d.Text[:m-1] + d.Text[m:]
			d.data = append(d.data[:m], d.data[m+1:]...)
		}
		if d.data[m].position < position {
			l = m + 1
		} else {
			r = m - 1
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
