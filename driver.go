// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package audio

// Driver is the platform audio output interface.
// Implementations exist per platform: WASAPI (Windows), CoreAudio (macOS),
// PulseAudio (Linux). The NullDriver is provided for testing and headless use.
type Driver interface {
	// Open opens the audio device with the given sample rate and channel count.
	Open(sampleRate, channels int) error

	// Write sends interleaved float32 PCM samples to the audio device.
	// It blocks until the hardware consumes the samples.
	Write(samples []float32) (int, error)

	// Close releases the audio device.
	Close() error
}
