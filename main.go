package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

/*
func main() {

	array := [][]int{
		{1, 2, 3, 4},
		{2, 3, 4, 1},
		{4, 2, 1, 3},
		{1, 3, 4, 2}}

	fmt.Println(array)
	output, error := Sort(array)
	if error == nil {
		fmt.Println(output)
	} else {
		fmt.Println(error)
	}
}

func Sort(input [][]int) ([][]int, error) {
	size := len(input)
	for i := 0; i < size; i++ {
		for j, value := range input[i] {
			if value == 1 {
				input[i][j], input[i][i] = input[i][i], input[i][j]
			}
		}
	}

	return input, nil
}
*/

type Credential struct {
	ID       int    `json:"id"`
	URL      string `json:"url"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

var db *sql.DB

func main() {

	var err error
	// Строка подключения к базе данных
	db, err = sql.Open("postgres", "host=localhost port=5432 dbname=KeyManager user=postgres password=OlgaK+15 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	initializeDatabase()
	router := gin.Default()

	router.Static("/static", "./static")

	router.GET("/credentials", listCredentials)
	router.POST("/credentials", addCredential)
	router.PUT("/credentials/:id", updateCredential)
	router.DELETE("/credentials/:id", deleteCredential)

	router.Run(":8080")
}

func initializeDatabase() {
	// Проверяем, есть ли уже данные в таблице
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM credentials").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	// Если в таблице нет данных, добавляем их
	if count == 0 {
		_, err := db.Exec(`
            INSERT INTO credentials (url, login, password)
            VALUES
            ('https://example.com', 'user1', 'pass1'),
            ('https://another.com', 'user2', 'pass2')
        `)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Инициализированы начальные данные.")
	}
}

func listCredentials(c *gin.Context) {
	rows, err := db.Query("SELECT id, url, login, password FROM credentials")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	creds := []Credential{}
	for rows.Next() {
		var cred Credential
		if err := rows.Scan(&cred.ID, &cred.URL, &cred.Login, &cred.Password); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		creds = append(creds, cred)
	}

	c.JSON(http.StatusOK, creds)
}
func addCredential(c *gin.Context) {
	var cred Credential
	if err := c.ShouldBindJSON(&cred); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("INSERT INTO credentials (url, login, password) VALUES ($1, $2, $3)", cred.URL, cred.Login, cred.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, cred)
}

func updateCredential(c *gin.Context) {
	id := c.Param("id")
	var cred Credential
	if err := c.ShouldBindJSON(&cred); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("UPDATE credentials SET url = $1, login = $2, password = $3 WHERE id = $4", cred.URL, cred.Login, cred.Password, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cred)
}

func deleteCredential(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("DELETE FROM credentials WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deleted": id})
}
