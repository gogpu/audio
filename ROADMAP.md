# Roadmap

## Vision

**audio** is a Pure Go audio engine for the [gogpu](https://github.com/gogpu) ecosystem — PCM output, decoding, mixing, with platform-native drivers. Zero CGO.

## Current: Unreleased

- Audio Context (singleton, sample rate/channels config)
- Player (io.Reader, play/pause/stop/volume)
- Driver interface + NullDriver
- WAV decoder (8/16/24/32-bit PCM, 32-bit float)
- Mixer (multi-channel, volume, clamping)

## Planned

### v0.1.0 — First Release

- [ ] WASAPI driver (Windows)
- [ ] CoreAudio driver (macOS)
- [ ] PulseAudio driver (Linux, Pure Go native protocol)
- [ ] CI workflow (3 platforms)
- [ ] Integration test with real audio device

### v0.2.0 — Formats + Features

- [ ] OGG Vorbis decoder (Pure Go)
- [ ] MP3 decoder (Pure Go)
- [ ] Streaming playback (large files)
- [ ] 3D spatial audio (basic panning)

### Future

- [ ] Audio recording/capture
- [ ] MIDI support
- [ ] Web Audio API (WASM target)
- [ ] PipeWire native driver (no PulseAudio compatibility layer)
