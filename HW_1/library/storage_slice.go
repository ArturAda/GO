package library

type StorageSlice struct {
	books [][]*Book
	count int
}

func NewStorageSlice() Storage {
	return &StorageSlice{}
}

func (storage *StorageSlice) add(book *Book, hash Hash) idType {
	storage.count++
	if len(storage.books) == 0 {
		book.id = [2]int{0, 0}
		storage.books = append(storage.books, []*Book{})
		storage.books[0] = append(storage.books[0], book)
		return book.id
	}
	if storage.count > len(storage.books) {
		newBooks := make([][]*Book, 2*len(storage.books))
		for i := 0; i < len(storage.books); i++ {
			for j := 0; j < len(storage.books[i]); j++ {
				storage.books[i][j].id[0] = hash(storage.books[i][j])[0] % len(newBooks)
				newBooks[storage.books[i][j].id[0]] = append(newBooks[storage.books[i][j].id[0]], storage.books[i][j])
				storage.books[i][j].id[1] = len(newBooks[storage.books[i][j].id[0]]) - 1
			}
		}
		storage.books = newBooks
	}
	book.id[0] = hash(book)[0] % len(storage.books)
	storage.books[book.id[0]] = append(storage.books[book.id[0]], book)
	book.id[1] = len(storage.books[book.id[0]]) - 1
	return book.id
}

func (storage *StorageSlice) find(id idType) Book {
	if id[0] >= len(storage.books) || id[1] >= len(storage.books[id[0]]) {
		return Book{}
	}
	return *storage.books[id[0]][id[1]]
}

func (storage *StorageSlice) all() []Book {
	library := make([]Book, storage.count)
	position := 0
	for i := 0; i < len(storage.books); i++ {
		for j := 0; j < len(storage.books[i]); j++ {
			library[position] = *storage.books[i][j]
			position++
		}
	}
	return library
}

func (storage *StorageSlice) remove(id idType) bool {
	if id[0] >= len(storage.books) || id[1] >= len(storage.books[id[0]]) {
		return false
	}
	storage.books[id[0]] = append(storage.books[id[0]][:id[1]], storage.books[id[0]][id[1]+1:]...)
	return true
}

func (storage *StorageSlice) findBook(book Book, hash Hash) *Book {
	id := hash(&book)[0] % len(storage.books)
	for i := 0; i < len(storage.books[id]); i++ {
		if book.Equal(*storage.books[id][i]) {
			return storage.books[id][i]
		}
	}
	return &Book{}
}

func (storage *StorageSlice) findAllTitle(title string, hash Hash) []Book {
	book := Book{}
	book.Title = title
	id := hash(&book)[0] % len(storage.books)
	books := make([]Book, 0)
	for i := 0; i < len(storage.books[id]); i++ {
		if storage.books[id][i].Title == title {
			books = append(books, *storage.books[id][i])
		}
	}
	return books
}
