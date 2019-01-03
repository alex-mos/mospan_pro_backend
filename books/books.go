package books

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var (
	dbPath string = "./database.db"
)

type Book struct {
	id int16
	author string
	title string
	edition string
	goodreads_link string
	cover_url string
	status string
	telegram string
}

// Добавить книгу к списку доступных для заказа
func Add(author string, title string, edition string, goodreads_link string) error {
	var newBook Book = Book{title: title, author: author, goodreads_link: goodreads_link}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	statement, err := db.Prepare("INSERT INTO books(author, title, edition, goodreads_link, cover_url, status, telegram) values(?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(newBook.author, newBook.title, newBook.edition, newBook.goodreads_link, "", "free", "")
	if err != nil {
		return err
	}
	db.Close()
	return nil
}

// Получить полный список книг
func GetAll() ([]Book, error) {
	var result []Book
	var book Book

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err := rows.Scan(&book.id, &book.author, &book.title, &book.edition, &book.goodreads_link, &book.cover_url, &book.status, &book.telegram)
		if err != nil {
			return nil, err
		}
		result = append(result, book)
	}
	rows.Close()
	db.Close()
	return result, nil
}

// Заказать книгу с выбранным id на указанный телеграм-аккаунт
func Order(id int16, telegram string) error  {
	db, err := sql.Open("sqlite3", dbPath)
	statement, err := db.Prepare("UPDATE books SET status=?, telegram=? where id=?")
	if err != nil {
		return err
	}
	_, err = statement.Exec("reserved", telegram, id)
	if err != nil {
		return err
	}
	db.Close()
	return nil
}


/* таблица с книгами
CREATE TABLE books2 (
	id integer primary key autoincrement,
	author text,
	title text,
	edition text,
	goodreads_link text,
	cover_url text,
	status text,
	telegram text
);
 */

/* пример добавления книги
err := Add("William Gibson", "Neuromancer", "мягкая обложка, жёлтые страницы", "https://goodreads.com")
if err != nil {
	panic(err)
}
*/

/* пример получения массива книг
allBooks, err := GetAll()
if err != nil {
	panic(err)
}
fmt.Println(allBooks)
*/

/* пример заказа книги
Order(4, "alexander_mospan")
*/
