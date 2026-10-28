package repository

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/rakot9/ya-practicum-shortener-url/internal/config"
)

type FileStorage struct {
	UUID        int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// FileStorgeRepository реализует URLRepository для хранения данных в оперативной памяти.
type FileStorgeRepository struct {
	configParams config.Flags
	file         *os.File
}

func NewFileStorgeRepository(configParams config.Flags) (*FileStorgeRepository, error) {

	file, err := os.OpenFile(configParams.FlagFileStoragePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open storage file: %w", err)
	}

	return &FileStorgeRepository{
		configParams: configParams,
		file:         file,
	}, nil
}

func (r *FileStorgeRepository) Save(key string, URL string) error {
	lines, err := r.getMaxCountRecords()

	if err != nil {
		return err
	}

	saveData := FileStorage{
		UUID:        lines + 1,
		ShortURL:    key,
		OriginalURL: URL,
	}

	data, err := json.Marshal(saveData)
	if err != nil {
		return err
	}

	data = append(data, '\n')

	_, err = r.file.Write(data)

	return err
}

func (r *FileStorgeRepository) Find(key string) (string, error) {
	if _, err := r.file.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(r.file)

	lineNum := 0

	for scanner.Scan() {
		lineNum++

		lineBytes := scanner.Bytes()

		var fileStorage FileStorage

		err := json.Unmarshal(lineBytes, &fileStorage)

		if err != nil {
			continue
		}

		if fileStorage.ShortURL == key {
			return fileStorage.OriginalURL, nil
		}

	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("key not found in storage file %s: %w", key, ErrNotFound)
}

func (r *FileStorgeRepository) getMaxCountRecords() (int, error) {
	if _, err := r.file.Seek(0, io.SeekStart); err != nil {
		return 0, err
	}

	scanner := bufio.NewScanner(r.file)
	count := 0

	// Сканер автоматически читает файл строка за строкой
	for scanner.Scan() {
		count++
	}

	// Обязательно проверяем, не возникло ли ошибок при чтении
	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *FileStorgeRepository) Close() error {
	if r.file != nil {
		return r.file.Close()
	}
	return nil
}
