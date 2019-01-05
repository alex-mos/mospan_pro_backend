package api

import (
	"fmt"
	"github.com/alex-mos/mospan_pro_backend/books"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func Serve() {
	router := httprouter.New()

	router.GET("/books", getAllBooks)
	router.POST("/order/:id", orderBook)

	fmt.Println("Server is running on port 8081")
	err := http.ListenAndServe(":8081", router)
	if err != nil {
		panic(err)
	}
}

// Получить список всех книг в джейсоне
func getAllBooks(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	allBooks, err := books.GetAll()
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, string(allBooks))
}

// Заказ выбранной книги на указанный телеграм-аккаунт
func orderBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "invalid request")
		return
	}
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "id is not a number")
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
