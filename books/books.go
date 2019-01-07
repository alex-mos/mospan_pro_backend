package books

/* таблица с книгами
CREATE TABLE `books` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `author` varchar(255) DEFAULT '',
  `title` varchar(255) DEFAULT '',
  `edition` varchar(255) DEFAULT '',
  `goodreads_link` varchar(255) DEFAULT '',
  `cover_url` varchar(255) DEFAULT '',
  `status` varchar(255) DEFAULT '',
  `telegram` varchar(255) DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
*/

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/alex-mos/mospan_pro_backend/email"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strconv"
)

var (
	dbPath string = "root:" + os.Getenv("MYSQL_ROOT_PASSWORD") + "@tcp(database:3306)/mospan_pro?charset=utf8"
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
// Пример: err := Add("William Gibson", "Neuromancer", "мягкая обложка, жёлтые страницы", "https://goodreads.com")
func Add(author string, title string, edition string, goodreads_link string) error {
	var newBook Book = Book{Title: title, Author: author, Goodreads_link: goodreads_link}

	db, err := sql.Open("mysql", dbPath)
	if err != nil {
		return err
	}
	err = db.Ping()
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
// Пример: allBooks, err := GetAll()
func GetAll() ([]byte, error) {
	var result BookSlice
	var resultJSON []byte
	var book Book

	db, err := sql.Open("mysql", dbPath)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("SELECT * FROM books order by status")
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
// Пример: Order(4, "alexander_mospan")
func Order(id int, telegram string) error {
	var book Book

	db, err := sql.Open("mysql", dbPath)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}

	// Получить значения, нужные для заказа
	rows, err := db.Query("SELECT author, title, status FROM books WHERE id=" + strconv.Itoa(id))
	if err != nil {
		return err
	}
	for rows.Next() {
		err := rows.Scan(&book.Author, &book.Title, &book.Status)
		if err != nil {
			return err
		}
	}
	rows.Close()

	// Проверить, не заказана ли уже эта книга
	if book.Status == "" {
		return errors.New("book is not found")
	}
	if book.Status != "free" {
		return errors.New("this book is already reserved")
	}

	// Отправить письмо о заказе книги
	err = email.SendBookRequest(book.Author + " — " + book.Title, telegram)
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
