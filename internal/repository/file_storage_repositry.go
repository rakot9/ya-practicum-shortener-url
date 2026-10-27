package repository

import (
	"bufio"
	"encoding/json"
	"fmt"
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
}

func NewFileStorgeRepository(configParams config.Flags) *FileStorgeRepository {
	return &FileStorgeRepository{
		configParams: configParams,
	}
}

func (r *FileStorgeRepository) Save(key string, URL string) error {
	file, err := os.OpenFile(r.configParams.FlagFileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

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

	_, err = file.Write(data)

	return err
}

func (r *FileStorgeRepository) Find(key string) (string, error) {
	fmt.Printf("File storage Find %s \n", key)

	file, err := os.OpenFile(r.configParams.FlagFileStoragePath, os.O_RDONLY, 0666)

	if err != nil {
		return "", fmt.Errorf("key not found in storage file %s: %w", key, ErrNotFound)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	lineNum := 0

	for scanner.Scan() {
		lineNum++

		lineBytes := scanner.Bytes()

		var fileStorage FileStorage

		err := json.Unmarshal(lineBytes, &fileStorage)

		if err != nil {
			fmt.Printf("Ошибка парсинга на строке %d: %v", lineNum, err)
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
	file, err := os.OpenFile(r.configParams.FlagFileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return 0, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
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
