package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Stats struct {
	ShortCode   string
	LongURL     string
	TTL         time.Duration
	CreatedAt   time.Time
	AccessCount int
}

type Shortener struct {
	shorten   map[string]*Stats
	urlToCode map[string]string
	codeToUrl map[string]string
}

var CHARSET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func NewShortener() *Shortener {
	return &Shortener{
		shorten:   make(map[string]*Stats),
		urlToCode: make(map[string]string),
		codeToUrl: make(map[string]string),
	}
}

func (s *Shortener) generateCode(length int) string {
	for {
		b := make([]byte, length)
		for i := range b {
			b[i] = CHARSET[seededRand.Intn(len(CHARSET))]
		}
		if _, exists := s.codeToUrl[string(b)]; !exists {
			return string(b)
		}
	}
}

func (s *Shortener) Shorten(url string) string {
	if val, exists := s.urlToCode[url]; exists {
		shorten, exists := s.shorten[val]
		if exists && shorten.TTL > 0 && time.Since(shorten.CreatedAt) > shorten.TTL {
			delete(s.shorten, val)
			delete(s.codeToUrl, val)
			delete(s.urlToCode, url)
		} else {
			return val
		}
	}

	code := s.generateCode(6)

	s.urlToCode[url] = code
	s.codeToUrl[code] = url

	s.shorten[code] = &Stats{
		ShortCode:   code,
		LongURL:     url,
		CreatedAt:   time.Now(),
		AccessCount: 0,
		TTL:         0, // unlimited
	}

	return code
}

func (s *Shortener) Resolve(code string) (string, error) {
	url, exists := s.codeToUrl[code]
	if !exists {
		return "", errors.New("short code not found")
	}

	stat, exists := s.shorten[code]
	if !exists {
		return "", errors.New("short code not found")
	}

	if stat.TTL > 0 && time.Since(stat.CreatedAt) > stat.TTL {
		delete(s.shorten, code)
		delete(s.codeToUrl, code)
		delete(s.urlToCode, url)

		return "", errors.New("short code not found")
	}

	stat.AccessCount += 1

	return url, nil
}

func (s *Shortener) GetStats(code string) *Stats {
	if stat, exists := s.shorten[code]; exists {
		if exists && stat.TTL > 0 && time.Since(stat.CreatedAt) > stat.TTL {
			delete(s.shorten, code)
			delete(s.codeToUrl, code)
			delete(s.urlToCode, stat.LongURL)
		} else {
			return stat
		}
	}
	return nil
}

func (s *Shortener) ShortenWithTTL(longURL string, ttl time.Duration) string {
	if val, exists := s.urlToCode[longURL]; exists {
		stat, exists := s.shorten[val]
		if exists && stat.TTL > 0 && time.Since(stat.CreatedAt) > stat.TTL {
			delete(s.shorten, val)
			delete(s.codeToUrl, val)
			delete(s.urlToCode, longURL)
		} else {
			return val
		}
	}

	code := s.generateCode(6)

	s.urlToCode[longURL] = code
	s.codeToUrl[code] = longURL

	s.shorten[code] = &Stats{
		ShortCode:   code,
		LongURL:     longURL,
		CreatedAt:   time.Now(),
		AccessCount: 0,
		TTL:         ttl,
	}

	return code
}

func main() {
	sh := NewShortener()
	short1 := sh.Shorten("https://www.google.com/search?q=golang+tutorials")
	fmt.Println(short1)
	short2 := sh.Shorten("https://www.google.com/search?q=golang+tutorials")
	fmt.Println(short2)
	short3 := sh.Shorten("https://www.google.com/search?q=golang+tutorials")
	fmt.Println(short3)

	sh.Resolve(short1)
	sh.Resolve(short1)
	sh.Resolve(short1)

	stat := sh.GetStats(short1)

	fmt.Printf("Stats : %v \n", stat)

	message, err := sh.Resolve("invalid")
	fmt.Printf("Invalid : %s, %e \n", message, err)
}

// ```

// **Example:**
// ```
// short1 := Shorten("https://www.google.com/search?q=golang+tutorials")
// // returns "xK9mQ2"

// short2 := Shorten("https://www.google.com/search?q=golang+tutorials")
// // returns "xK9mQ2" (same URL = same code)

// short3 := Shorten("https://github.com")
// // returns "pL3nR8" (different URL = different code)

// Resolve("xK9mQ2")
// // returns "https://www.google.com/search?q=golang+tutorials", nil

// Resolve("xK9mQ2")
// // returns same URL (and increments access count)

// Resolve("invalid")
// // returns "", error("short code not found")

// GetStats("xK9mQ2")
// // returns &Stats{
// //     ShortCode:   "xK9mQ2",
// //     LongURL:     "https://www.google.com/search?q=golang+tutorials",
// //     CreatedAt:   2024-01-15 10:30:00,
// //     AccessCount: 2,
// // }
