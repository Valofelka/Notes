package services

import (
	"encoding/csv"
	"fmt"
	"notes_project/models"
	"notes_project/services/store"
	"os"
	"strconv"
	"time"

	"github.com/gocarina/gocsv"
)

type NoteService struct {
	storage *store.CSVStore
}

func NewNoteService(storage *store.CSVStore) *NoteService {
	return &NoteService{storage: storage}

}

func (s *NoteService) CreateNote(title, text string) *models.Note {
	return &models.Note{
		Title:     title,
		Text:      text,
		CreatedAt: time.Now(),
	}

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

	file, err := os.Create(s.storage.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return gocsv.MarshalFile(&notes, file)
}

func (s *NoteService) LastID() error {
	file, err := os.Open(s.storage.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			s.storage.LastID = 0
			return nil
		}
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	maxId := 0
	for i, record := range records {
		if i == 0 {
			continue
		}
		id, err := strconv.Atoi(record[0])
		if err != nil {
			return fmt.Errorf("invalid id")
		}
		if id > maxId {
			maxId = id
		}
	}
	s.storage.LastID = maxId
	return nil

}

func (s *NoteService) GetAllNotes() ([]models.Note, error) {
	return s.storage.GetAllNotes()
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
	notes, err := s.storage.GetAllNotes()
	if err != nil {
		return err
	}

	var newNotes []models.Note
	found := false

	for _, note := range notes {
		if note.Id != id {
			newNotes = append(newNotes, note)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("note not found")
	}

	return s.storage.SaveAll(newNotes)
}

//test
