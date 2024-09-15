# Mixorama

The `mixorama` package provides different ways of mixing `[]int16` audio samples.

A couple of different mixing strategies are provided: linear summation, weighted summation and Root Mean Square (RMS).

* NOTE: this package is experimental and a work in progress!

## Functions

### `func LinearSummation(samples ...[]int16) ([]int16, error)`
- **Description**:
    - This function adds multiple audio samples together. It automatically clamps the sum to ensure that it stays within the valid range of `int16` values, thus avoiding overflow and distortion.
- **Parameters**:
    - `samples`: A variable number of slices where each slice contains `int16` audio samples.
- **Returns**:
    - A slice of `int16` containing the combined audio samples.
    - An error if there are no input samples or if the lengths of the samples are mismatched.
- **Usage**:
    ```go
    combined, err := LinearSummation(wave1, wave2, wave3)
    ```

### `func WeightedSummation(weights []float64, samples ...[]int16) ([]int16, error)`
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

### `func RMSMixing(samples ...[]int16) ([]int16, error)`
- **Description**:
    - This function uses the Root Mean Square (RMS) method to mix multiple audio samples. It squares each sample, calculates the mean of the squares, and then takes the square root of the result. This technique helps provide a more balanced perception of loudness when mixing.
- **Parameters**:
    - `samples`: A variable number of slices where each slice contains `int16` audio samples.
- **Returns**:
    - A slice of `int16` containing the RMS-mixed audio samples.
    - An error if there are no input samples or if the lengths of the samples are mismatched.
- **Usage**:
    ```go
    combined, err := RMSMixing(wave1, wave2)
    ```

## Example use

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

## General info

* License: MIT
* Version: 0.0.1
