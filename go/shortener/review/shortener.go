package main

import (
	"crypto/rand"
	"errors"
	"math/big"
	"sync"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var (
	ErrNotFound    = errors.New("short code not found")
	ErrExpired     = errors.New("short code expired")
	ErrCodeTaken   = errors.New("custom code already in use")
	ErrInvalidCode = errors.New("invalid code format")
	ErrInvalidURL  = errors.New("invalid URL")
)

type Stats struct {
	ShortCode   string
	LongURL     string
	CreatedAt   time.Time
	ExpiresAt   time.Time // Zero means never expires
	AccessCount int
}

func (s *Stats) isExpired() bool {
	return !s.ExpiresAt.IsZero() && time.Now().After(s.ExpiresAt)
}

type Shortener struct {
	mu        sync.RWMutex
	stats     map[string]*Stats // code -> stats
	urlToCode map[string]string // url -> code
}

func NewShortener() *Shortener {
	return &Shortener{
		stats:     make(map[string]*Stats),
		urlToCode: make(map[string]string),
	}
}

func (s *Shortener) generateCode(length int) (string, error) {
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[n.Int64()]
	}
	return string(b), nil
}

func (s *Shortener) generateUniqueCode(length int) (string, error) {
	for attempts := 0; attempts < 100; attempts++ {
		code, err := s.generateCode(length)
		if err != nil {
			return "", err
		}
		if _, exists := s.stats[code]; !exists {
			return code, nil
		}
	}
	return "", errors.New("failed to generate unique code")
}

func (s *Shortener) cleanup(code string) {
	if stat, exists := s.stats[code]; exists {
		delete(s.urlToCode, stat.LongURL)
		delete(s.stats, code)
	}
}

func (s *Shortener) shorten(url string, ttl time.Duration, customCode string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if existingCode, exists := s.urlToCode[url]; exists {
		stat := s.stats[existingCode]
		if !stat.isExpired() {
			return existingCode, nil
		}
		s.cleanup(existingCode)
	}

	var code string
	var err error

	if customCode != "" {
		if _, exists := s.stats[customCode]; exists {
			return "", ErrCodeTaken
		}
		code = customCode
	} else {
		code, err = s.generateUniqueCode(6)
		if err != nil {
			return "", err
		}
	}

	var expiresAt time.Time
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl)
	}

	s.stats[code] = &Stats{
		ShortCode:   code,
		LongURL:     url,
		CreatedAt:   time.Now(),
		ExpiresAt:   expiresAt,
		AccessCount: 0,
	}
	s.urlToCode[url] = code

	return code, nil
}

func (s *Shortener) Shorten(url string) (string, error) {
	return s.shorten(url, 0, "")
}

func (s *Shortener) ShortenWithTTL(url string, ttl time.Duration) (string, error) {
	return s.shorten(url, ttl, "")
}

func (s *Shortener) ShortenCustom(url string, customCode string) (string, error) {
	return s.shorten(url, 0, customCode)
}

func (s *Shortener) Resolve(code string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stat, exists := s.stats[code]
	if !exists {
		return "", ErrNotFound
	}

	if stat.isExpired() {
		s.cleanup(code)
		return "", ErrExpired
	}

	stat.AccessCount++
	return stat.LongURL, nil
}

func (s *Shortener) GetStats(code string) (*Stats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stat, exists := s.stats[code]
	if !exists {
		return nil, ErrNotFound
	}

	if stat.isExpired() {
		return nil, ErrExpired
	}

	copy := *stat
	return &copy, nil
}
