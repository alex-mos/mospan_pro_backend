package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var (
	dbPath string = "../database.db"
)

type Book struct {
	id int16
	title string
	author string
	goodreads_link string
	cover_url string
	is_booked bool
	booker_telegram string
}

// Добавить книгу к списку доступных для заказа
func Add(title string, author string, goodreads_link string ) error {
	var newBook Book = Book{title: title, author: author, goodreads_link: goodreads_link}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	statement, err := db.Prepare("INSERT INTO books(title, author, goodreads_link, cover_url, is_booked, booker_telegram) values(?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(newBook.title, newBook.author, newBook.goodreads_link, "", "0", "")
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
		err := rows.Scan(&book.id, &book.title, &book.author, &book.goodreads_link, &book.cover_url, &book.is_booked, &book.booker_telegram)
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
func Order(id int16, booker_telegram string) error  {
	db, err := sql.Open("sqlite3", dbPath)
	statement, err := db.Prepare("UPDATE books SET is_booked=1, booker_telegram=? where id=?")
	if err != nil {
		return err
	}
	_, err = statement.Exec("alex_mos", id)
	if err != nil {
		return err
	}
	db.Close()
	return nil
}

func main() {
	//CREATE TABLE books (id integer primary key autoincrement, title text, author text, goodreads_link text, cover_url text, is_booked integer, booker_telegram text);

	// тест добавления книги
	//err := Add("Neuromancer", "William Gibson", "https://goodreads.com")
	//if err != nil {
	//	panic(err)
	//}

	// тест получения массива книг
	//allBooks, err := GetAll()
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println(allBooks)

	// тест заказа книги
	//Order(4, "alexander_mospan")
}
