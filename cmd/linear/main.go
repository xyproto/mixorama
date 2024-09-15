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

// normalizeSamples scales the combined samples so that the peak amplitude matches the target peak amplitude
func normalizeSamples(samples []int16, targetPeak int16) []int16 {
	currentPeak := findPeakAmplitude(samples)
	if currentPeak == 0 || targetPeak == 0 {
		fmt.Println("Skipping normalization due to zero or invalid peak amplitude.")
		return samples // Avoid division by zero or scaling silent audio
	}

	// Calculate scaling factor based on peak amplitude
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

	// Find the loudest peak across all input files
	loudestPeak := findPeakAmplitude(combined)

	// Process additional files and mix them using weighted summation
	for _, inputFile := range inputFiles[1:] {
		// Load the next file
		wave, sr, err := mixorama.LoadWav(inputFile)
		if err != nil {
			log.Fatalf("Failed to load %s: %v", inputFile, err)
		}

		// Ensure the sample rate matches
		if sr != sampleRate {
			log.Fatalf("Sample rate mismatch between %s and %s", firstFile, inputFile)
		}

		// Find the peak amplitude in the current file and track the loudest peak
		peak := findPeakAmplitude(wave)
		if peak > loudestPeak {
			loudestPeak = peak
		}

		// Pad the shorter sample with zeros
		combined, wave = mixorama.PadSamples(combined, wave)

		// Perform weighted summation (reduce contribution of each input to avoid clipping)
		for i := 0; i < len(combined); i++ {
			weightedSum := (int32(combined[i]) + int32(wave[i])) / 2 // Adjust weight (averaging)
			if weightedSum > math.MaxInt16 {
				combined[i] = math.MaxInt16
			} else if weightedSum < math.MinInt16 {
				combined[i] = math.MinInt16
			} else {
				combined[i] = int16(weightedSum)
			}
		}
	}

	// Apply low-pass filter using a reasonable cutoff frequency (e.g., 15kHz to remove high-frequency noise)
	fmt.Println("Applying low-pass filter to combined audio.")
	combined = lowPassFilter(combined, sampleRate, 15000) // Cut off frequencies above 15kHz

	// Normalize the final combined samples based on the loudest peak value
	fmt.Printf("Normalizing combined file to match the loudest input peak: %d\n", loudestPeak)
	combined = normalizeSamples(combined, loudestPeak)

	// Save the final combined result to the output file
	if err := mixorama.SaveWav(*outputFile, combined, sampleRate); err != nil {
		log.Fatalf("Failed to save %s: %v", *outputFile, err)
	}

	fmt.Printf("Successfully mixed %d files into %s\n", len(inputFiles), *outputFile)
}
