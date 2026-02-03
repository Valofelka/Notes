package services

import (
	"encoding/csv"
	"fmt"
	"notes_project/models"
	"os"
	"strconv"
	"time"

	"github.com/gocarina/gocsv"
)

type NoteService struct {
	filePath string
	lastId   int
}

func NewNoteService(filePath string) *NoteService {
	return &NoteService{filePath: filePath}

}

func (s *NoteService) CreateNote(title, text string) *models.Note {
	return &models.Note{
		Title:     title,
		Text:      text,
		CreatedAt: time.Now(),
	}

}

func (s *NoteService) LastID() error {
	file, err := os.Open(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			s.lastId = 0
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
	s.lastId = maxId
	return nil

}

func (s *NoteService) nextID() int {
	s.lastId++
	return s.lastId
}

func (s *NoteService) AddNote(note *models.Note) error {
	notes, err := s.GetAllNotes()
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	note.Id = s.nextID()
	notes = append(notes, note)

	file, err := os.Create(s.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return gocsv.MarshalFile(&notes, file)
}

func (s *NoteService) DeleteNote(id int) error {
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
		result = append(result, note)

	}
	if !found {
		return fmt.Errorf("not found id")
	}

	file, err := os.Create(s.filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	return gocsv.MarshalFile(&result, file)

}

func (s *NoteService) UpdateNote(id int, title, text string) (*models.Note, error) {
	file, err := os.Open(s.filePath)
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

	newFile, err := os.Create(s.filePath)
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

func (s *NoteService) GetAllNotes() ([]*models.Note, error) {
	file, err := os.Open(s.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var notes []*models.Note
	if err := gocsv.UnmarshalFile(file, &notes); err != nil {
		return nil, err
	}
	return notes, nil
}

func (s *NoteService) GetNoteByID(id int) (*models.Note, error) {
	notes, err := s.GetAllNotes()
	if err != nil {
		return nil, err
	}

	for _, note := range notes {
		if note.Id == id {
			return note, nil
		}
	}
	return nil, fmt.Errorf("not found id")
}

// test
