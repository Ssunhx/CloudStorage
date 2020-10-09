package main

import (
	"cloudstorage/handler"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/file/upload", handler.UploadHandler)
	err := http.ListenAndServe(":9234", nil)
	if err != nil {
		fmt.Println("failed start server")
	}
}
