package kindle

import (
	"sync"
	"testing"
	"time"
)

// Test Books
var (
	book1 = Book{
		ID:         "book-1",
		Title:      "The Go Programming Language",
		Author:     "Alan Donovan",
		TotalPages: 380,
	}
	book2 = Book{
		ID:         "book-2",
		Title:      "Clean Code",
		Author:     "Robert C. Martin",
		TotalPages: 464,
	}
	book3 = Book{
		ID:         "book-3",
		Title:      "Design Patterns",
		Author:     "Gang of Four",
		TotalPages: 395,
	}
)

// ==================== Basic Tests ====================

func TestAddBook(t *testing.T) {
	lib := NewLibrary()

	err := lib.AddBook("user-1", book1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	books, _ := lib.GetUserBooks("user-1")
	if len(books) != 1 {
		t.Fatalf("expected 1 book, got %d", len(books))
	}
	if books[0].ID != book1.ID {
		t.Fatalf("expected book ID %s, got %s", book1.ID, books[0].ID)
	}
}

func TestAddBookDuplicate(t *testing.T) {
	lib := NewLibrary()

	lib.AddBook("user-1", book1)
	err := lib.AddBook("user-1", book1)

	if err != ErrBookAlreadyExists {
		t.Fatalf("expected ErrBookAlreadyExists, got %v", err)
	}
}

func TestAddSameBookDifferentUsers(t *testing.T) {
	lib := NewLibrary()

	err1 := lib.AddBook("user-1", book1)
	err2 := lib.AddBook("user-2", book1)

	if err1 != nil || err2 != nil {
		t.Fatalf("expected no errors, got %v and %v", err1, err2)
	}

	books1, _ := lib.GetUserBooks("user-1")
	books2, _ := lib.GetUserBooks("user-2")

	if len(books1) != 1 || len(books2) != 1 {
		t.Fatalf("each user should have 1 book")
	}
}

func TestRemoveBook(t *testing.T) {
	lib := NewLibrary()

	lib.AddBook("user-1", book1)
	lib.AddBook("user-1", book2)

	err := lib.RemoveBook("user-1", book1.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	books, _ := lib.GetUserBooks("user-1")
	if len(books) != 1 {
		t.Fatalf("expected 1 book remaining, got %d", len(books))
	}
}

func TestRemoveBookNotFound(t *testing.T) {
	lib := NewLibrary()

	lib.AddBook("user-1", book1)

	err := lib.RemoveBook("user-1", "non-existent")
	if err != ErrBookNotInLibrary {
		t.Fatalf("expected ErrBookNotInLibrary, got %v", err)
	}
}

func TestRemoveBookUserNotFound(t *testing.T) {
	lib := NewLibrary()

	err := lib.RemoveBook("non-existent", book1.ID)
	if err != ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestGetUserBooksEmpty(t *testing.T) {
	lib := NewLibrary()

	books, err := lib.GetUserBooks("user-1")
	if err != nil {
		t.Fatalf("expected no error for non-existent user, got %v", err)
	}
	if len(books) != 0 {
		t.Fatalf("expected empty slice, got %d books", len(books))
	}
}

// ==================== Reading Session Tests ====================

func TestOpenBook(t *testing.T) {
	lib := NewLibrary()

	lib.AddBook("user-1", book1)

	session, err := lib.OpenBook("user-1", book1.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if session == nil {
		t.Fatal("expected session, got nil")
	}
	if session.Book.ID != book1.ID {
		t.Fatalf("expected book ID %s, got %s", book1.ID, session.Book.ID)
	}
	if session.CurrentPage != 0 {
		t.Fatalf("expected page 0 for new book, got %d", session.CurrentPage)
	}
}

func TestOpenBookNotInLibrary(t *testing.T) {
	lib := NewLibrary()

	_, err := lib.OpenBook("user-1", book1.ID)
	if err != ErrBookNotInLibrary {
		t.Fatalf("expected ErrBookNotInLibrary, got %v", err)
	}
}

func TestOpenAnotherBookWhileActive(t *testing.T) {
	lib := NewLibrary()

	lib.AddBook("user-1", book1)
	lib.AddBook("user-1", book2)

	lib.OpenBook("user-1", book1.ID)
	_, err := lib.OpenBook("user-1", book2.ID)

	if err != ErrAnotherBookActive {
		t.Fatalf("expected ErrAnotherBookActive, got %v", err)
	}
}

func TestOpenSameBookTwice(t *testing.T) {
	lib := NewLibrary()

	lib.AddBook("user-1", book1)

	session1, _ := lib.OpenBook("user-1", book1.ID)
	lib.UpdateProgress("user-1", 50)

	session2, err := lib.OpenBook("user-1", book1.ID)
	if err != nil {
		t.Fatalf("expected no error when opening same book, got %v", err)
	}
	if session2.CurrentPage != 50 {
		t.Fatalf("expected page 50, got %d", session2.CurrentPage)
	}
	if session1.Book.ID != session2.Book.ID {
		t.Fatal("sessions should reference same book")
	}
}

func TestResumeReading(t *testing.T) {
	lib := NewLibrary()

	lib.AddBook("user-1", book1)

	// First session
	lib.OpenBook("user-1", book1.ID)
	lib.UpdateProgress("user-1", 100)
	lib.CloseBook("user-1")

	// Second session - should resume
	session, _ := lib.OpenBook("user-1", book1.ID)
	if session.CurrentPage != 100 {
		t.Fatalf("expected to resume at page 100, got %d", session.CurrentPage)
	}
}

// ==================== Progress Tests ====================

func TestUpdateProgress(t *testing.T) {
	lib := NewLibrary()

	lib.AddBook("user-1", book1)
	lib.OpenBook("user-1", book1.ID)

	err := lib.UpdateProgress("user-1", 50)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	session, _ := lib.GetActiveBook("user-1")
	if session.CurrentPage != 50 {
		t.Fatalf("expected page 50, got %d", session.CurrentPage)
	}
}

func TestUpdateProgressNoActiveBook(t *testing.T) {
	lib := NewLibrary()

	lib.AddBook("user-1", book1)

	err := lib.UpdateProgress("user-1", 50)
	if err != ErrNoActiveBook {
		t.Fatalf("expected ErrNoActiveBook, got %v", err)
	}
}

func TestUpdateProgressInvalidPage(t *testing.T) {
	lib := NewLibrary()

	lib.AddBook("user-1", book1)
	lib.OpenBook("user-1", book1.ID)

	// Negative page
	err := lib.UpdateProgress("user-1", -1)
	if err != ErrInvalidPage {
		t.Fatalf("expected ErrInvalidPage for negative page, got %v", err)
	}

	// Page beyond total
	err = lib.UpdateProgress("user-1", book1.TotalPages+1)
	if err != ErrInvalidPage {
		t.Fatalf("expected ErrInvalidPage for page beyond total, got %v", err)
	}
}

func TestGetReadingProgress(t *testing.T) {
	lib := NewLibrary()

	lib.AddBook("user-1", book1)
	lib.OpenBook("user-1", book1.ID)
	lib.UpdateProgress("user-1", 190) // 50% of 380 pages
	lib.CloseBook("user-1")

	progress, err := lib.GetReadingProgress("user-1", book1.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if progress.CurrentPage != 190 {
		t.Fatalf("expected page 190, got %d", progress.CurrentPage)
	}
	if progress.Percentage != 50.0 {
		t.Fatalf("expected 50%% progress, got %.2f%%", progress.Percentage)
	}
}

func TestGetReadingProgressNeverOpened(t *testing.T) {
	lib := NewLibrary()

	lib.AddBook("user-1", book1)

	progress, err := lib.GetReadingProgress("user-1", book1.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if progress != nil {
		t.Fatalf("expected nil progress for unopened book, got %+v", progress)
	}
}

// ==================== Close Book Tests ====================

func TestCloseBook(t *testing.T) {
	lib := NewLibrary()

	lib.AddBook("user-1", book1)
	lib.OpenBook("user-1", book1.ID)
	lib.UpdateProgress("user-1", 75)

	err := lib.CloseBook("user-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	session, _ := lib.GetActiveBook("user-1")
	if session != nil {
		t.Fatal("expected no active session after close")
	}

	// Verify progress was saved
	progress, _ := lib.GetReadingProgress("user-1", book1.ID)
	if progress.CurrentPage != 75 {
		t.Fatalf("expected progress saved at page 75, got %d", progress.CurrentPage)
	}
}

func TestCloseBookNoActive(t *testing.T) {
	lib := NewLibrary()

	err := lib.CloseBook("user-1")
	if err != ErrNoActiveBook {
		t.Fatalf("expected ErrNoActiveBook, got %v", err)
	}
}

func TestRemoveActiveBook(t *testing.T) {
	lib := NewLibrary()

	lib.AddBook("user-1", book1)
	lib.OpenBook("user-1", book1.ID)
	lib.UpdateProgress("user-1", 50)

	err := lib.RemoveBook("user-1", book1.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Should have no active book
	session, _ := lib.GetActiveBook("user-1")
	if session != nil {
		t.Fatal("expected no active session after removing active book")
	}
}

// ==================== Get Active Book Tests ====================

func TestGetActiveBook(t *testing.T) {
	lib := NewLibrary()

	lib.AddBook("user-1", book1)
	lib.OpenBook("user-1", book1.ID)

	session, err := lib.GetActiveBook("user-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if session == nil {
		t.Fatal("expected active session")
	}
	if session.Book.ID != book1.ID {
		t.Fatalf("expected book %s, got %s", book1.ID, session.Book.ID)
	}
}

func TestGetActiveBookNone(t *testing.T) {
	lib := NewLibrary()

	session, err := lib.GetActiveBook("user-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if session != nil {
		t.Fatal("expected nil session when no active book")
	}
}

// ==================== Isolation Tests ====================

func TestUserIsolation(t *testing.T) {
	lib := NewLibrary()

	// User 1 has book1
	lib.AddBook("user-1", book1)
	lib.OpenBook("user-1", book1.ID)
	lib.UpdateProgress("user-1", 100)
	lib.CloseBook("user-1")

	// User 2 has same book
	lib.AddBook("user-2", book1)
	lib.OpenBook("user-2", book1.ID)
	lib.UpdateProgress("user-2", 200)
	lib.CloseBook("user-2")

	// Check isolation
	prog1, _ := lib.GetReadingProgress("user-1", book1.ID)
	prog2, _ := lib.GetReadingProgress("user-2", book1.ID)

	if prog1.CurrentPage != 100 {
		t.Fatalf("user-1 progress should be 100, got %d", prog1.CurrentPage)
	}
	if prog2.CurrentPage != 200 {
		t.Fatalf("user-2 progress should be 200, got %d", prog2.CurrentPage)
	}
}

// ==================== Concurrency Tests ====================

func TestConcurrentAddBooks(t *testing.T) {
	lib := NewLibrary()
	var wg sync.WaitGroup

	// Multiple goroutines adding books to the same user
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			book := Book{
				ID:         string(rune('a'+idx%26)) + string(rune('0'+idx/26)),
				Title:      "Book",
				Author:     "Author",
				TotalPages: 100,
			}
			lib.AddBook("user-1", book)
		}(i)
	}

	wg.Wait()

	books, _ := lib.GetUserBooks("user-1")
	if len(books) == 0 {
		t.Fatal("expected books to be added")
	}
}

func TestConcurrentOpenAndUpdate(t *testing.T) {
	lib := NewLibrary()
	lib.AddBook("user-1", book1)

	var wg sync.WaitGroup

	// Open book
	lib.OpenBook("user-1", book1.ID)

	// Multiple goroutines updating progress
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			lib.UpdateProgress("user-1", page)
		}(i)
	}

	wg.Wait()

	session, _ := lib.GetActiveBook("user-1")
	if session == nil {
		t.Fatal("expected active session")
	}
	// Page should be between 0 and 49
	if session.CurrentPage < 0 || session.CurrentPage > 49 {
		t.Fatalf("unexpected page %d", session.CurrentPage)
	}
}

func TestConcurrentMultipleUsers(t *testing.T) {
	lib := NewLibrary()
	var wg sync.WaitGroup

	// 10 users, each doing operations concurrently
	for u := 0; u < 10; u++ {
		wg.Add(1)
		go func(userID string) {
			defer wg.Done()

			lib.AddBook(userID, book1)
			lib.AddBook(userID, book2)

			lib.OpenBook(userID, book1.ID)
			lib.UpdateProgress(userID, 50)
			lib.CloseBook(userID)

			lib.OpenBook(userID, book2.ID)
			lib.UpdateProgress(userID, 100)
			lib.CloseBook(userID)

		}(string(rune('A' + u)))
	}

	wg.Wait()

	// Verify each user has correct state
	for u := 0; u < 10; u++ {
		userID := string(rune('A' + u))
		books, _ := lib.GetUserBooks(userID)
		if len(books) != 2 {
			t.Fatalf("user %s should have 2 books, got %d", userID, len(books))
		}
	}
}

// ==================== Timestamp Tests ====================

func TestStartedAtTimestamp(t *testing.T) {
	lib := NewLibrary()
	lib.AddBook("user-1", book1)

	beforeOpen := time.Now()
	lib.OpenBook("user-1", book1.ID)
	afterOpen := time.Now()

	lib.CloseBook("user-1")

	progress, _ := lib.GetReadingProgress("user-1", book1.ID)
	if progress.StartedAt.Before(beforeOpen) || progress.StartedAt.After(afterOpen) {
		t.Fatal("StartedAt should be set when first opening book")
	}

	// Open again - StartedAt should not change
	originalStartedAt := progress.StartedAt
	time.Sleep(10 * time.Millisecond)

	lib.OpenBook("user-1", book1.ID)
	lib.CloseBook("user-1")

	progress, _ = lib.GetReadingProgress("user-1", book1.ID)
	if !progress.StartedAt.Equal(originalStartedAt) {
		t.Fatal("StartedAt should not change on subsequent opens")
	}
}

func TestLastReadAtTimestamp(t *testing.T) {
	lib := NewLibrary()
	lib.AddBook("user-1", book1)

	lib.OpenBook("user-1", book1.ID)
	beforeUpdate := time.Now()
	lib.UpdateProgress("user-1", 50)
	afterUpdate := time.Now()
	lib.CloseBook("user-1")

	progress, _ := lib.GetReadingProgress("user-1", book1.ID)
	if progress.LastReadAt.Before(beforeUpdate) || progress.LastReadAt.After(afterUpdate) {
		t.Fatal("LastReadAt should be updated on progress update")
	}
}

// ==================== Edge Cases ====================

func TestUpdateProgressToLastPage(t *testing.T) {
	lib := NewLibrary()
	lib.AddBook("user-1", book1)
	lib.OpenBook("user-1", book1.ID)

	// Should be able to set to exactly total pages (finished book)
	err := lib.UpdateProgress("user-1", book1.TotalPages)
	if err != nil {
		t.Fatalf("should be able to set progress to total pages, got %v", err)
	}

	progress, _ := lib.GetReadingProgress("user-1", book1.ID)
	if progress.Percentage != 100.0 {
		t.Fatalf("expected 100%% progress, got %.2f%%", progress.Percentage)
	}
}

func TestMultipleBooksProgress(t *testing.T) {
	lib := NewLibrary()

	lib.AddBook("user-1", book1)
	lib.AddBook("user-1", book2)
	lib.AddBook("user-1", book3)

	// Read book1
	lib.OpenBook("user-1", book1.ID)
	lib.UpdateProgress("user-1", 100)
	lib.CloseBook("user-1")

	// Read book2
	lib.OpenBook("user-1", book2.ID)
	lib.UpdateProgress("user-1", 200)
	lib.CloseBook("user-1")

	// book3 never opened

	prog1, _ := lib.GetReadingProgress("user-1", book1.ID)
	prog2, _ := lib.GetReadingProgress("user-1", book2.ID)
	prog3, _ := lib.GetReadingProgress("user-1", book3.ID)

	if prog1.CurrentPage != 100 {
		t.Fatalf("book1 progress should be 100, got %d", prog1.CurrentPage)
	}
	if prog2.CurrentPage != 200 {
		t.Fatalf("book2 progress should be 200, got %d", prog2.CurrentPage)
	}
	if prog3 != nil {
		t.Fatal("book3 should have nil progress (never opened)")
	}
}
