package services

import (
	"notes_project/models"
)

type NoteStore interface {
	CreateNote(note *models.Note) error
	GetAllNotes() ([]models.Note, error)
	GetNoteByID(id uint) (*models.Note, error)
	UpdateNote(note *models.Note) error
	DeleteNote(id uint) error
}

type NoteService struct {
	storage NoteStore //шаг сделать на интерфейсе
}

func NewNoteService(storage NoteStore) *NoteService {
	return &NoteService{storage: storage}

}

func (s *NoteService) NewNote(title, text string) *models.Note { // конструкторы начинаются с New +
	return &models.Note{
		Title: title,
		Text:  text,
	}

}

func (s *NoteService) GetAllNotes() ([]models.Note, error) {
	return s.storage.GetAllNotes()
}

func (s *NoteService) GetNoteByID(id uint) (*models.Note, error) {
	return s.storage.GetNoteByID(id)
}

func (s *NoteService) AddNote(note *models.Note) error {
	return s.storage.CreateNote(note)
}

func (s *NoteService) UpdateNote(id uint, title, text string) (*models.Note, error) {
	note, err := s.storage.GetNoteByID(id)

	if err != nil {
		return nil, err
	}
	note.Title = title
	note.Text = text

	if err := s.storage.UpdateNote(note); err != nil {
		return nil, err
	}

	return note, nil

}

func (s *NoteService) DeleteNote(id uint) error {
	return s.storage.DeleteNote(id)
}
