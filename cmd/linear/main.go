package main

import (
	"flag"
	"fmt"
	"log"
	"math"

	"github.com/xyproto/mixorama"
)

const version = "0.0.1"

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
	loudestPeak := mixorama.FindPeakAmplitude(combined)

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
		peak := mixorama.FindPeakAmplitude(wave)
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
	combined = mixorama.LowPassFilter(combined, sampleRate, 15000) // Cut off frequencies above 15kHz

	// Normalize the final combined samples based on the loudest peak value
	fmt.Printf("Normalizing combined file to match the loudest input peak: %d\n", loudestPeak)
	combined = mixorama.NormalizeSamples(combined, loudestPeak)

	// Save the final combined result to the output file
	if err := mixorama.SaveWav(*outputFile, combined, sampleRate); err != nil {
		log.Fatalf("Failed to save %s: %v", *outputFile, err)
	}

	fmt.Printf("Successfully mixed %d files into %s\n", len(inputFiles), *outputFile)
}
