package businnesLogic

type DocumentRepository interface {
	GetById(id int) (*Document, error)
	CreateDocument() (*Document, error)
	UpdateDocument(document *Document) error
}
