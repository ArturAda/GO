package library

import (
	"math/rand"
	"time"
	"unicode/utf8"
)

type idType [2]int

type Book struct {
	id      idType
	Title   string
	Authors []string
	Edition int
	Year    int
	Pages   int
}

func (book *Book) Update(Title string, Authors []string, Edition int, Year int, Pages int) Book {
	book.Title = Title
	book.Authors = Authors
	book.Edition = Edition
	book.Year = Year
	book.Pages = Pages
	return *book
}

func (bookF Book) Equal(bookS Book) bool {
	if bookF.Edition != bookS.Edition {
		return false
	}
	if bookF.Year != bookS.Year {
		return false
	}
	if bookF.Pages != bookS.Pages {
		return false
	}
	if bookF.Title != bookS.Title {
		return false
	}
	if len(bookF.Authors) != len(bookS.Authors) {
		return false
	}
	for i := 0; i < len(bookF.Authors); i++ {
		if bookF.Authors[i] != bookS.Authors[i] {
			return false
		}
	}
	return true
}

type Hash func(book *Book) idType

func GetHash(book *Book) idType {
	var hash [2]int64
	for i := 0; i < len(book.Title); {
		r, size := utf8.DecodeRuneInString(book.Title[i:])
		addSymbol(&hash, int64(r))
		i += size
	}
	/*
		for i := 0; i < len(book.Authors); i++ {
			for j := 0; j < len(book.Authors[i]); {
				r, size := utf8.DecodeRuneInString(book.Authors[i][j:])
				addSymbol(&hash, int64(r))
				j += size
			}
		}
	*/
	return [2]int{int(hash[0]), int(hash[1])}
}

func NewHash(book *Book) idType {
	return [2]int{rand.New(rand.NewSource(time.Now().UnixNano())).Intn(int(x)), 0}
}
