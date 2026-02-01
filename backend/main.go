package main

import (
	"fmt"
	"net/http"
	"os"

	"music-app/pkg/analyzer"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	r.MaxMultipartMemory = 8 << 20 // 8 MiB

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/analyze", func(c *gin.Context) {
		// Single file
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
			return
		}

		// Upload the file to specific dst.
		// For now, we'll just acknowledge receipt
		filename := file.Filename
		fmt.Printf("Received file: %s\n", filename)

		// Create a temporary file to process
		tmpFile, err := os.CreateTemp("", "upload-*.tmp")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temp file"})
			return
		}
		defer os.Remove(tmpFile.Name()) // clean up
		defer tmpFile.Close()

		if err := c.SaveUploadedFile(file, tmpFile.Name()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		// Analyze the file
		result, err := analyzer.Analyze(tmpFile.Name())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Analysis failed: %v", err)})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   "success",
			"filename": filename,
			"bpm":      result.BPM,
			"key":      result.Key,
		})
	})

	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
