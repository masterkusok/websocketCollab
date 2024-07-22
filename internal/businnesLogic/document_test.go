package businnesLogic

import "testing"

func TestCreateDocument(t *testing.T) {
	doc := CreateDocument()

	// data array should contain 2 special chars indicating start and end of data
	if len(doc.data) != 2 {
		t.Logf("Wrong data length for empty document!\n")
		t.Fail()
	}

	if len(doc.Text) != 0 {
		t.Fail()
	}
}

func TestDocument_Insert(t *testing.T) {
	doc := CreateDocument()
	exampleString := "i love chicken nuggets"

	for i := 0; i < len(exampleString); i++ {
		// insert i-th char of example string on position
		position := 1 / float32(len(exampleString)+1) * float32(i+1)
		doc.Insert(position, exampleString[i])
	}

	if doc.GetText() != exampleString {
		t.Logf("Texts differ!\nGot: %s\nExpected: %s\n", doc.GetText(), exampleString)
		t.Fail()
	}

	insertWord := "really "
	startPosition := 1 / float32(len(exampleString)+1) * 2.0
	for i := 1; i < len(insertWord)+1; i++ {
		position := startPosition + float32(i)*0.000012
		doc.Insert(position, insertWord[i-1])
	}

	if doc.GetText() != "i really love chicken nuggets" {
		t.Logf("Texts differ!\nGot: %s\nExpected: %s\n", doc.GetText(), "i really love chicken nuggets")
		t.Fail()
	}
}

func TestDocument_Delete(t *testing.T) {
	doc := CreateDocument()
	doc.Insert(0.2, 's')
	doc.Insert(0.4, 'i')
	doc.Insert(0.6, 'g')
	doc.Insert(0.8, 'm')
	doc.Insert(0.9, 'a')

	doc.Delete(0.4)
	if doc.GetText() != "sgma" {
		t.Fail()
	}

	doc.Delete(0.2)
	if doc.GetText() != "gma" {
		t.Fail()
	}

	doc.Insert(0.92, 'i')
	doc.Insert(0.95, 'l')
	if doc.GetText() != "gmail" {
		t.Fail()
	}

	// try to delete special char
	doc.Delete(0)
	if !doc.data[0].isSpecialChar {
		t.Fail()
	}

	doc.Delete(1)
	if !doc.data[len(doc.data)-1].isSpecialChar {
		t.Fail()
	}
}

func TestDocument_PullData(t *testing.T) {
	doc := CreateDocument()
	exampleText := "oh no, cringe!"
	doc.Text = exampleText
	doc.PullData()

	if len(doc.data) != 2+len(exampleText) {
		t.Logf("Incorrect length of data array after pulling!\n")
		t.Fail()
	}

	for i, elem := range doc.data {
		if i == 0 || i == len(doc.data)-1 {
			if !elem.isSpecialChar {
				t.Fail()
			}
			continue
		}
		if elem.value != exampleText[i-1] {
			t.Logf("Data and example text differs at position %d\nExpected: %b\nGot: %b\n", i, exampleText[i-1], elem.value)
			t.Fail()
		}
	}
}
