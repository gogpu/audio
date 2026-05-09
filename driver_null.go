// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package audio

// NullDriver discards all audio output.
// It is used for testing and headless environments where no audio
// hardware is available.
type NullDriver struct{}

// Open is a no-op for the null driver.
func (d *NullDriver) Open(sampleRate, channels int) error { return nil }

// Write accepts and discards all samples, returning the full count.
func (d *NullDriver) Write(samples []float32) (int, error) { return len(samples), nil }

// Close is a no-op for the null driver.
func (d *NullDriver) Close() error { return nil }
