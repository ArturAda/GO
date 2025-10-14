package library

import (
	"hash/fnv"
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

type Hash func(book *Book) idType

func FirstHash(book *Book) idType {
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

func SecondHash(book *Book) idType {
	hashFunc := fnv.New64a()
	hashFunc.Write([]byte(book.Title))
	return [2]int{int(hashFunc.Sum64() % uint64(x)), 0}
}
