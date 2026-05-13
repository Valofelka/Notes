package store

import (
	"notes_project/models"

	"gorm.io/gorm"
)

type PostgresStore struct {
	db *gorm.DB
}

func NewPostgresStore(db *gorm.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func (s *PostgresStore) CreateNote(note *models.Note) error {
	result := s.db.Create(note)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *PostgresStore) GetAllNotes() ([]models.Note, error) {
	var notes []models.Note

	result := s.db.Find(&notes)

	if result.Error != nil {
		return nil, result.Error
	}
	return notes, nil
}

func (s *PostgresStore) UpdateNote(note *models.Note) error {
	result := s.db.Model(note).Updates(note)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *PostgresStore) DeleteNote(id uint) error {
	result := s.db.Delete(&models.Note{}, id)

	if result.Error != nil {
		return result.Error
	}
	return nil
}
