package functions

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Summary Load a recording
// @Description Uploads and saves a recording with the specified name.
// @ID loadRecording
// @Accept multipart/form-data
// @Produce json
// @Tags Recording
// @Param recording formData file true "Recording file (WAV format)"
// @Param name query string true "Name to save the recording as"
// @Success 200 {string} json "{"response": "Saved successfully"}"
// @Failure 400 {string} json "{"response": "Bad Request"}"
// @Failure 415 {string} json "{"response": "Unsupported Media Type. Only WAV files are allowed."}"
// @Failure 500 {string} json "{"response": "Internal Server Error"}"
// @Router /load-recording [POST]
func (h *Handler) LoadRecording(c *gin.Context) {
	var (
		folderPath = "recordings/"
		name       = c.Query("name")
	)

	file, err := c.FormFile("recording")
	if err != nil {
		h.handleResponse(c, 400, "Bad Request")
		return
	}

	if strings.ToLower(filepath.Ext(file.Filename)) != ".wav" {
		h.handleResponse(c, 415, "Unsupported Media Type. Only WAV files are allowed.")
		return
	}

	destinationFilePath := folderPath + name + ".wav"

	destinationFile, err := os.Create(destinationFilePath)
	if err != nil {
		h.handleResponse(c, 500, "Could not create destination file")
		return
	}
	defer destinationFile.Close()

	sourceFile, err := file.Open()
	if err != nil {
		h.handleResponse(c, 500, "Could not open source file")
		return
	}

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		h.handleResponse(c, 500, "Could not copy source file to destination file")
		return
	}

	h.handleResponse(c, 200, "Saved successfully")
}
