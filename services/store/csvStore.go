package store

import (
	"encoding/csv"
	"fmt"
	"notes_project/models"
	"os"
	"strconv"
	"time"

	"github.com/gocarina/gocsv"
)

type CSVStore struct {
	FilePath string
	LastID   int
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
	file, err := os.Open(s.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return gocsv.MarshalFile(&notes, file)
}

func (s *CSVStore) UpdateNote(id int, title, text string) (*models.Note, error) {
	file, err := os.Open(s.FilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var (
		newRecords [][]string
		updated    bool
		updateNote *models.Note
	)

	for i, record := range records {
		if i == 0 {
			newRecords = append(newRecords, record)
			continue
		}
		recordID, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, err
		}
		if recordID == id {
			record[1] = title
			record[2] = text
			updated = true

			createdAt, _ := time.Parse(time.RFC1123Z, record[3])

			updateNote = &models.Note{
				Id:        id,
				Title:     title,
				Text:      text,
				CreatedAt: createdAt,
			}
		}
		newRecords = append(newRecords, record)

	}

	if !updated {
		return nil, fmt.Errorf("note with id %d not found", id)
	}

	newFile, err := os.Create(s.FilePath)
	if err != nil {
		return nil, err
	}
	defer newFile.Close()

	writer := csv.NewWriter(newFile)
	defer writer.Flush()

	if err := writer.WriteAll(newRecords); err != nil {
		return nil, err
	}

	return updateNote, nil
}

func (s *CSVStore) DeleteNote(id int) error {
	notes, err := s.GetAllNotes()
	if err != nil {
		return err
	}

	var result []*models.Note
	found := false

	for _, note := range notes {
		if note.Id == id {
			found = true
			continue
		}
		result = append(result, &note)

	}
	if !found {
		return fmt.Errorf("not found id")
	}

	file, err := os.Create(s.FilePath)
	if err != nil {
		return err
	}

	defer file.Close()

	return gocsv.MarshalFile(&result, file)

}
