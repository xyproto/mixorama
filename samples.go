package mixorama

import (
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
