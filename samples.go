package mixorama

import (
	"math"
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

// LoadWav loads a .wav file and returns its samples as []int16 (stereo) along with the sample rate.
// If the file is mono, it converts it to stereo by duplicating the mono channel to both the left and right channels.
func LoadWav(filename string) ([]int16, int, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, 0, err
	}
	defer f.Close()

	decoder := wav.NewDecoder(f)
	buffer, err := decoder.FullPCMBuffer()
	if err != nil {
		return nil, 0, err
	}

	intBuffer := buffer
	numChannels := intBuffer.Format.NumChannels

	if numChannels == 1 {
		// Convert mono to stereo by duplicating the mono channel
		l := len(intBuffer.Data)
		stereoSamples := make([]int16, l*2)
		for i := 0; i < l; i++ {
			monoSample := int16(intBuffer.Data[i])
			// Copy the mono sample to both left and right channels
			stereoSamples[2*i] = monoSample   // Left channel
			stereoSamples[2*i+1] = monoSample // Right channel
		}
		return stereoSamples, intBuffer.Format.SampleRate, nil
	}

	// If stereo, just convert to []int16 directly
	l := len(intBuffer.Data)
	stereoSamples := make([]int16, l)
	for i := 0; i < l; i++ {
		stereoSamples[i] = int16(intBuffer.Data[i])
	}

	return stereoSamples, intBuffer.Format.SampleRate, nil
}

// SaveWav saves a slice of int16 samples as a .wav file
func SaveWav(filename string, samples []int16, sampleRate int) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := wav.NewEncoder(f, sampleRate, 16, 1, 1)
	intBuffer := &audio.IntBuffer{
		Data:           make([]int, len(samples)),
		Format:         &audio.Format{SampleRate: sampleRate, NumChannels: 1},
		SourceBitDepth: 16,
	}
	for i, sample := range samples {
		intBuffer.Data[i] = int(sample)
	}

	if err := encoder.Write(intBuffer); err != nil {
		return err
	}
	return encoder.Close()
}

// PadSamples pads the shorter sample with zeros (silence) so that both samples have the same length.
func PadSamples(wave1, wave2 []int16) ([]int16, []int16) {
	length1 := len(wave1)
	length2 := len(wave2)

	if length1 == length2 {
		return wave1, wave2
	}

	if length1 < length2 {
		// Pad wave1 with zeros
		paddedWave1 := make([]int16, length2)
		copy(paddedWave1, wave1)
		return paddedWave1, wave2
	}

	// Pad wave2 with zeros
	paddedWave2 := make([]int16, length1)
	copy(paddedWave2, wave2)
	return wave1, paddedWave2
}

// LowPassFilter is a simple low-pass filter that can remove high frequencies
func LowPassFilter(samples []int16, sampleRate int, cutoffFrequency float64) []int16 {
	rc := 1.0 / (2.0 * math.Pi * cutoffFrequency)
	dt := 1.0 / float64(sampleRate)
	alpha := dt / (rc + dt)

	filteredSamples := make([]int16, len(samples))
	filteredSamples[0] = samples[0]

	for i := 1; i < len(samples); i++ {
		filteredSamples[i] = int16(float64(filteredSamples[i-1]) + alpha*(float64(samples[i])-float64(filteredSamples[i-1])))
	}

	return filteredSamples
}

// NormalizeSamples scales the samples so the peak amplitude matches the given max amplitude
func NormalizeSamples(samples []int16, targetPeak int16) []int16 {
	// Find the current peak amplitude
	currentPeak := FindPeakAmplitude(samples)

	// Calculate scaling factor
	if currentPeak == 0 {
		return samples // Avoid division by zero
	}

	scale := float64(targetPeak) / float64(currentPeak)

	l := len(samples)

	// Apply scaling to all samples
	normalizedSamples := make([]int16, l)
	for i := 0; i < l; i++ {
		normalized := float64(samples[i]) * scale
		if normalized > math.MaxInt16 {
			normalizedSamples[i] = math.MaxInt16
		} else if normalized < math.MinInt16 {
			normalizedSamples[i] = math.MinInt16
		} else {
			normalizedSamples[i] = int16(normalized)
		}
	}

	return normalizedSamples
}

// FindPeakAmplitude returns the maximum absolute amplitude in the sample set
func FindPeakAmplitude(samples []int16) int16 {
	maxAmplitude := int16(0)
	for _, sample := range samples {
		if abs := int16(math.Abs(float64(sample))); abs > maxAmplitude {
			maxAmplitude = abs
		}
	}
	return maxAmplitude
}

// AnalyzeHighestFrequency estimates the highest frequency in the audio signal
func AnalyzeHighestFrequency(samples []int16, sampleRate int) float64 {
	// Simple estimation: check zero crossings
	zeroCrossings := 0
	l := len(samples)
	for i := 1; i < l; i++ {
		if (samples[i-1] > 0 && samples[i] < 0) || (samples[i-1] < 0 && samples[i] > 0) {
			zeroCrossings++
		}
	}

	// Highest frequency estimation based on zero crossings
	duration := float64(l) / float64(sampleRate)
	frequency := float64(zeroCrossings) / (2 * duration)

	return frequency
}
