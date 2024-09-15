package mix

import (
	"math"
	"testing"
)

// Helper function to create a simple waveform for testing
func createTestWaveform(value int16, numSamples int) []int16 {
	waveform := make([]int16, numSamples)
	for i := 0; i < numSamples; i++ {
		waveform[i] = value
	}
	return waveform
}

// TestLinearSummation checks if the linear summation mixing works as expected
func TestLinearSummation(t *testing.T) {
	wave1 := createTestWaveform(1000, 10)
	wave2 := createTestWaveform(2000, 10)
	expected := createTestWaveform(3000, 10)

	result, err := LinearSummation(wave1, wave2)
	if err != nil {
		t.Fatalf("Error in LinearSummation: %v", err)
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("LinearSummation failed at index %d: expected %d, got %d", i, expected[i], v)
		}
	}
}

// TestWeightedSummation checks if the weighted summation mixing works as expected
func TestWeightedSummation(t *testing.T) {
	wave1 := createTestWaveform(1000, 10)
	wave2 := createTestWaveform(2000, 10)
	weights := []float64{0.5, 0.5}
	expected := createTestWaveform(1500, 10)

	result, err := WeightedSummation(weights, wave1, wave2)
	if err != nil {
		t.Fatalf("Error in WeightedSummation: %v", err)
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("WeightedSummation failed at index %d: expected %d, got %d", i, expected[i], v)
		}
	}
}

// TestRMSMixing checks if the RMS mixing works as expected
func TestRMSMixing(t *testing.T) {
	wave1 := createTestWaveform(1000, 10)
	wave2 := createTestWaveform(2000, 10)

	// Calculate the expected RMS value
	expectedRMS := int16(math.Sqrt((1000*1000 + 2000*2000) / 2))

	result, err := RMSMixing(wave1, wave2)
	if err != nil {
		t.Fatalf("Error in RMSMixing: %v", err)
	}

	for i, v := range result {
		if v != expectedRMS {
			t.Errorf("RMSMixing failed at index %d: expected %d, got %d", i, expectedRMS, v)
		}
	}
}

// TestErrorCases tests that the functions handle error cases correctly
func TestErrorCases(t *testing.T) {
	// Mismatched sample lengths
	wave1 := createTestWaveform(1000, 10)
	wave2 := createTestWaveform(1000, 5)

	_, err := LinearSummation(wave1, wave2)
	if err == nil {
		t.Error("Expected error for mismatched sample lengths in LinearSummation")
	}

	_, err = WeightedSummation([]float64{0.5, 0.5}, wave1, wave2)
	if err == nil {
		t.Error("Expected error for mismatched sample lengths in WeightedSummation")
	}

	_, err = RMSMixing(wave1, wave2)
	if err == nil {
		t.Error("Expected error for mismatched sample lengths in RMSMixing")
	}

	// Mismatched weights
	wave3 := createTestWaveform(1000, 10)
	_, err = WeightedSummation([]float64{0.5}, wave1, wave3)
	if err == nil {
		t.Error("Expected error for mismatched number of weights and samples in WeightedSummation")
	}
}
