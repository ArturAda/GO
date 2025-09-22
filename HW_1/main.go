package main

import (
	"fmt"

	"github.com/ArturAda/GO/HW_1/library"
)

func main() {
	book1 := library.Book{
		Title:   "Clean Code",
		Authors: []string{"Robert C. Martin"},
		Edition: 1,
		Year:    2008,
		Pages:   464,
	}

	book2 := library.Book{
		Title:   "The Pragmatic Programmer",
		Authors: []string{"Andrew Hunt", "David Thomas"},
		Edition: 2,
		Year:    2019,
		Pages:   352,
	}

	book3 := library.Book{
		Title:   "Effective Go",
		Authors: []string{"Go Team"},
		Edition: 1,
		Year:    2023,
		Pages:   250,
	}
	book4 := library.Book{
		Title:   "Introduction to Algorithms",
		Authors: []string{"Thomas H. Cormen", "Charles E. Leiserson", "Ronald L. Rivest", "Clifford Stein"},
		Edition: 3,
		Year:    2009,
		Pages:   1312,
	}
	book5 := library.Book{
		Title:   "Design Patterns: Elements of Reusable Object-Oriented Software",
		Authors: []string{"Erich Gamma", "Richard Helm", "Ralph Johnson", "John Vlissides"},
		Edition: 1,
		Year:    1994,
		Pages:   395,
	}
	book6 := library.Book{
		Title:   "Преступление и наказание",
		Authors: []string{"Фёдор Достоевский"},
		Edition: 1,
		Year:    1866,
		Pages:   672,
	}
	fmt.Println(book1)
	lib := library.Library{}
	lib.SetStorage(library.NewStorageSlice())
	lib.SetHash(library.GetHash)
	lib.Add(book1)
	lib.Add(book2)
	lib.Add(book3)
	lib.Add(book4)
	lib.Add(book5)
	fmt.Println(book5)
	fmt.Println(lib.AllBooks())
	fmt.Println(lib.Remove(book6))
	fmt.Println(lib.AllBooks())
	fmt.Println(lib.FindBooks("Design Patterns: Elements of Reusable Object-Oriented Software"))
	fmt.Println(lib.Remove(book5))
	fmt.Println(lib.FindBooks("Design Patterns: Elements of Reusable Object-Oriented Software"))
	fmt.Println(lib.AllBooks())
}
