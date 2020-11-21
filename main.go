package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"log"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

type Customer struct {
	ID     int		`json:"id"`
	Name  string 	`json:"name"`
	Email string	`json:"email"`
	Status string	`json:"status"`
}

func init() {
	createTable()
}

func createTable() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatal("Connect to Database Error",err)
	}

	defer db.Close()

	createTb := `CREATE TABLE IF NOT EXISTS customers (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT,
		status TEXT
	);
	`
	_, err = db.Exec(createTb)

	if err != nil {
		log.Fatal("Cannot Create Table: ", err)
	}

	fmt.Println("Create Table Success")
}

func createCustomerHandler(c *gin.Context) {
	cust := Customer{}
	
	if err := c.ShouldBindJSON(&cust); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatal("Connect to Database Error",err)
	}

	defer db.Close()

	row := db.QueryRow("INSERT INTO customers (name, email, status) values ($1, $2, $3) RETURNING id", cust.Name, cust.Email ,cust.Status)

	var id int

	err = row.Scan(&id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}

	fmt.Println("Insert Customer Success id: ",id)

	cust.ID = id

	c.JSON(http.StatusCreated, cust)
}

func getCustomerHandler(c *gin.Context) {
	custs := []Customer{}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatal("Connect to Database Error",err)
	}


	statement, err := db.Prepare("SELECT id,name,email,status FROM customers")
	if err != nil {
		log.Fatal("Cannot Prepare Statement Error",err)
	}

	rows, err := statement.Query()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}

	for rows.Next() {
		var id int
		var name, email, status string

		err := rows.Scan(&id, &name, &email, &status)
		if err != nil {
			log.Fatal("can't Scan row into variable: %s", err)
		}
		cust := Customer{id, name, email, status}
		custs = append(custs, cust)
	}
	fmt.Println("Query Customers success")

	c.JSON(http.StatusOK, custs)
}

func getCustomerByIDHandler(c *gin.Context) {

	searchID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatal("Connect to Database Error",err)
	}

	statement, err := db.Prepare("SELECT id, name, email, status FROM customers where id=$1")
	if err != nil {
		log.Fatal("Cannot Prepare Statement", err)
	}

	row := statement.QueryRow(searchID)
	var id int
	var name, email, status string

	err = row.Scan(&id, &name, &email, &status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}

	cust := Customer{id, name, email, status}
	fmt.Println("Query Customers success")

	c.JSON(http.StatusOK, cust)
}

func updateCustomerHandler(c *gin.Context) {
	cust := Customer{}
	searchID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if err := c.ShouldBindJSON(&cust); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatal("Connect to Database Error",err)
	}

	statement, err := db.Prepare("UPDATE customers SET name=$2, email= $3, status=$4 where id=$1")

	if err != nil {
		log.Fatal("Cannot Prepare Statement", err)
	}

	cust.ID = searchID

	_, err = statement.Exec(cust.ID, cust.Name, cust.Email, cust.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}

	fmt.Println("Update Customers success")

	c.JSON(http.StatusOK, cust)

}

func deleteCustomerHandler(c *gin.Context) {
	searchID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatal("Connect to Database Error",err)
	}

	statement, err := db.Prepare("DELETE FROM customers WHERE id=$1")

	if err != nil {
		log.Fatal("Cannot Prepare Statement", err)
	}

	_, err = statement.Exec(searchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
	}
	fmt.Println("Delete success")

	c.JSON(http.StatusOK, gin.H{"message":"customer deleted"})
}

func main() {
	fmt.Println("customer service")
	r := gin.Default()
	r.POST("/customers", createCustomerHandler)
	r.GET("/customers", getCustomerHandler)
	r.GET("/customers/:id", getCustomerByIDHandler)
	r.PUT("/customers/:id", updateCustomerHandler)
	r.DELETE("/customers/:id", deleteCustomerHandler)
	r.Run(":2009")
}
