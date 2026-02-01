package analyzer

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/benjojo/bpm"
	"github.com/go-audio/wav"
)

type AnalysisResult struct {
	BPM float64 `json:"bpm"`
	Key string  `json:"key"`
}

// Analyze processes the audio file at the given path and returns the extracted metrics.
func Analyze(inputPath string) (*AnalysisResult, error) {
	// 1. Convert to WAV (Mono, 44.1kHz)
	tempWav := inputPath + ".temp.wav"
	if err := convertToWav(inputPath, tempWav); err != nil {
		return nil, fmt.Errorf("failed to convert audio: %w", err)
	}
	defer os.Remove(tempWav)

	// 2. Read the WAV file
	f, err := os.Open(tempWav)
	if err != nil {
		return nil, fmt.Errorf("failed to open converted wav: %w", err)
	}
	defer f.Close()

	decoder := wav.NewDecoder(f)
	if !decoder.IsValidFile() {
		return nil, fmt.Errorf("invalid wav file")
	}

	buf, err := decoder.FullPCMBuffer()
	if err != nil {
		return nil, fmt.Errorf("failed to read wav buffer: %w", err)
	}

	// 3. Detect BPM
	floats32 := make([]float32, len(buf.Data))
	for i, sample := range buf.Data {
		floats32[i] = float32(sample)
	}

	nrg := bpm.ReadFloatArray(floats32)
	// Scan from 60 to 200 BPM
	detectedBPM := bpm.ScanForBpm(nrg, 60, 200, 1024, 2048)

	// 4. Detect Key
	detectedKey, err := DetectKey(floats32, int(decoder.SampleRate))
	if err != nil {
		fmt.Printf("Key detection failed: %v\n", err)
		detectedKey = "Unknown"
	}

	return &AnalysisResult{
		BPM: detectedBPM,
		Key: detectedKey,
	}, nil
}

func convertToWav(input, output string) error {
	cmd := exec.Command("ffmpeg", "-y", "-i", input, "-ac", "1", "-ar", "44100", "-f", "wav", output)
	return cmd.Run()
}
