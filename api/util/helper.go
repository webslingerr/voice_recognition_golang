package util

import (
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-audio/wav"
	"github.com/mjibson/go-dsp/fft"
)

const (
	sampleRate      = 16000
	frameSize       = 400
	frameStep       = 160
	numCoefficients = 13
)

func ReadWavFile(filePath string) ([]float64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := wav.NewDecoder(file)
	buf, err := decoder.FullPCMBuffer()
	if err != nil {
		return nil, err
	}

	data := make([]float64, len(buf.Data))
	for i, val := range buf.Data {
		data[i] = float64(val)
	}

	return data, nil
}

func CalculateMFCC(data []float64) []float64 {
	fftResult := fft.FFTReal(data)

	magnitudes := make([]float64, len(fftResult)/2)
	for i := range magnitudes {
		realPart := real(fftResult[i])
		imagPart := imag(fftResult[i])
		magnitudes[i] = math.Sqrt(realPart*realPart + imagPart*imagPart)
	}

	melFilterbank := CreateMelFilterbank(len(magnitudes), numCoefficients)
	melFiltered := applyFilterbank(magnitudes, melFilterbank)

	logMelFiltered := applyLog(melFiltered)

	dctResult := applyDCT(logMelFiltered)

	return dctResult[:numCoefficients]
}

func CompareMFCCs(mfcc1, mfcc2 []float64) float64 {
	distance := 0.0
	for i := 0; i < len(mfcc1); i++ {
		distance += (mfcc1[i] - mfcc2[i]) * (mfcc1[i] - mfcc2[i])
	}
	distance = math.Sqrt(distance)

	// You may need to adjust the threshold based on your specific use case
	//threshold := 10.0

	return distance
}

func CreateMelFilterbank(size, numFilters int) [][]float64 {
	filterbank := make([][]float64, numFilters)
	for i := range filterbank {
		filterbank[i] = make([]float64, size)
	}

	mel := func(frequency float64) float64 {
		return 2595 * math.Log10(1+frequency/700)
	}

	invMel := func(mel float64) float64 {
		return 700 * (math.Pow(10, mel/2595) - 1)
	}

	melStart := mel(300)
	melEnd := mel(sampleRate / 2)
	melStep := (melEnd - melStart) / float64(numFilters+1)

	for i := 0; i < numFilters; i++ {
		centerMel := melStart + float64(i+1)*melStep

		centerFreq := invMel(centerMel)
		leftFreq := invMel(centerMel - melStep)
		rightFreq := invMel(centerMel + melStep)

		for j := 0; j < size; j++ {
			frequency := float64(j) * float64(sampleRate) / float64(size)
			if frequency >= leftFreq && frequency <= rightFreq {
				filterbank[i][j] = (frequency - leftFreq) / (centerFreq - leftFreq)
			} else if frequency >= rightFreq && frequency <= centerFreq {
				filterbank[i][j] = (centerFreq - frequency) / (centerFreq - rightFreq)
			}
		}
	}

	return filterbank
}

func applyFilterbank(data []float64, filterbank [][]float64) []float64 {
	result := make([]float64, len(filterbank))
	for i := range filterbank {
		sum := 0.0
		for j, v := range data {
			sum += v * filterbank[i][j]
		}
		result[i] = sum
	}
	return result
}

func applyLog(data []float64) []float64 {
	result := make([]float64, len(data))
	for i, v := range data {
		result[i] = math.Log(v)
	}
	return result
}

func applyDCT(data []float64) []float64 {
	result := make([]float64, len(data))
	for i := range data {
		sum := 0.0
		for j, v := range data {
			sum += v * math.Cos(math.Pi/float64(len(data))*(float64(j)+0.5)*float64(i))
		}
		result[i] = sum
	}
	return result
}

func ExtractBase(filePath string) string {
	baseNameWithExtension := filepath.Base(filePath)
	baseName := strings.TrimSuffix(baseNameWithExtension, filepath.Ext(baseNameWithExtension))
	return baseName
}
