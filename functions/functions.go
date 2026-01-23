package functions

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type Note struct {
	Text     string
	Id       int
	Title    string
	CreateAt time.Time
}

func Create(note *Note, title, text string) {
	note.Title = title
	note.Text = text
	note.CreateAt = time.Now()

}

func ChangeID() (int, error) {
	file, err := os.Open("notes.csv")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	maxId := 0
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return 0, err
	}
	for i, record := range records {
		if i == 0 {
			continue
		}
		id, err := strconv.Atoi(record[0])
		if err != nil {
			continue
		}
		if id > maxId {
			maxId = id
		}

	}
	return maxId + 1, nil
}

func AddNote(note *Note) error {
	file, err := os.OpenFile("notes.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		return err
	}
	defer file.Close()
	id, err := ChangeID()
	if err != nil {
		return err
	}
	note.Id = id
	idStr := strconv.Itoa(note.Id)
	createAdStr := note.CreateAt.Format(time.RFC1123Z)
	record := []string{
		idStr,
		note.Title,
		note.Text,
		createAdStr,
	}

	writer := csv.NewWriter(file)

	if err := writer.Write(record); err != nil {
		log.Fatalf("Failed to transfer data: %v", err)
	}
	writer.Flush()
	return nil
}

func DeleteNote(note *Note) error {
	file, err := os.Open("notes.csv")
	if err != nil {
		log.Fatalf("Reading error: %v", err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	var newRecords [][]string
	var flagId bool
	for i, record := range records {
		if i == 0 {
			newRecords = append(newRecords, record)
			continue
		}

		if record[0] == strconv.Itoa(note.Id) {
			flagId = true
			continue
		}
		newRecords = append(newRecords, record)
	}

	wrFile, err := os.Create("notes.csv")
	if err != nil {
		log.Fatalf("Create error: %v", err)
	}
	defer wrFile.Close()

	if !flagId {
		return fmt.Errorf("note with id %d not found", note.Id)
	}

	writer := csv.NewWriter(file)
	err = writer.WriteAll(newRecords)
	if err != nil {
		return err
	}
	writer.Flush()

	return nil
}

func UpdateNote(note *Note) error {
	file, err := os.Open("notes.csv")
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

	newFile, err := os.Create("notes.csv")
	if err != nil {
		return err
	}
	defer newFile.Close()

	writer := csv.NewWriter(newFile)
	err = writer.WriteAll(newRecords)
	if err != nil {
		return err
	}
	writer.Flush()

	return nil
}

func ReadNote() ([][]string, error) {
	file, err := os.Open("notes.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	return reader.ReadAll()
}

func ReadNoteId(id int) (*Note, error) {
	file, err := os.Open("notes.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	for i, record := range records {
		if i == 0 {
			continue
		}
		recordId, err := strconv.Atoi(record[0])
		if err != nil {
			continue
		}

		if recordId == id {
			createdAt, _ := time.Parse(time.RFC1123Z, record[3])

			return &Note{
				Id:       recordId,
				Title:    record[1],
				Text:     record[2],
				CreateAt: createdAt,
			}, nil
		}
	}
	return nil, fmt.Errorf("Not found id")
}
