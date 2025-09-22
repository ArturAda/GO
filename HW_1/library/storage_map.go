package library

type StorageMap struct {
	books map[int][]*Book
	count int
}

func NewStorageMap() Storage {
	return &StorageMap{}
}

func (storage *StorageMap) add(book *Book, hash Hash) idType {
	book.id[0] = hash(book)[0]
	storage.books[book.id[0]] = append(storage.books[book.id[0]], book)
	book.id[1] = len(storage.books[book.id[0]]) - 1
	storage.count++
	return book.id
}

func (storage *StorageMap) find(id idType) Book {
	if len(storage.books[id[0]]) == 0 {
		return Book{}
	}
	return *storage.books[id[0]][id[1]]
}

func (storage *StorageMap) all() []Book {
	library := make([]Book, storage.count)
	position := 0
	for key := range storage.books {
		for i := 0; i < len(storage.books[key]); i++ {
			library[position] = *storage.books[key][i]
		}
	}
	return library
}

func (storage *StorageMap) remove(id idType) bool {
	if len(storage.books[id[0]]) == 0 {
		return false
	}
	storage.books[id[0]] = append(storage.books[id[0]][:id[1]], storage.books[id[0]][id[1]+1:]...)
	if len(storage.books[id[0]]) == 0 {
		delete(storage.books, id[0])
	}
	return true
}

func (storage *StorageMap) findBook(book Book, hash Hash) *Book {
	id := hash(&book)[0]
	for i := 0; i < len(storage.books[id]); i++ {
		if book.Equal(*storage.books[id][i]) {
			return storage.books[id][i]
		}
	}
	return &Book{}
}

func (storage *StorageMap) findAllTitle(title string, hash Hash) []Book {
	book := Book{}
	book.Title = title
	id := hash(&book)[0]
	books := make([]Book, 0)
	for i := 0; i < len(storage.books[id]); i++ {
		if storage.books[id][i].Title == title {
			books = append(books, *storage.books[id][i])
		}
	}
	return books
}
