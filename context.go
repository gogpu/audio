// Copyright 2026 The gogpu Authors
// SPDX-License-Identifier: MIT

package audio

import (
	"errors"
	"io"
	"sync"
)

const (
	defaultSampleRate = 44100
	defaultChannels   = 2
)

// Context manages audio output for an application. It owns the platform
// audio device and coordinates player lifecycle and mixing. Typically
// one Context exists per application.
type Context struct {
	mu         sync.Mutex
	sampleRate int
	channels   int
	driver     Driver
	mixer      *Mixer
	players    []*Player
	running    bool
	stopCh     chan struct{}
	mixBuf     []float32
}

// Option configures the audio context during construction.
type Option func(*Context)

// WithSampleRate sets the audio sample rate in Hz. Default is 44100.
func WithSampleRate(rate int) Option {
	return func(c *Context) {
		if rate > 0 {
			c.sampleRate = rate
		}
	}
}

// WithChannels sets the number of output channels. Default is 2 (stereo).
func WithChannels(ch int) Option {
	return func(c *Context) {
		if ch > 0 {
			c.channels = ch
		}
	}
}

// WithDriver sets a specific audio driver. If not set, NullDriver is used
// until platform drivers are implemented.
func WithDriver(d Driver) Option {
	return func(c *Context) {
		c.driver = d
	}
}

// NewContext creates an audio context and opens the audio device.
func NewContext(opts ...Option) (*Context, error) {
	c := &Context{
		sampleRate: defaultSampleRate,
		channels:   defaultChannels,
		mixer:      NewMixer(),
		stopCh:     make(chan struct{}),
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.driver == nil {
		c.driver = &NullDriver{}
	}

	if err := c.driver.Open(c.sampleRate, c.channels); err != nil {
		return nil, errors.New("audio: failed to open driver: " + err.Error())
	}

	// Buffer holds one callback worth of samples.
	// 2048 frames * channels = common low-latency buffer size.
	const framesPerBuffer = 2048
	c.mixBuf = make([]float32, framesPerBuffer*c.channels)

	c.running = true
	go c.outputLoop()

	return c, nil
}

// Close stops the audio output loop and releases the audio device.
func (c *Context) Close() error {
	c.mu.Lock()
	if !c.running {
		c.mu.Unlock()
		return nil
	}
	c.running = false
	c.mu.Unlock()

	close(c.stopCh)

	return c.driver.Close()
}

// NewPlayer creates a player from a PCM audio stream. The reader should
// provide interleaved float32 PCM samples encoded as little-endian bytes
// (4 bytes per sample, same layout as WAVDecoder.Read output).
func (c *Context) NewPlayer(src io.Reader) *Player {
	p := &Player{
		ctx:    c,
		src:    src,
		volume: 1.0,
	}

	c.mu.Lock()
	c.players = append(c.players, p)
	c.mixer.AddSource(p)
	c.mu.Unlock()

	return p
}

// PlayWAV decodes WAV data and starts playback immediately.
// Returns the player for volume/pause/stop control.
func (c *Context) PlayWAV(data []byte) (*Player, error) {
	dec, err := DecodeWAV(data)
	if err != nil {
		return nil, err
	}
	p := c.NewPlayer(dec)
	p.Play()
	return p, nil
}

// SampleRate returns the context's configured sample rate.
func (c *Context) SampleRate() int {
	return c.sampleRate
}

// Channels returns the context's configured channel count.
func (c *Context) Channels() int {
	return c.channels
}

// outputLoop runs in a goroutine, continuously mixing active players
// and writing the result to the driver.
func (c *Context) outputLoop() {
	for {
		select {
		case <-c.stopCh:
			return
		default:
		}

		c.mixer.Mix(c.mixBuf)

		if _, err := c.driver.Write(c.mixBuf); err != nil {
			// Driver write failed -- stop the loop.
			// In production drivers, this typically means the device was lost.
			return
		}
	}
}
