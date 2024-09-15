package main

import (
	"flag"
	"fmt"
	"log"
	"math"

	"github.com/xyproto/mixorama"
)

const version = "0.0.1"

// Simple low-pass filter to remove high frequencies
func lowPassFilter(samples []int16, sampleRate int, cutoffFrequency float64) []int16 {
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

// normalizeSamples scales the samples so the peak amplitude matches the given max amplitude
func normalizeSamples(samples []int16, targetPeak int16) []int16 {
	// Find the current peak amplitude
	currentPeak := findPeakAmplitude(samples)

	// Calculate scaling factor
	if currentPeak == 0 {
		return samples // Avoid division by zero
	}

	scale := float64(targetPeak) / float64(currentPeak)

	// Apply scaling to all samples
	normalizedSamples := make([]int16, len(samples))
	for i := range samples {
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

// findPeakAmplitude returns the maximum absolute amplitude in the sample set
func findPeakAmplitude(samples []int16) int16 {
	maxAmplitude := int16(0)
	for _, sample := range samples {
		if abs := int16(math.Abs(float64(sample))); abs > maxAmplitude {
			maxAmplitude = abs
		}
	}
	return maxAmplitude
}

func main() {
	// Define flags
	outputFile := flag.String("o", "combined.wav", "Specify the output file")
	showVersion := flag.Bool("version", false, "Show the version and exit")
	showHelp := flag.Bool("help", false, "Show help")

	// Parse flags
	flag.Parse()

	// Show version and exit if --version is passed
	if *showVersion {
		fmt.Printf("rms version %s\n", version)
		return
	}

	// Show help and exit if --help is passed
	if *showHelp {
		flag.Usage()
		return
	}

	// Expect at least two input files
	if flag.NArg() < 2 {
		fmt.Println("Usage: rms [options] <input1.wav> <input2.wav> [additional input files...]")
		flag.Usage()
		return
	}

	// Load the first input file to initialize the combined samples and sample rate
	inputFiles := flag.Args()
	firstFile := inputFiles[0]
	combined, sampleRate, err := mixorama.LoadWav(firstFile)
	if err != nil {
		log.Fatalf("Failed to load %s: %v", firstFile, err)
	}

	// Find the highest frequency across all files (initialize to 0) and track the loudest sample
	highestFrequency := 0.0
	loudestPeak := findPeakAmplitude(combined)

	// Process additional files and mix them using RMSMixing
	for _, inputFile := range inputFiles {
		// Load the next file
		wave, sr, err := mixorama.LoadWav(inputFile)
		if err != nil {
			log.Fatalf("Failed to load %s: %v", inputFile, err)
		}

		// Ensure the sample rate matches
		if sr != sampleRate {
			log.Fatalf("Sample rate mismatch between %s and %s", firstFile, inputFile)
		}

		// Determine the highest frequency in the current file
		currentHighestFrequency := analyzeHighestFrequency(wave, sr)
		if currentHighestFrequency > highestFrequency {
			highestFrequency = currentHighestFrequency
		}

		// Find the peak amplitude in the current file and track the loudest peak
		peak := findPeakAmplitude(wave)
		if peak > loudestPeak {
			loudestPeak = peak
		}

		// Pad the shorter sample with zeros
		combined, wave = mixorama.PadSamples(combined, wave)

		// Mix the current combined samples with the newly loaded samples
		combined, err = mixorama.RMSMixing(combined, wave)
		if err != nil {
			log.Fatalf("Error during RMS mixing of %s: %v", inputFile, err)
		}
	}

	// Apply low-pass filter using the highest detected frequency
	fmt.Printf("Applying low-pass filter with cutoff frequency: %.2f Hz\n", highestFrequency)
	combined = lowPassFilter(combined, sampleRate, highestFrequency)

	// Normalize the final combined samples to the loudest input sample's peak
	fmt.Printf("Normalizing loudness to the loudest peak: %d\n", loudestPeak)
	combined = normalizeSamples(combined, loudestPeak)

	// Save the final combined result to the output file
	if err := mixorama.SaveWav(*outputFile, combined, sampleRate); err != nil {
		log.Fatalf("Failed to save %s: %v", *outputFile, err)
	}

	fmt.Printf("Successfully mixed %d files into %s\n", len(inputFiles), *outputFile)
}

// analyzeHighestFrequency estimates the highest frequency in the audio signal.
func analyzeHighestFrequency(samples []int16, sampleRate int) float64 {
	// Simple estimation: check zero crossings
	zeroCrossings := 0
	for i := 1; i < len(samples); i++ {
		if (samples[i-1] > 0 && samples[i] < 0) || (samples[i-1] < 0 && samples[i] > 0) {
			zeroCrossings++
		}
	}

	// Highest frequency estimation based on zero crossings
	duration := float64(len(samples)) / float64(sampleRate)
	frequency := float64(zeroCrossings) / (2 * duration)

	return frequency
}
