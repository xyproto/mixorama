
# Mixorama

The `mixorama` package provides several ways of manipulating and mixing `[]int16` audio samples, including linear summation, weighted summation, and Root Mean Square (RMS) mixing. Additionally, it provides functions for loading and saving `.wav` files, padding audio samples, applying low-pass filtering, and normalizing the audio signal.

* NOTE: This package is experimental and a work in progress!

## Functions

### Mixing Functions

#### `func LinearSummation(samples ...[]int16) ([]int16, error)`
- **Description**:
    - This function adds multiple audio samples together. It automatically clamps the sum to ensure that it stays within the valid range of `int16` values, avoiding overflow and distortion.
- **Parameters**:
    - `samples`: A variable number of slices where each slice contains `int16` audio samples.
- **Returns**:
    - A slice of `int16` containing the combined audio samples.
    - An error if there are no input samples or if the lengths of the samples are mismatched.
- **Usage**:
    ```go
    combined, err := LinearSummation(wave1, wave2, wave3)
    ```

#### `func WeightedSummation(weights []float64, samples ...[]int16) ([]int16, error)`
- **Description**:
    - This function allows for weighted summation of multiple audio samples. Each sample is scaled by its corresponding weight before being summed together. This provides control over the relative volumes of each input.
- **Parameters**:
    - `weights`: A slice of `float64` values representing the weights for each input sample.
    - `samples`: A variable number of slices where each slice contains `int16` audio samples.
- **Returns**:
    - A slice of `int16` containing the combined audio samples after applying the weights.
    - An error if the number of weights does not match the number of samples, or if the sample lengths are mismatched.
- **Usage**:
    ```go
    weights := []float64{0.5, 0.8, 0.6}
    combined, err := WeightedSummation(weights, wave1, wave2, wave3)
    ```

#### `func RMSMixing(samples ...[]int16) ([]int16, error)`
- **Description**:
    - This function mixes audio samples using the Root Mean Square (RMS) method. It squares each sample, calculates the mean of the squares, and then takes the square root of the result. This technique helps provide a more balanced perception of loudness when mixing.
- **Parameters**:
    - `samples`: A variable number of slices where each slice contains `int16` audio samples.
- **Returns**:
    - A slice of `int16` containing the RMS-mixed audio samples.
    - An error if there are no input samples or if the lengths of the samples are mismatched.
- **Usage**:
    ```go
    combined, err := RMSMixing(wave1, wave2)
    ```

### Utility Functions

#### `func LoadWav(filename string) ([]int16, int, error)`
- **Description**:
    - Loads a `.wav` file and returns the audio samples as `[]int16` (stereo), along with the sample rate. If the file is mono, it duplicates the mono channel to create stereo output.
- **Parameters**:
    - `filename`: The path to the `.wav` file.
- **Returns**:
    - A slice of `int16` containing the audio samples.
    - The sample rate as an `int`.
    - An error if the file could not be loaded.
- **Usage**:
    ```go
    samples, sampleRate, err := LoadWav("input.wav")
    ```

#### `func SaveWav(filename string, samples []int16, sampleRate int) error`
- **Description**:
    - Saves a slice of `int16` audio samples as a `.wav` file.
- **Parameters**:
    - `filename`: The path where the `.wav` file will be saved.
    - `samples`: A slice of `int16` containing the audio samples.
    - `sampleRate`: The sample rate of the audio.
- **Returns**:
    - An error if the file could not be saved.
- **Usage**:
    ```go
    err := SaveWav("output.wav", samples, sampleRate)
    ```

#### `func PadSamples(wave1, wave2 []int16) ([]int16, []int16)`
- **Description**:
    - Pads the shorter sample with zeros (silence) so that both samples have the same length.
- **Parameters**:
    - `wave1`, `wave2`: Two slices of `int16` audio samples.
- **Returns**:
    - Two slices of `int16`, both with the same length after padding.
- **Usage**:
    ```go
    paddedWave1, paddedWave2 := PadSamples(wave1, wave2)
    ```

#### `func LowPassFilter(samples []int16, sampleRate int, cutoffFrequency float64) []int16`
- **Description**:
    - Applies a low-pass filter to remove high-frequency noise from the audio samples.
- **Parameters**:
    - `samples`: A slice of `int16` containing the audio samples.
    - `sampleRate`: The sample rate of the audio.
    - `cutoffFrequency`: The frequency above which audio will be filtered out.
- **Returns**:
    - A slice of `int16` containing the filtered audio samples.
- **Usage**:
    ```go
    filteredSamples := LowPassFilter(samples, 44100, 5000) // Low-pass filter with 5kHz cutoff
    ```

#### `func NormalizeSamples(samples []int16, targetPeak int16) []int16`
- **Description**:
    - Normalizes the audio samples so the peak amplitude matches the given `targetPeak`.
- **Parameters**:
    - `samples`: A slice of `int16` containing the audio samples.
    - `targetPeak`: The desired peak amplitude.
- **Returns**:
    - A slice of `int16` containing the normalized audio samples.
- **Usage**:
    ```go
    normalizedSamples := NormalizeSamples(samples, 30000)
    ```

#### `func FindPeakAmplitude(samples []int16) int16`
- **Description**:
    - Finds the peak amplitude in the audio samples.
- **Parameters**:
    - `samples`: A slice of `int16` containing the audio samples.
- **Returns**:
    - The peak amplitude as an `int16`.
- **Usage**:
    ```go
    peak := FindPeakAmplitude(samples)
    ```

#### `func AnalyzeHighestFrequency(samples []int16, sampleRate int) float64`
- **Description**:
    - Estimates the highest frequency in the audio signal by analyzing zero crossings.
- **Parameters**:
    - `samples`: A slice of `int16` containing the audio samples.
    - `sampleRate`: The sample rate of the audio.
- **Returns**:
    - The estimated highest frequency as a `float64`.
- **Usage**:
    ```go
    highestFrequency := AnalyzeHighestFrequency(samples, 44100)
    ```

## Example Use

```go
package main

import (
    "fmt"
    "github.com/xyproto/mixorama"
)

func main() {
    wave1 := []int16{1000, 2000, 3000}
    wave2 := []int16{1500, 2500, 3500}

    // Linear Summation
    combined, err := mixorama.LinearSummation(wave1, wave2)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Linear Summation:", combined)

    // Weighted Summation
    weights := []float64{0.5, 0.5}
    combined, err = mixorama.WeightedSummation(weights, wave1, wave2)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Weighted Summation:", combined)

    // RMS Mixing
    combined, err = mixorama.RMSMixing(wave1, wave2)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("RMS Mixing:", combined)
}
```

## General Info

- License: MIT
- Version: 0.1.0
