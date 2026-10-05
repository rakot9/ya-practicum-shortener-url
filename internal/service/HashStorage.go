package service

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

const StorageFile = "storage.txt"
const DELIMITER = ";"

type Record struct {
	HashURL string
	URL     string
}

// Todo: сделать проверку на дубли
func Save(URL string) (string, error) {

	hash := md5.Sum([]byte(URL))
	Hash := hex.EncodeToString(hash[:])

	data := Record{
		HashURL: Hash,
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
		fmt.Println("Error write to file storage:", err)
		return "", fmt.Errorf("error write to file storage: %s", err)
	}

	return Hash, nil
}

func Find(hash string) (string, error) {

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
		if currentKey == hash {
			return parts[1], nil
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Print("Error scanner")
		return "", err
	}

	return "", fmt.Errorf("key %s not found", hash)
}
