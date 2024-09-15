package mixorama

import (
	"math"
	"os"
	"testing"
)

func TestLoadWav(t *testing.T) {
	// Assuming you have a test.wav file available in your test directory for this test
	_, sampleRate, err := LoadWav("test.wav")
	if err != nil {
		t.Fatalf("Failed to load test.wav: %v", err)
	}
	if sampleRate <= 0 {
		t.Errorf("Expected valid sample rate, got %d", sampleRate)
	}
}

func TestSaveWav(t *testing.T) {
	// Create test samples
	samples := []int16{1000, -1000, 2000, -2000}
	filename := "test_output.wav"
	defer os.Remove(filename) // Cleanup after test

	// Save the samples to a file
	err := SaveWav(filename, samples, 44100)
	if err != nil {
		t.Fatalf("Failed to save WAV file: %v", err)
	}

	// Check if the file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatalf("Expected file to exist: %s", filename)
	}
}

func TestPadSamples(t *testing.T) {
	wave1 := []int16{100, 200, 300}
	wave2 := []int16{400, 500}
	padded1, padded2 := PadSamples(wave1, wave2)

	if len(padded1) != len(padded2) {
		t.Errorf("Expected padded waves to have the same length, got %d and %d", len(padded1), len(padded2))
	}

	// Ensure the shorter wave is padded with zeros
	if padded2[2] != 0 {
		t.Errorf("Expected padding to be zero, got %d", padded2[2])
	}
}

func TestLowPassFilter(t *testing.T) {
	samples := []int16{100, 200, 300, 400, 500}
	filtered := LowPassFilter(samples, 44100, 1000) // Apply a low-pass filter with 1kHz cutoff

	if len(filtered) != len(samples) {
		t.Errorf("Expected filtered samples to have the same length, got %d and %d", len(filtered), len(samples))
	}

	// Check that the filter has smoothed the samples
	if math.Abs(float64(filtered[1])-float64(filtered[0])) > 100 {
		t.Errorf("Low-pass filter did not smooth the values as expected")
	}
}

func TestNormalizeSamples(t *testing.T) {
	samples := []int16{100, 200, -300}
	targetPeak := int16(1000)
	normalized := NormalizeSamples(samples, targetPeak)

	// Check that the peak amplitude matches the target
	peak := FindPeakAmplitude(normalized)
	if peak != targetPeak {
		t.Errorf("Expected peak amplitude %d, got %d", targetPeak, peak)
	}
}

func TestFindPeakAmplitude(t *testing.T) {
	samples := []int16{100, 200, -300}
	expectedPeak := int16(300)
	peak := FindPeakAmplitude(samples)

	if peak != expectedPeak {
		t.Errorf("Expected peak amplitude %d, got %d", expectedPeak, peak)
	}
}

func TestAnalyzeHighestFrequency(t *testing.T) {
	// Simple samples with alternating values for zero-crossing detection
	samples := []int16{1000, -1000, 1000, -1000}
	sampleRate := 44100
	frequency := AnalyzeHighestFrequency(samples, sampleRate)

	// Check if the estimated frequency is reasonable for the test case
	if frequency <= 0 {
		t.Errorf("Expected a positive frequency, got %.2f", frequency)
	}
}
