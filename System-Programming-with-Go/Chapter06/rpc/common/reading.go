package common

import (
	"fmt"
	"strings"
)

// List of errors
var (
	ErrISBN      = fmt.Errorf("missing ISBN")
	ErrDuplicate = fmt.Errorf("duplicate book")
	ErrMissing   = fmt.Errorf("missing book")
)

// Course represents a book entry
type Course struct {
	ISBN          string
	Title, Author string
	Year, Pages   int
}

func (b Course) String() string {
	s := strings.Builder{}
	fmt.Fprintf(&s, "%s - %s", b.Title, b.Author)
	if b.Year != 0 {
		fmt.Fprintf(&s, " (%d)", b.Year)
	}
	if b.Pages != 0 {
		fmt.Fprintf(&s, " %d pages", b.Pages)
	}
	return s.String()
}

// ReadingList keeps tracks of books and pages read
type ReadingList struct {
	Courses  []Course
	Progress []int
}

func (r *ReadingList) bookIndex(isbn string) int {
	for i := range r.Courses {
		if isbn == r.Courses[i].ISBN {
			return i
		}
	}
	return -1
}

// AddCourse checks if the book is not present and adds it
func (r *ReadingList) AddCourse(b Course) error {
	if b.ISBN == "" {
		return ErrISBN
	}
	if r.bookIndex(b.ISBN) != -1 {
		return ErrDuplicate
	}
	r.Courses = append(r.Courses, b)
	r.Progress = append(r.Progress, 0)
	return nil
}

// RemoveCourse removes the book from list and forgets its progress
func (r *ReadingList) RemoveCourse(isbn string) error {
	if isbn == "" {
		return ErrISBN
	}
	i := r.bookIndex(isbn)
	if i == -1 {
		return ErrMissing
	}
	// replace the deleted book with the last of the list
	r.Courses[i] = r.Courses[len(r.Courses)-1]
	r.Progress[i] = r.Progress[len(r.Progress)-1]
	// shrink the list of 1 element to remove the duplicate
	r.Courses = r.Courses[:len(r.Courses)-1]
	r.Progress = r.Progress[:len(r.Progress)-1]
	return nil
}

// GetProgress returns the progress of abbok
func (r *ReadingList) GetProgress(isbn string) (int, error) {
	if isbn == "" {
		return -1, ErrISBN
	}
	i := r.bookIndex(isbn)
	if i == -1 {
		return -1, ErrMissing
	}
	return r.Progress[i], nil
}

// SetProgress changes the progress of a book
func (r *ReadingList) SetProgress(isbn string, pages int) error {
	if isbn == "" {
		return ErrISBN
	}
	i := r.bookIndex(isbn)
	if i == -1 {
		return ErrMissing
	}
	if p := r.Courses[i].Pages; pages > p {
		pages = p
	}
	r.Progress[i] = pages
	return nil
}

// AdvanceProgress adds the pages the progress of a book
func (r *ReadingList) AdvanceProgress(isbn string, pages int) error {
	if isbn == "" {
		return ErrISBN
	}
	i := r.bookIndex(isbn)
	if i == -1 {
		return ErrMissing
	}
	if p := r.Courses[i].Pages - r.Progress[i]; p < pages {
		pages = p
	}
	r.Progress[i] += pages
	return nil
}
