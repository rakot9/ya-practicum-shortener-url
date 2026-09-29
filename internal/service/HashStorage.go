package service

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

const STORAGE_FILE = "storage.txt"
const DELIMITER = ";"

type Record struct {
	HashUrl string
	Url     string
}

// Todo: сделать проверку на дубли
func Save(Url string) (string, error) {

	hash := md5.Sum([]byte(Url))
	Hash := hex.EncodeToString(hash[:])

	data := Record{
		HashUrl: Hash,
		Url:     Url,
	}

	file, err := os.OpenFile(STORAGE_FILE, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return "", fmt.Errorf("Ошибка создания файла хранилища: %s", err)
	}
	defer file.Close()

	// Записываем поля структуры, разделяя их ; (%s;%s)
	_, err = fmt.Fprintf(file, "%s%s%s\n", data.HashUrl, DELIMITER, data.Url)
	if err != nil {
		fmt.Println("Ошибка записи в файл хранилища:", err)
		return "", fmt.Errorf("Ошибка записи в файл хранилища: %s", err)
	}

	return Hash, nil
}

func Find(hash string) (string, error) {

	file, err := os.OpenFile(STORAGE_FILE, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

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
		return "", err
	}

	return "", fmt.Errorf("ключ %s не найден", hash)
}
