package store

import (
	"fmt"
	"notes_project/models"
	"os"
	"sync"

	"github.com/gocarina/gocsv"
)

type CSVStore struct {
	FilePath string
	LastID   int
	mu       sync.Mutex
}

func NewCSVStore(path string) *CSVStore {
	return &CSVStore{FilePath: path}
}

func (s *CSVStore) GetAllNotes() ([]models.Note, error) {
	file, err := os.Open(s.FilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var notes []models.Note
	if err := gocsv.UnmarshalFile(file, &notes); err != nil {
		return nil, err
	}

	return notes, nil
}
func (s *CSVStore) CreateNote(note *models.Note) error {
	notes, err := s.GetAllNotes()
	if err != nil {
		notes = []models.Note{}
	}

	var maxID uint
	for _, n := range notes {
		if n.ID > maxID {
			maxID = n.ID
		}
	}

	note.ID = maxID + 1

	notes = append(notes, *note)

	return s.SaveAll(notes)
}

func (s *CSVStore) GetNoteByID(id uint) (*models.Note, error) {
	notes, err := s.GetAllNotes()
	if err != nil {
		return nil, err
	}

	for _, note := range notes {
		if note.ID == id {
			return &note, nil
		}
	}

	return nil, fmt.Errorf("note not found")
}

func (s *CSVStore) SaveAll(notes []models.Note) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	file, err := os.Create(s.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return gocsv.MarshalFile(&notes, file)
}
func (s *CSVStore) UpdateNote(note *models.Note) error {
	notes, err := s.GetAllNotes()
	if err != nil {
		return err
	}

	found := false

	for i := range notes {
		if notes[i].ID == note.ID {
			notes[i] = *note
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("note not found")
	}

	return s.SaveAll(notes)
}

func (s *CSVStore) DeleteNote(id uint) error {
	notes, err := s.GetAllNotes()
	if err != nil {
		return err
	}

	var result []models.Note

	for _, note := range notes {
		if note.ID != id {
			result = append(result, note)
		}
	}

	return s.SaveAll(result)
}
