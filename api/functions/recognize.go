package functions

import (
	"fmt"
	"math"
	"path/filepath"
	"project/api/models"
	"project/api/util"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Summary Recognize a voice
// @Description Recognizes a voice from an input recording.
// @ID recognizeVoice
// @Accept multipart/form-data
// @Produce json
// @Tags Input
// @Param input_recording formData file true "Recording file (WAV format)"
// @Success 200 {object} models.Response
// @Failure 400 {string} json "{"response": "Bad Request"}"
// @Failure 415 {string} json "{"response": "Unsupported Media Type. Only WAV files are allowed."}"
// @Failure 500 {string} json "{"response": "Internal Server Error"}"
// @Router /recognize-voice [POST]
func (h *Handler) RecognizeVoice(c *gin.Context) {
	inputRecording, err := c.FormFile("input_recording")
	if err != nil {
		h.handleResponse(c, 400, "Bad Request")
		return
	}

	if strings.ToLower(filepath.Ext(inputRecording.Filename)) != ".wav" {
		h.handleResponse(c, 415, "Unsupported Media Type. Only WAV files are allowed.")
		return
	}

	if err := c.SaveUploadedFile(inputRecording, h.cfg.InputFilePath); err != nil {
		h.handleResponse(c, 500, "Could not save uploaded file")
		return
	}

	inputRecordingData, err := util.ReadWavFile(h.cfg.InputFilePath)
	if err != nil {
		h.handleResponse(c, 500, "Could not read uploaded file")
		return
	}

	inputRecordingMFCC := util.CalculateMFCC(inputRecordingData)

	audioFiles, err := filepath.Glob(filepath.Join(h.cfg.RecordingsDir, "*.wav"))
	if err != nil {
		h.handleResponse(c, 500, "Could not read audio files")
	}

	resp := models.Response{}
	dataMap := make(map[string]interface{})

	bestMatch := ""
	var minDistance float64 = math.Inf(1)

	for _, file := range audioFiles {
		fmt.Print("\n\n==================================")
		fmt.Printf("\nAnalyzing %v's voice ...\n", util.ExtractBase(file))
		fmt.Print("==================================")

		refData, err := util.ReadWavFile(file)
		if err != nil {
			h.handleResponse(c, 500, "Could not read reference file")
			continue
		}

		refMFCC := util.CalculateMFCC(refData)
		distance := util.CompareMFCCs(inputRecordingMFCC, refMFCC)

		dataMap[util.ExtractBase(file)] = distance

		if distance < minDistance && distance < 13 {
			bestMatch = util.ExtractBase(file)
			minDistance = distance
		}
	}

	fmt.Print("\n\n")
	resp.Data = dataMap
	resp.Match = bestMatch
	h.handleResponse(c, 200, resp)
}
