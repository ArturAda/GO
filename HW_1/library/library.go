package library

type Library struct {
	storage Storage
	hash    Hash
}

func UpdateLibrary(storage Storage, hash Hash) *Library {
	return &Library{storage, hash}
}

func (library *Library) SetStorage(storage Storage) {
	library.storage = storage
}

func (library *Library) SetHash(hash Hash) {
	library.hash = hash
}

func (library *Library) Add(book Book) {
	library.storage.add(&book, library.hash)
}

func (library *Library) AllBooks() []Book {
	return library.storage.all()
}

func (library *Library) Remove(book Book) bool {
	realBook := library.storage.findBook(book, library.hash)
	if (Book{}).Equal(*realBook) {
		return false
	}
	library.storage.remove(realBook.id)
	return true
}

func (library *Library) FindBooks(title string) []Book {
	return library.storage.findAllTitle(title, library.hash)
}
