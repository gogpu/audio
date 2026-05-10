// Example: Mozart's "Eine Kleine Nachtmusik" K.525 — Citizen watch style.
//
// Pure Go WASAPI audio engine, zero CGO.
//
// Run:
//
//	go run ./examples/mozart
package main

import (
	"fmt"
	"math"
	"time"

	"github.com/gogpu/audio"
)

const (
	A3 = 220.00
	B3 = 246.94
	C4 = 261.63
	D4 = 293.66
	E4 = 329.63
	F4 = 349.23
	G4 = 392.00
	A4 = 440.00
	B4 = 493.88
	C5 = 523.25
	D5 = 587.33
	E5 = 659.25
	F5 = 698.46
	G5 = 783.99
	R  = 0
)

type note struct {
	f float64
	d float64
}

func main() {
	fmt.Println("♪ Mozart — Eine Kleine Nachtmusik")
	fmt.Println("  Citizen watch style · Pure Go WASAPI")
	fmt.Println()

	bpm := 138.0

	melody := []note{
		// Bar 1: C5(q) r G4(8) C5(q) r G4(8) — from MusicXML
		{C5, 1}, {R, 0.5}, {G4, 0.5},
		{C5, 1}, {R, 0.5}, {G4, 0.5},
		// Bar 2: C5 G4 C5 E5 G5(q) r
		{C5, 0.5}, {G4, 0.5}, {C5, 0.5}, {E5, 0.5},
		{G5, 1}, {R, 1},
		// Bar 3: F5(q) r D5(8) F5(q) r D5(8)
		{F5, 1}, {R, 0.5}, {D5, 0.5},
		{F5, 1}, {R, 0.5}, {D5, 0.5},
		// Bar 4: F5 D5 B4 D5 G4(q) r
		{F5, 0.5}, {D5, 0.5}, {B4, 0.5}, {D5, 0.5},
		{G4, 1}, {R, 1},
		// Bar 5: C5(8) r C5(dq) E5 D5 C5
		{C5, 0.5}, {R, 0.5}, {C5, 1.5}, {E5, 0.5},
		{D5, 0.5}, {C5, 0.5},
		// Bar 6: C5 B4 B4(dq) D5 F5 B4
		{C5, 0.5}, {B4, 0.5}, {B4, 1.5}, {D5, 0.5},
		{F5, 0.5}, {B4, 0.5},
		// Bar 7: D5 C5 C5(dq) E5 D5 C5
		{D5, 0.5}, {C5, 0.5}, {C5, 1.5}, {E5, 0.5},
		{D5, 0.5}, {C5, 0.5},
		// Bar 8: C5 B4 B4(dq) D5 F5 B4
		{C5, 0.5}, {B4, 0.5}, {B4, 1.5}, {D5, 0.5},
		{F5, 0.5}, {B4, 0.5},
		// Bar 9: running 16ths
		{C5, 0.5}, {C5, 0.5},
		{C5, 0.25}, {B4, 0.25}, {A4, 0.25}, {B4, 0.25},
		{C5, 0.5}, {C5, 0.5},
		{E5, 0.25}, {D5, 0.25}, {C5, 0.25}, {D5, 0.25},
		// Bar 10: E5 E5 ... G5(q) r
		{E5, 0.5}, {E5, 0.5},
		{G5, 0.25}, {F5, 0.25}, {E5, 0.25}, {F5, 0.25},
		{G5, 1}, {R, 1},
	}

	wav := renderMelody(melody, bpm, 44100)

	ctx, err := audio.NewContext()
	if err != nil {
		fmt.Printf("Audio error: %v\n", err)
		return
	}
	defer ctx.Close()

	player, err := ctx.PlayWAV(wav)
	if err != nil {
		fmt.Printf("PlayWAV error: %v\n", err)
		return
	}

	fmt.Print("Playing")
	for player.IsPlaying() {
		fmt.Print(" ♪")
		time.Sleep(300 * time.Millisecond)
	}
	fmt.Println("\nBravo! 🎹")
}

func renderMelody(notes []note, bpm float64, sr int) []byte {
	beatDur := 60.0 / bpm
	ch := 2
	total := 0.0
	for _, n := range notes {
		total += n.d * beatDur
	}
	samples := int(total * float64(sr))
	dataSize := samples * ch * 4
	buf := make([]byte, 44+dataSize)

	copy(buf[0:], "RIFF")
	le32(buf[4:], uint32(len(buf)-8))
	copy(buf[8:], "WAVE")
	copy(buf[12:], "fmt ")
	le32(buf[16:], 16)
	le16(buf[20:], 3)
	le16(buf[22:], uint16(ch))
	le32(buf[24:], uint32(sr))
	le32(buf[28:], uint32(sr*ch*4))
	le16(buf[32:], uint16(ch*4))
	le16(buf[34:], 32)
	copy(buf[36:], "data")
	le32(buf[40:], uint32(dataSize))

	off := 44
	for _, n := range notes {
		ns := int(n.d * beatDur * float64(sr))
		for i := 0; i < ns; i++ {
			var v float32
			if n.f > 0 {
				t := float64(i) / float64(sr)
				phase := math.Mod(t*n.f, 1.0)
				if phase < 0.45 {
					v = 0.22
				} else {
					v = -0.22
				}
				pos := float64(i) / float64(ns)
				v *= float32(math.Min(pos*40, 1.0) * math.Min((1-pos)*25, 1.0))
			}
			putF32(buf[off:], v)
			off += 4
			putF32(buf[off:], v)
			off += 4
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
func putF32(b []byte, v float32) { le32(b, math.Float32bits(v)) }
