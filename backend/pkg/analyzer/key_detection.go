package analyzer

import (
	"math"
	"math/cmplx"

	"gonum.org/v1/gonum/dsp/fourier"
)

// MajorKeyProfiles - Krumhansl-Schmuckler
var MajorKeyProfiles = []float64{6.35, 2.23, 3.48, 2.33, 4.38, 4.09, 2.52, 5.19, 2.39, 3.66, 2.29, 2.88}

// MinorKeyProfiles - Krumhansl-Schmuckler
var MinorKeyProfiles = []float64{6.33, 2.68, 3.52, 5.38, 2.60, 3.53, 2.54, 4.75, 3.98, 2.69, 3.34, 3.17}

var NoteNames = []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}

// DetectKey estimates the musical key of the audio samples.
func DetectKey(samples []float32, sampleRate int) (string, error) {
	// 1. Calculate Chromagram
	chroma := calculateChromagram(samples, sampleRate)

	// 2. Correlate with profiles
	bestKey := ""
	maxCorr := -1.0

	// Check Major keys
	for i := 0; i < 12; i++ {
		profile := rotate(MajorKeyProfiles, i)
		corr := correlation(chroma, profile)
		if corr > maxCorr {
			maxCorr = corr
			bestKey = NoteNames[i] + " Major"
		}
	}

	// Check Minor keys
	for i := 0; i < 12; i++ {
		profile := rotate(MinorKeyProfiles, i)
		corr := correlation(chroma, profile)
		if corr > maxCorr {
			maxCorr = corr
			bestKey = NoteNames[i] + " Minor"
		}
	}

	return bestKey, nil
}

func calculateChromagram(samples []float32, sampleRate int) []float64 {
	// Use a simplified approach: FFT on chunks, map bins to chroma
	// This is computationally intensive, so we process a subset or downsample
	// For this app, let's take a 30s slice from the middle if possible, or just chunks.

	fftSize := 4096
	chroma := make([]float64, 12)

	// cast to float64
	input := make([]float64, len(samples))
	for i, v := range samples {
		input[i] = float64(v)
	}

	numChunks := len(input) / fftSize
	if numChunks > 100 {
		numChunks = 100 // limit to 100 chunks to save time
	}

	step := len(input) / numChunks
	if step < fftSize {
		step = fftSize
	}

	fft := fourier.NewFFT(fftSize)

	for i := 0; i < numChunks; i++ {
		start := i * step
		if start+fftSize > len(input) {
			break
		}

		chunk := input[start : start+fftSize]
		coeffs := fft.Coefficients(nil, chunk)

		for j, c := range coeffs {
			mag := cmplx.Abs(c)
			freq := float64(j) * float64(sampleRate) / float64(fftSize)

			if freq < 27.5 { // Skip below A0
				continue
			}
			if freq > 4186 { // Skip above C8
				continue
			}

			// Map frequency to pitch class
			// MIDI note = 69 + 12 * log2(freq / 440)
			midiNote := 69 + 12*math.Log2(freq/440.0)
			pitchClass := int(math.Round(midiNote)) % 12
			if pitchClass < 0 {
				pitchClass += 12
			}

			chroma[pitchClass] += mag
		}
	}

	// Normalize?
	// optional, correlation handles scaling usually, but good for debug

	return chroma
}

func rotate(slice []float64, shift int) []float64 {
	rotated := make([]float64, len(slice))
	for i := 0; i < len(slice); i++ {
		rotated[i] = slice[(i-shift+12)%12]
	}
	return rotated
}

func correlation(a, b []float64) float64 {
	n := float64(len(a))
	sumA, sumB, sumAA, sumBB, sumAB := 0.0, 0.0, 0.0, 0.0, 0.0

	for i := 0; i < len(a); i++ {
		sumA += a[i]
		sumB += b[i]
		sumAA += a[i] * a[i]
		sumBB += b[i] * b[i]
		sumAB += a[i] * b[i]
	}

	// Pearson correlation
	numerator := sumAB - (sumA*sumB)/n
	denominator := math.Sqrt((sumAA - (sumA*sumA)/n) * (sumBB - (sumB*sumB)/n))

	if denominator == 0 {
		return 0
	}
	return numerator / denominator
}
