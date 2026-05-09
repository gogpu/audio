package main

import (
	"fmt"
	"math"
	"time"

	"github.com/gogpu/audio"
)

func main() {
	fmt.Println("GoGPU Audio Engine — WASAPI Test")

	// Generate a 440Hz sine wave WAV (1 second, stereo, 44100Hz, float32)
	sampleRate := 44100
	duration := 1.0
	freq := 440.0
	samples := int(float64(sampleRate) * duration)

	// WAV header for float32 stereo
	wav := makeWAV(sampleRate, 2, samples, func(i int) (float32, float32) {
		t := float64(i) / float64(sampleRate)
		v := float32(0.3 * math.Sin(2*math.Pi*freq*t))
		return v, v
	})

	ctx, err := audio.NewContext()
	if err != nil {
		fmt.Printf("NewContext error: %v\n", err)
		return
	}
	defer ctx.Close()

	player, err := ctx.PlayWAV(wav)
	if err != nil {
		fmt.Printf("PlayWAV error: %v\n", err)
		return
	}

	fmt.Printf("Playing 440Hz sine wave for 1 second...\n")
	for player.IsPlaying() {
		time.Sleep(50 * time.Millisecond)
	}
	fmt.Println("Done!")
}

func makeWAV(sampleRate, channels, samples int, gen func(int) (float32, float32)) []byte {
	dataSize := samples * channels * 4
	fileSize := 44 + dataSize

	buf := make([]byte, fileSize)

	// RIFF header
	copy(buf[0:4], "RIFF")
	le32(buf[4:], uint32(fileSize-8))
	copy(buf[8:12], "WAVE")

	// fmt chunk
	copy(buf[12:16], "fmt ")
	le32(buf[16:], 16) // chunk size
	le16(buf[20:], 3)  // format = IEEE float
	le16(buf[22:], uint16(channels))
	le32(buf[24:], uint32(sampleRate))
	le32(buf[28:], uint32(sampleRate*channels*4)) // byte rate
	le16(buf[32:], uint16(channels*4))            // block align
	le16(buf[34:], 32)                            // bits per sample

	// data chunk
	copy(buf[36:40], "data")
	le32(buf[40:], uint32(dataSize))

	offset := 44
	for i := 0; i < samples; i++ {
		l, r := gen(i)
		putFloat32(buf[offset:], l)
		offset += 4
		if channels == 2 {
			putFloat32(buf[offset:], r)
			offset += 4
		}
	}

	return buf
}

func le16(b []byte, v uint16) { b[0] = byte(v); b[1] = byte(v >> 8) }
func le32(b []byte, v uint32) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
}
func putFloat32(b []byte, v float32) { le32(b, math.Float32bits(v)) }
