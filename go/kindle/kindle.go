package kindle

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrBookNotFound      = errors.New("book not found")
	ErrBookAlreadyExists = errors.New("book already exists in library")
	ErrUserNotFound      = errors.New("user not found")
	ErrNoActiveBook      = errors.New("no active book")
	ErrBookNotInLibrary  = errors.New("book not in user's library")
	ErrInvalidPage       = errors.New("invalid page number")
	ErrAnotherBookActive = errors.New("another book is currently active")
)

type Library interface {
	// Book Management
	AddBook(userID string, book Book) error
	RemoveBook(userID string, bookID string) error
	GetUserBooks(userID string) ([]Book, error)

	// Reading Session
	OpenBook(userID string, bookID string) (*ReadingSession, error)
	UpdateProgress(userID string, currentPage int) error
	CloseBook(userID string) error

	// Progress & Stats
	GetReadingProgress(userID string, bookID string) (*Progress, error)
	GetActiveBook(userID string) (*ReadingSession, error)
}

// Book represents a book in the system
type Book struct {
	ID         string
	Title      string
	Author     string
	TotalPages int
}

// Progress represents reading progress for a specific book
type Progress struct {
	BookID      string
	CurrentPage int
	TotalPages  int
	Percentage  float64
	LastReadAt  time.Time
	StartedAt   time.Time
}

// ReadingSession represents an active reading session
type ReadingSession struct {
	Book        Book
	CurrentPage int
	StartedAt   time.Time
	OpenedAt    time.Time
}

// UserLibrary stores a user's books and reading data
type UserLibrary struct {
	Books         map[string]Book      // bookID -> Book
	Progress      map[string]*Progress // bookID -> Progress
	ActiveBookID  *string              // currently active book (nil if none)
	ActiveSession *ReadingSession
}

// KindleLibrary implements the Library interface
type KindleLibrary struct {
	mu    sync.RWMutex
	users map[string]*UserLibrary
}

// NewLibrary creates a new KindleLibrary instance
func NewLibrary() Library {
	return &KindleLibrary{
		users: make(map[string]*UserLibrary),
	}
}

// getOrCreateUser is a helper to get or create a user's library
// Hint: This helper might be useful
func (k *KindleLibrary) getOrCreateUser(userID string) *UserLibrary {
	user, ok := k.users[userID]
	if ok {
		return user
	}

	newLibrary := &UserLibrary{
		Books:         make(map[string]Book),
		Progress:      make(map[string]*Progress),
		ActiveBookID:  nil,
		ActiveSession: nil,
	}
	k.users[userID] = newLibrary
	return newLibrary
}

// AddBook adds a book to a user's library
func (k *KindleLibrary) AddBook(userID string, book Book) error {
	// TODO: Implement this
	//
	// Requirements:
	// 1. Create user library if doesn't exist
	// 2. Return ErrBookAlreadyExists if book already in library
	// 3. Store the book in user's library
	//
	// Hint: Use mutex for thread safety
	k.mu.Lock()
	defer k.mu.Unlock()

	userLib := k.getOrCreateUser(userID)

	if _, exists := userLib.Books[book.ID]; exists {
		return ErrBookAlreadyExists
	}

	userLib.Books[book.ID] = book
	return nil
}

// RemoveBook removes a book from a user's library
func (k *KindleLibrary) RemoveBook(userID string, bookID string) error {
	// TODO: Implement this
	//
	// Requirements:
	// 1. Return ErrUserNotFound if user doesn't exist
	// 2. Return ErrBookNotInLibrary if book not in library
	// 3. If this book is currently active, close it first
	// 4. Remove book and its progress
	k.mu.Lock()
	defer k.mu.Unlock()

	userLib, exists := k.users[userID]
	if !exists {
		return ErrUserNotFound
	}

	if _, exists = userLib.Books[bookID]; !exists {
		return ErrBookNotInLibrary
	}

	if userLib.ActiveBookID != nil && *userLib.ActiveBookID == bookID {
		userLib.ActiveBookID = nil
		userLib.ActiveSession = nil
	}

	delete(userLib.Books, bookID)
	delete(userLib.Progress, bookID)

	return nil
}

// GetUserBooks returns all books in a user's library
func (k *KindleLibrary) GetUserBooks(userID string) ([]Book, error) {
	// TODO: Implement this
	//
	// Requirements:
	// 1. Return empty slice if user doesn't exist (not an error)
	// 2. Return all books in user's library
	k.mu.RLock()
	defer k.mu.RUnlock()

	userLib, exists := k.users[userID]
	if !exists {
		return []Book{}, nil
	}

	books := make([]Book, 0, len(userLib.Books))

	for _, book := range userLib.Books {
		books = append(books, book)
	}

	return books, nil
}

