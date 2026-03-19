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

func (s *CSVStore) UpdateNote(id int, title, text string) (*models.Note, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	file, err := os.Open(s.FilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var notes []*models.Note // использовать гоцсв анмаршал +
	err = gocsv.Unmarshal(file, &notes)
	if err != nil {
		return nil, err
	}

	var updatedNote *models.Note
	var updated bool

	for _, n := range notes {
		if n.Id == id {
			n.Title = title
			n.Text = text
			updatedNote = n
			updated = true
		}
	}
	if !updated {
		return nil, fmt.Errorf("note with id %d not found", id)
	}

	tmpFile := s.FilePath + ".tmp"

	newFile, err := os.Create(tmpFile)
	if err != nil {
		return nil, err
	}

	err = gocsv.Marshal(&notes, newFile)

	newFile.Close()

	if err != nil {
		return nil, err
	}

	err = os.Rename(tmpFile, s.FilePath)

	return updatedNote, nil
}

func (s *CSVStore) DeleteNote(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	notes, err := s.GetAllNotes()
	if err != nil {
		return err
	}
	var result []models.Note
	found := false

	for _, note := range notes {
		if note.Id != id {
			result = append(result, note)

		} else {
			found = true
		}

	}
	if !found {
		return fmt.Errorf("not found id")
	}

	tmpFile := s.FilePath + ".tmp"

	file, err := os.Create(tmpFile)
	if err != nil {
		return err
	}

	err = gocsv.MarshalFile(result, file)
	file.Close()
	if err != nil {
		return err
	}
	return os.Rename(tmpFile, s.FilePath)

}
