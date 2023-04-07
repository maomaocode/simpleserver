package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Saver struct {
	path        string
	currentFile string
	file        *os.File

	addressSet map[string]struct{}
}

func NewSaver(path string) *Saver {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.MkdirAll(path, os.ModePerm)
	}

	addressSet := make(map[string]struct{})

	filepath.Walk(path, func(file string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		data, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		for _, row := range strings.Split(string(data), "\n") {
			r := UploadFormReq{}
			if err := json.Unmarshal([]byte(row), &r); err != nil {
				continue
			}

			addressSet[r.ID] = struct{}{}
		}

		return nil
	})

	return &Saver{path: path, addressSet: addressSet}
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

func (s *Saver) write(addr, content string) {
	f, err := s.GetFile()
	if err != nil {
		log.Println(err)
		return
	}

	s.addressSet[addr] = struct{}{}

	if _, err = f.WriteString(content); err != nil {
		log.Println(err)
	}
	f.WriteString("\n")
}

func (s *Saver) Close() {
	_ = s.file.Close()
}

func (s *Saver) Exist(addr string) bool {
	_, ok := s.addressSet[addr]
	return ok
}