// OpenBook opens a book for reading, resuming from last position
func (k *KindleLibrary) OpenBook(userID string, bookID string) (*ReadingSession, error) {
	// TODO: Implement this
	//
	// Requirements:
	// 1. Return ErrBookNotInLibrary if book not in user's library
	// 2. Return ErrAnotherBookActive if user has another book open
	// 3. If same book is already active, return current session
	// 4. Resume from last read page (or page 0 if never read)
	// 5. Create/update reading session
	// 6. Initialize Progress.StartedAt on first open only
	//
	// Hint: Check if progress exists to determine if it's first open
	k.mu.Lock()
	defer k.mu.Unlock()

	userLib, exists := k.users[userID]
	if !exists {
		return nil, ErrBookNotInLibrary
	}

	book, exists := userLib.Books[bookID]
	if !exists {
		return nil, ErrBookNotInLibrary
	}

	if userLib.ActiveBookID != nil && *userLib.ActiveBookID != bookID {
		return nil, ErrAnotherBookActive
	}

	if userLib.ActiveBookID != nil && *userLib.ActiveBookID == bookID {
		return userLib.ActiveSession, nil
	}

	currentPage := 0
	var startedAt time.Time

	if progress, hasProgress := userLib.Progress[bookID]; hasProgress {
		currentPage = progress.CurrentPage
		startedAt = progress.StartedAt
	} else {
		startedAt = now()
	}

	newSession := &ReadingSession{
		Book:        book,
		CurrentPage: currentPage,
		StartedAt:   startedAt,
		OpenedAt:    now(),
	}

	if _, hasProgress := userLib.Progress[bookID]; !hasProgress {
		userLib.Progress[bookID] = &Progress{
			BookID:      bookID,
			CurrentPage: currentPage,
			TotalPages:  book.TotalPages,
			Percentage:  calculatePercentage(currentPage, book.TotalPages),
			StartedAt:   startedAt,
			LastReadAt:  now(),
		}
	}

	userLib.ActiveBookID = &bookID
	userLib.ActiveSession = newSession

	return newSession, nil
}

// UpdateProgress updates the current reading position
func (k *KindleLibrary) UpdateProgress(userID string, currentPage int) error {
	// TODO: Implement this
	//
	// Requirements:
	// 1. Return ErrNoActiveBook if no book is open
	// 2. Return ErrInvalidPage if page < 0 or page > totalPages
	// 3. Update current page in session and progress
	// 4. Update LastReadAt timestamp
	// 5. Calculate percentage: (currentPage / totalPages) * 100
	k.mu.Lock()
	defer k.mu.Unlock()

	userLib, exists := k.users[userID]
	if !exists {
		return ErrNoActiveBook
	}

	if userLib.ActiveBookID == nil || userLib.ActiveSession == nil {
		return ErrNoActiveBook
	}

	book := userLib.ActiveSession.Book

	if currentPage > book.TotalPages || currentPage < 0 {
		return ErrInvalidPage
	}

	userLib.ActiveSession.CurrentPage = currentPage

	bookID := *userLib.ActiveBookID
	userLib.Progress[bookID].CurrentPage = currentPage
	userLib.Progress[bookID].LastReadAt = now()
	userLib.Progress[bookID].Percentage = calculatePercentage(currentPage, book.TotalPages)

	return nil
}

// CloseBook closes the currently active book
func (k *KindleLibrary) CloseBook(userID string) error {
	// TODO: Implement this
	//
	// Requirements:
	// 1. Return ErrNoActiveBook if no book is open
	// 2. Clear active session (progress is already saved via UpdateProgress)
	k.mu.Lock()
	defer k.mu.Unlock()

	userLib, exists := k.users[userID]
	if !exists {
		return ErrNoActiveBook
	}

	if userLib.ActiveBookID == nil || userLib.ActiveSession == nil {
		return ErrNoActiveBook
	}

	userLib.ActiveBookID = nil
	userLib.ActiveSession = nil

	return nil
}

// GetReadingProgress returns progress for a specific book
func (k *KindleLibrary) GetReadingProgress(userID string, bookID string) (*Progress, error) {
	// TODO: Implement this
	//
	// Requirements:
	// 1. Return ErrUserNotFound if user doesn't exist
	// 2. Return ErrBookNotInLibrary if book not in library
	// 3. Return nil progress if book was never opened (not an error)
	k.mu.Lock()
	defer k.mu.Unlock()

	userLib, exists := k.users[userID]
	if !exists {
		return nil, ErrUserNotFound
	}

	if _, exists := userLib.Books[bookID]; !exists {
		return nil, ErrBookNotInLibrary
	}

	progress, exists := userLib.Progress[bookID]
	if !exists {
		return nil, nil
	}

	return progress, nil
}

// GetActiveBook returns the currently active reading session
func (k *KindleLibrary) GetActiveBook(userID string) (*ReadingSession, error) {
	// TODO: Implement this
	//
	// Requirements:
	// 1. Return nil, nil if no active book (not an error)
	// 2. Return current session if book is open
	k.mu.Lock()
	defer k.mu.Unlock()

	userLib, exists := k.users[userID]
	if !exists {
		return nil, nil
	}

	if userLib.ActiveBookID == nil || userLib.ActiveSession == nil {
		return nil, nil
	}

	return userLib.ActiveSession, nil
}

// Helper function to get current time (useful for testing)
var now = time.Now

// Helper function to calculate percentage
func calculatePercentage(currentPage, totalPages int) float64 {
	if totalPages == 0 {
		return 0
	}
	return float64(currentPage) / float64(totalPages) * 100
}
