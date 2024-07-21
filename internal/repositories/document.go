package repositories

import (
	"github.com/masterkusok/websocketCollab/internal/businnesLogic"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewDocumentRepository(db *gorm.DB) businnesLogic.DocumentRepository {
	return &Repository{db: db}
}

func (r *Repository) GetById(id int) (*businnesLogic.Document, error) {
	var doc businnesLogic.Document
	context := r.db.First(&doc, id)
	if context.Error != nil {
		return nil, context.Error
	}
	doc.PullData()
	return &doc, nil
}

func (r *Repository) CreateDocument() (*businnesLogic.Document, error) {
	doc := businnesLogic.CreateDocument()
	context := r.db.Create(doc)
	if context.Error != nil {
		return nil, context.Error
	}
	return doc, nil
}

func (r *Repository) UpdateDocument(doc *businnesLogic.Document) error {
	context := r.db.Save(doc)
	if context.Error != nil {
		return context.Error
	}
	return nil
}
