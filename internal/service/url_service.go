package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/rakot9/ya-practicum-shortener-url/internal/model"
	"github.com/rakot9/ya-practicum-shortener-url/internal/repository"
	"github.com/rakot9/ya-practicum-shortener-url/internal/service/txctx"
	"go.uber.org/zap"
	"math/big"
	"net/url"
)

const base62Alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var ErrCountAttemptExceed = errors.New("count attempt exceed")
var ErrOriginalURLAlreadyExist = errors.New("original URL already processed")
var ErrInternalServerError = errors.New("internal server error")
var ErrServerError = errors.New("server error")
var ErrBadRequest = errors.New("bad request")

type BatchError struct {
	IsOrigURLInStorage bool
	Err                error
}

type URLRepository interface {
	Save(ctx context.Context, key string, url string, correlationID string) (string, error)
	Find(ctx context.Context, key string) (string, error)
	FindByOriginalURL(ctx context.Context, URL string) (string, error)
	Close() error
	GetLastUserId() (int, error)
}

// URLService управляет бизнес-логикой генерации ключей и координирует работу с репозиторием.
type URLService struct {
	repo URLRepository
	Txm  *txctx.TxManager
}

// NewURLService принимает любой репозиторий, удовлетворяющий интерфейсу URLRepository.
func NewURLService(repo URLRepository, txm *txctx.TxManager) *URLService {
	return &URLService{
		repo: repo,
		Txm:  txm,
	}
}

// ShortenURL генерирует уникальный ключ и сохраняет его в репозиторий.
func (s *URLService) ShortenURL(ctx context.Context, url string, correlationID string) (string, error) {
	for i := 10; i > 0; i-- {

		shortURL, err := s.repo.FindByOriginalURL(ctx, url)

		if err != nil {
			if !errors.Is(err, repository.ErrNotFound) {
				return "", fmt.Errorf("%w", err)
			}
		}

		if shortURL != "" {
			return shortURL, fmt.Errorf("%w", ErrOriginalURLAlreadyExist)
		}

		key, err := s.generateKey()
		if err != nil {
			return "", fmt.Errorf("failed to generate random id: %w", err)
		}

		// Проверяем, существует ли уже ключ в репозитории
		_, err = s.repo.Find(ctx, key)
		if errors.Is(err, repository.ErrNotFound) {
			// Ключ свободен, сохраняем и выходим из цикла
			if key, err = s.repo.Save(ctx, key, url, correlationID); err != nil {
				return "", fmt.Errorf("failed to save url: %w", err)
			}
			return key, nil
		}
		if err != nil {
			// Если произошла другая ошибка (например, упала БД), возвращаем её
			return "", fmt.Errorf("failed to check key existence: %w", err)
		}
	}

	return "", fmt.Errorf("failed generate key 10 attempt exceed for url %s error %w", url, ErrCountAttemptExceed)
}

func (s *URLService) ShortenBatch(ctx context.Context, batchRequestItems []model.BatchRequestItem, logger *zap.Logger, domain string) ([]model.BatchResponseItem, error) {

	var BatchURLsItems []model.BatchResponseItem

	isSetOriginalURLInStorage := false

	err := s.Txm.RunInTx(ctx, logger, func(ctx context.Context) error {

		var err error

		for _, item := range batchRequestItems {
			var key string

			key, err = s.ShortenURL(ctx, item.OriginalURL, item.CorrelationID)

			if err != nil {
				if errors.Is(err, ErrOriginalURLAlreadyExist) {
					isSetOriginalURLInStorage = true
				}
				if errors.Is(err, ErrCountAttemptExceed) {
					logger.Error("count attempt (10) exceed for generate key. URL", zap.Error(err))
					return fmt.Errorf("count attempt (10) exceed for generate key. URL: %w", ErrServerError)
				}
				if !isSetOriginalURLInStorage {
					return fmt.Errorf("internal error: %w", ErrInternalServerError)
				}
			}

			url, err := url.JoinPath(domain, key)

			if err != nil {
				return fmt.Errorf("join short url: %w", ErrBadRequest)
			}

			var responseItems = model.BatchResponseItem{
				CorrelationID: item.CorrelationID,
				ShortURL:      url,
			}

			BatchURLsItems = append(BatchURLsItems, responseItems)
		}

		return err
	})

	if err != nil {
		err = &BatchError{
			IsOrigURLInStorage: isSetOriginalURLInStorage,
			Err:                err,
		}
	}

	return BatchURLsItems, err
}

// GetURL возвращает оригинальный URL по его короткому ключу.
func (s *URLService) GetURL(ctx context.Context, key string) (string, error) {
	return s.repo.Find(ctx, key)
}

func (s *URLService) GetUser() (int, error) {
	return s.repo.GetLastUserId()
}

// Вспомогательный метод генерации ключа остался внутри сервиса, так как это бизнес-логика.
func (s *URLService) generateKey() (string, error) {
	lengthBig, err := rand.Int(rand.Reader, big.NewInt(4)) // случайное число от 0 до 3
	if err != nil {
		return "", err
	}
	length := int(lengthBig.Int64()) + 6

	b := make([]byte, length)
	alphabetLen := big.NewInt(int64(len(base62Alphabet)))

	for i := 0; i < length; i++ {
		idx, err := rand.Int(rand.Reader, alphabetLen)
		if err != nil {
			return "", err
		}
		b[i] = base62Alphabet[idx.Int64()]
	}

	return string(b), nil
}

func (e *BatchError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return fmt.Sprintf("batch error: %v", e.Err)
}

func (e *BatchError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}
