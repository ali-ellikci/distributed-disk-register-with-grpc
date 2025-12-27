package storage

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const baseDir = "messages"

func WriteMessage(id int, msg string) error {
	// messages klasörünü oluştur
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return err
	}

	// id.int → string
	path := filepath.Join(baseDir, fmt.Sprintf("%d.msg", id))

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(msg)
	if err != nil {
		return err
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	return file.Sync()
}

func ReadMessage(id int) (string, error) {
	path := filepath.Join(baseDir, fmt.Sprintf("%d.msg", id))

	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	msgBytes, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(msgBytes), nil
}
