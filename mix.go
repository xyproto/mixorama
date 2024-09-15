package mix

import (
	"errors"
	"math"
)

// LinearSummation mixes multiple audio samples by adding them together.
// It automatically clamps the sum to avoid overflow and distortion.
func LinearSummation(samples ...[]int16) ([]int16, error) {
	if len(samples) == 0 {
		return nil, errors.New("no samples provided")
	}

	numSamples := len(samples[0])
	combined := make([]int16, numSamples)

	for i := 0; i < numSamples; i++ {
		sum := int32(0)
		for _, sample := range samples {
			if len(sample) != numSamples {
				return nil, errors.New("mismatched sample lengths")
			}
			sum += int32(sample[i])
		}
		// Clamp the result to avoid overflow
		if sum > math.MaxInt16 {
			sum = math.MaxInt16
		} else if sum < math.MinInt16 {
			sum = math.MinInt16
		}
		combined[i] = int16(sum)
	}

	return combined, nil
}

// WeightedSummation mixes multiple audio samples by applying a weight to each sample.
// Each sample's amplitude is scaled by its corresponding weight before summing.
func WeightedSummation(weights []float64, samples ...[]int16) ([]int16, error) {
	if len(weights) != len(samples) {
		return nil, errors.New("number of weights must match number of samples")
	}

	if len(samples) == 0 {
		return nil, errors.New("no samples provided")
	}

	numSamples := len(samples[0])
	combined := make([]int16, numSamples)

	for i := 0; i < numSamples; i++ {
		sum := float64(0)
		for j, sample := range samples {
			if len(sample) != numSamples {
				return nil, errors.New("mismatched sample lengths")
			}
			sum += float64(sample[i]) * weights[j]
		}
		// Clamp the result to avoid overflow
		if sum > math.MaxInt16 {
			sum = math.MaxInt16
		} else if sum < math.MinInt16 {
			sum = math.MinInt16
		}
		combined[i] = int16(sum)
	}

	return combined, nil
}

// RMSMixing correctly mixes audio samples using the Root Mean Square method.
func RMSMixing(samples ...[]int16) ([]int16, error) {
	if len(samples) == 0 {
		return nil, errors.New("no samples provided")
	}

	numSamples := len(samples[0])
	combined := make([]int16, numSamples)

	for i := 0; i < numSamples; i++ {
		sumSquares := float64(0)
		for _, sample := range samples {
			if len(sample) != numSamples {
				return nil, errors.New("mismatched sample lengths")
			}
			// Square the sample value and accumulate
			sumSquares += float64(sample[i]) * float64(sample[i])
		}
		// Calculate RMS by taking the square root of the mean of squares
		rms := math.Sqrt(sumSquares / float64(len(samples)))

		// Clamp the result to int16 range
		if rms > float64(math.MaxInt16) {
			rms = float64(math.MaxInt16)
		} else if rms < float64(math.MinInt16) {
			rms = float64(math.MinInt16)
		}
		combined[i] = int16(rms)
	}

	return combined, nil
}
