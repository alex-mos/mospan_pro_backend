package books

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/alex-mos/mospan_pro_backend/email"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
)

var (
	dbPath string = "./database.db"
)

type Book struct {
	Id             int16  `json:"id"`
	Author         string `json:"author"`
	Title          string `json:"title"`
	Edition        string `json:"edition"`
	Goodreads_link string `json:"goodreads_link"`
	Cover_url      string `json:"cover_url"`
	Status         string `json:"status"`
	telegram       string
}

type BookSlice struct {
	Data []Book `json:"data"`
}

// Добавить книгу к списку доступных для заказа
func Add(author string, title string, edition string, goodreads_link string) error {
	var newBook Book = Book{Title: title, Author: author, Goodreads_link: goodreads_link}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	statement, err := db.Prepare("INSERT INTO books(author, title, edition, goodreads_link, cover_url, status, telegram) values(?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(newBook.Author, newBook.Title, newBook.Edition, newBook.Goodreads_link, "", "free", "")
	if err != nil {
		return err
	}
	db.Close()
	return nil
}

// Получить джейсон с полным списком книг
func GetAll() ([]byte, error) {
	var result BookSlice
	var resultJSON []byte
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
		err := rows.Scan(&book.Id, &book.Author, &book.Title, &book.Edition, &book.Goodreads_link, &book.Cover_url, &book.Status, &book.telegram)
		if err != nil {
			return nil, err
		}

		result.Data = append(result.Data, book)
	}
	rows.Close()
	db.Close()

	// Преобразовать массив книг в джейсон
	resultJSON, err = json.Marshal(result)
	if err != nil {
		return nil, err
	}
	return resultJSON, nil
}

// Заказать книгу с выбранным id на указанный телеграм-аккаунт
func Order(id int, telegram string) error {
	var book Book

	db, err := sql.Open("sqlite3", dbPath)

	// Проверить, не заказана ли уже эта книга
	rows, err := db.Query("SELECT status FROM books WHERE id=" + strconv.Itoa(id))
	if err != nil {
		return err
	}
	for rows.Next() {
		err := rows.Scan(&book.Status)
		if err != nil {
			return err
		}
	}
	rows.Close()
	if book.Status == "" {
		return errors.New("book is not found")
	}
	if book.Status != "free" {
		return errors.New("this book is already reserved")
	}

	// Отправить письмо о заказе книги
	title, err := getTitleById(id)
	if err != nil {
		return err
	}
	err = email.SendBookRequest(title, telegram)
	if err != nil {
		return err
	}

	// Пометить в базе, что книга заказана
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

// Получить строку вида "Автор — Название" по id книги. Нужно для формирования письма.
func getTitleById(id int) (string, error) {
	var book Book
	var result string

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return "", err
	}
	rows, err := db.Query("SELECT author, title FROM books WHERE id=" + strconv.Itoa(id))
	if err != nil {
		return "", err
	}
	for rows.Next() {
		err := rows.Scan(&book.Author, &book.Title)
		if err != nil {
			return "", err
		}
	}
	rows.Close()
	result = book.Author + " — " + book.Title
	db.Close()
	return result, nil
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
