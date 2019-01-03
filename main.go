package main

import (
	"./api"
	//"./books"
	//"fmt"
)


func main() {
	//err := books.Order(3, "wintermute1")
	//if err != nil {
	//	panic(err)
	//}

	//allBooks, err := books.GetAll()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(allBooks)

	api.Serve()
}
