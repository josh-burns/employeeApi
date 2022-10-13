package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Employee struct {
	IdNumber int
	FName    string
	SName    string
}

func EmployeeScan() []Employee {

	db, err := sql.Open("mysql", "root:Cypress123!!@tcp(localhost:3306)/classicmodels")

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	var (
		employeeNumber int
		lastName       string
		firstName      string
	)

	rows, err := db.Query("SELECT employeeNumber, lastName, firstName FROM employees;")

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var employeeList []Employee

	for rows.Next() {
		err := rows.Scan(&employeeNumber, &lastName, &firstName)

		if err != nil {
			log.Fatal(err)
		}

		person := new(Employee)
		person.IdNumber = employeeNumber
		person.FName = firstName
		person.SName = lastName

		employeeList = append(employeeList, *person)

	}

	defer db.Close()

	return employeeList
}

func addNewEmployee(request *http.Request) string {
	defer request.Body.Close()
	var newEmployee Employee
	postRequestBody, _ := io.ReadAll(request.Body)
	json.Unmarshal(postRequestBody, &newEmployee)

	db, err := sql.Open("mysql", "root:PASSWORD@tcp(localhost:3306)/classicmodels")

	if err != nil {
		log.Fatal(err)
	}

	queryStringArray := []string{"INSERT INTO `classicmodels`.`employees` (`employeeNumber`, `lastName`, `firstName`, `extension`, `email`, `officeCode`, `reportsTo`, `jobTitle`) VALUES ('", strconv.Itoa(newEmployee.IdNumber), "', '", newEmployee.FName, "', '", newEmployee.SName, "', 'test', 'jost@test.com', '1', '1143', 'Software');"}

	queryString := strings.Join(queryStringArray, "")
	insertResult, err := db.ExecContext(context.Background(), queryString)
	if err != nil {
		log.Fatal(err)
	}

	id, err := insertResult.LastInsertId()
	fmt.Println("success adding ", newEmployee.IdNumber, id)
	return newEmployee.FName
}

func HandleEmployeeRequest(w http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case "GET":
		employees := EmployeeScan()
		var employeesReturnValue []string

		for i := range employees {
			jsonBody := Employee{FName: employees[i].FName, SName: employees[i].SName, IdNumber: employees[i].IdNumber}
			jsonReturn, _ := json.Marshal(jsonBody)
			employeesReturnValue = append(employeesReturnValue, string(jsonReturn))
		}
		fmt.Fprintf(w, strings.Join(employeesReturnValue, ","))
	case "POST":
		successfulPost := addNewEmployee(req)
		fmt.Fprintf(w, successfulPost)
	}
}
