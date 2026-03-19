package services

import (
	"fmt"
	"notes_project/models"
	"time"
)

type NoteStore interface {
	GetAllNotes() ([]models.Note, error)
	SaveAll(notes []models.Note) error
	UpdateNote(id int, title, text string) (*models.Note, error)
	DeleteNote(id int) error
}

type NoteService struct {
	storage NoteStore //шаг сделать на интерфейсе
}

func NewNoteService(storage NoteStore) *NoteService {
	return &NoteService{storage: storage}

}

func (s *NoteService) NewNote(title, text string) *models.Note { // конструкторы начинаются с New +
	return &models.Note{
		Title:     title,
		Text:      text,
		CreatedAt: time.Now(),
	}

}

func (s *NoteService) GetAllNotes() ([]models.Note, error) {
	return s.storage.GetAllNotes()
}

func (s *NoteService) nextID(notes []models.Note) int {
	maxID := 0
	for _, note := range notes {
		if note.Id > maxID {
			maxID = note.Id
		}
	}
	return maxID + 1
}

func (s *NoteService) GetNoteByID(id int) (*models.Note, error) {
	notes, err := s.storage.GetAllNotes()
	if err != nil {
		return nil, err
	}

	for _, note := range notes {
		if note.Id == id {
			return &note, nil
		}
	}
	return nil, fmt.Errorf("not found id")
}

func (s *NoteService) AddNote(note *models.Note) error {
	notes, err := s.storage.GetAllNotes()
	if err != nil {
		notes = []models.Note{}
	}

	note.Id = s.nextID(notes)
	notes = append(notes, *note)

	return s.storage.SaveAll(notes)
}

func (s *NoteService) UpdateNote(id int, title, text string) (*models.Note, error) {
	notes, err := s.storage.GetAllNotes()

	var updatedNote *models.Note

	if err != nil {
		return updatedNote, err
	}
	found := false

	for i := range notes {
		if notes[i].Id == id {
			notes[i].Title = title
			notes[i].Text = text
			updatedNote = &notes[i]
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("note not found")
	}

	if err := s.storage.SaveAll(notes); err != nil {
		return nil, err
	}

	return updatedNote, nil

}

func (s *NoteService) DeleteNote(id int) error {
	return s.storage.DeleteNote(id)
}

//test
