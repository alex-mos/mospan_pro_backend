package api

import (
	"fmt"
	"strconv"
	"net/http"
	"../books"
)

func Serve() {
	http.HandleFunc("/books", getAllBooks)
	http.HandleFunc("/order", orderBook)

	err := http.ListenAndServe(":8081", nil) // set listen port
	if err != nil {
		panic(err)
	}
	fmt.Println("Server is running on port 8081")
}

// Получить список всех книг в джейсоне
func getAllBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(405)
		fmt.Fprintf(w,  "Method is not allowed")
		return
	}
	allBooks, err := books.GetAll()
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, string(allBooks))
}

// Заказ выбранной книги на указанный телеграм-аккаунт
func orderBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		fmt.Fprintf(w,  "method is not allowed")
		return
	}
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w,  "invalid request")
		return
	}
	id, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w,  "id is not a number")
		return
	}
	telegram := r.Form.Get("telegram")
	err = books.Order(id, telegram)
	if err != nil {
		w.WriteHeader(422)
		fmt.Fprintf(w, err.Error())
		return
	}
	w.WriteHeader(200)
	return
}