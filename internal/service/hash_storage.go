package service

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"strings"
)

const StorageFile = "storage.txt"
const DELIMITER = ";"
const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type URLStorage struct {
}

type Record struct {
	HashURL string
	URL     string
}

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (s URLStorage) Save(URL string) (string, error) {

	key, err := genStorageKey(URL, false, 1)

	if err != nil {
		slog.Error("error generate storage key.", slog.Any("error", err))
		return "", fmt.Errorf("error generate storage key: %s", err)
	}

	// В хранилище не найден Hash, сохраняем его
	data := Record{
		HashURL: key,
		URL:     URL,
	}

	file, err := os.OpenFile(StorageFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return "", fmt.Errorf("error create file storage: %s", err)
	}
	defer file.Close()

	// Записываем поля структуры, разделяя их ; (%s;%s)
	_, err = fmt.Fprintf(file, "%s%s%s\n", data.HashURL, DELIMITER, data.URL)
	if err != nil {
		slog.Error("error write to file storage:. Error: ", slog.Any("error", err))
		return "", fmt.Errorf("error write to file storage: %s", err)
	}

	return key, nil
}

func (s URLStorage) Find(key string) (string, error) {
	url, err := findByKey(key)

	if err != nil {
		slog.Error("error find record by key. Error: ", slog.Any("error", err))
		return "", err
	}

	if url == "" {
		slog.Warn("record exist, but url empty", slog.Any("key", key))
		return "", fmt.Errorf("record exist, but url empty.")
	}

	return url, nil
}

func findByKey(key string) (string, error) {

	file, err := os.Open(StorageFile)

	if err != nil {
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

func genStorageKey(URL string, regenKey bool, countRandLetter int) (string, error) {
	hash := md5.Sum([]byte(URL))
	Hash := hex.EncodeToString(hash[:])

	if regenKey {
		Hash += randString(countRandLetter + 1)
	}

	result, err := findByKey(Hash)

	if err != nil {
		slog.Error("error find record by key. Error: ", slog.Any("error", err))
		return "", fmt.Errorf("error find record by key.: %s", err)
	}

	if result != "" {
		slog.Warn("key alreatdy exist. Go regenerate key.", slog.Any("key", Hash))
		return genStorageKey(URL, true, 1)
	}

	return Hash, nil
}
