package service

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

const StorageFile = "storage.txt"
const DELIMITER = ";"

type URLStorage struct {
}

type Record struct {
	HashURL string
	URL     string
}

func (s URLStorage) Save(URL string, key string) (bool, error) {

	findResult, err := s.Find(key)

	if findResult != "" {
		return false, nil
	}

	// В хранилище не найден Hash, сохраняем его
	data := Record{
		HashURL: key,
		URL:     URL,
	}

	file, err := os.OpenFile(StorageFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return false, fmt.Errorf("error create file storage: %s", err)
	}
	defer file.Close()

	// Записываем поля структуры, разделяя их ; (%s;%s)
	_, err = fmt.Fprintf(file, "%s%s%s\n", data.HashURL, DELIMITER, data.URL)
	if err != nil {
		slog.Error("error write to file storage:. Error: ", slog.Any("error", err))
		return false, fmt.Errorf("error write to file storage: %s", err)
	}

	return true, nil
}

func (s URLStorage) Find(key string) (string, error) {
	url, err := findByKey(key)

	if err != nil {
		slog.Error("error find record by key. Error: ", slog.Any("error", err))
		return "", err
	}

	if url == "" {
		slog.Warn("record exist, but url empty", slog.Any("key", key))
		return "", fmt.Errorf("record exist, but url empty with key %s", key)
	}

	return url, nil
}

func findByKey(key string) (string, error) {
	file, err := os.OpenFile(StorageFile, os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		slog.Error("findByKey error open file.", slog.Any("error", err))
		return "", err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Split(line, DELIMITER)

		if len(parts) == 0 {
			continue
		}

		currentKey := strings.TrimSpace(parts[0])
		if currentKey == key {
			return parts[1], nil
		}
	}

	if err := scanner.Err(); err != nil {
		slog.Error("error scanner", slog.Any("error", err))
		return "", err
	}

	return "", nil
}
