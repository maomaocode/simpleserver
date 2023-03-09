package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Saver struct {
	path        string
	currentFile string
	file        *os.File
}

func NewSaver(path string) *Saver {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.MkdirAll(path, os.ModePerm)
	}
	return &Saver{path: path}
}

func (s *Saver) GetFile() (*os.File, error) {
	fileName := fmt.Sprintf("%s/%s.log", s.path, time.Now().Format("2006-01-02"))

	if s.currentFile == fileName {
		return s.file, nil
	}

	s.currentFile = fileName

	_ = s.file.Close()

	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return nil, err
	}

	s.file = f

	return s.file, nil
}

func (s *Saver) write(content string) {
	f, err := s.GetFile()
	if err != nil {
		log.Println(err)
		return
	}

	if _, err = f.WriteString(content); err != nil {
		log.Println(err)
	}

}

func (s *Saver) Close() {
	_ = s.file.Close()
}
