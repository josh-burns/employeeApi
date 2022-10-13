package main

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	http.HandleFunc("/employees", HandleEmployeeRequest)
	http.ListenAndServe(":8090", nil)
}
