package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := getDBConfig()

	db, err := OpenDB(cfg)
	if err != nil {
		log.Fatalf("koneksi DB gagal: %v", err)
	}
	defer db.Close()

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		indexHandler(c, db)
	})

	r.POST("/add", func(c *gin.Context) {
		addHandler(c, db)
	})

	r.POST("/delete", func(c *gin.Context) {
		deleteHandler(c, db)
	})

	r.GET("/edit/:id", func(c *gin.Context) {
		editPageHandler(c, db)
	})

	r.POST("/edit", func(c *gin.Context) {
		updateHandler(c, db)
	})

	log.Println("Server berjalan di http://localhost:8080")
	r.Run(":8080")
}
