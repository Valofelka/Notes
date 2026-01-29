package services

import (
	"encoding/csv"
	"fmt"
	"github.com/gocarina/gocsv"
	"notes_project/models"
	"os"
	"strconv"
	"time"
)

type NoteService struct {
	filePath string
	lastId   int
}

func NewNoteService(filePath string) (*NoteService, error) {
	service := &NoteService{filePath: filePath}

	if err := service.LastID(); err != nil {
		return nil, err
	}
	return service, nil
}

func (s *NoteService) Create(title, text string) *models.Note {
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

func (s *NoteService) UpdateNote(note *models.Note) error {
	file, err := os.Open(s.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	var newRecords [][]string
	var _ bool

	for i, record := range records {
		if i == 0 {
			newRecords = append(newRecords, record)
			continue
		}
		if record[0] == strconv.Itoa(note.Id) {
			record[1] = note.Title
			record[2] = note.Text
			_ = true
		}
		newRecords = append(newRecords, record)
	}

	newFile, err := os.Create(s.filePath)
	if err != nil {
		return err
	}
	defer newFile.Close()

	writer := csv.NewWriter(newFile)
	err = writer.WriteAll(newRecords)
	if err != nil {
		return err
	}
	defer writer.Flush()

	return nil
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
