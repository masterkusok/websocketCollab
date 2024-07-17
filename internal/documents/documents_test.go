package documents

import "testing"

func TestInsert(t *testing.T) {
	doc := CreateDocument("Test document!")
	doc.Insert(0.5, 'i')
	doc.Insert(0.75, 'a')
	doc.Insert(0.87, 'm')

	result := doc.Text()
	expected := "iam"
	if result != expected {
		t.Fail()
	}
}

func TestDelete(t *testing.T) {
	doc := CreateDocument("Test document!")
	doc.Insert(0.5, 'i')
	doc.Insert(0.75, 'a')
	doc.Insert(0.87, 'm')

	result := doc.Text()

	if result != "iam" {
		t.Fail()
	}

	doc.Delete(0.5)

	result = doc.Text()
	if result != "am" {
		t.Fail()
	}
}
