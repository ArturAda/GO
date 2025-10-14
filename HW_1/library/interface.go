package library

type Storage interface {
	add(b *Book, hash Hash) idType /* возвращаем id книги после добавления в библиотеку */
	find(id idType) Book
	all() []Book
	remove(id idType) bool /* вернем флаг если нашлась такая книга */
	findBook(book Book, hash Hash) *Book
	findAllTitle(title string, hash Hash) []Book
}
