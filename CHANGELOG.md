# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **Audio Context** — singleton manager, configurable sample rate and channels
- **Player** — io.Reader PCM stream playback with play/pause/stop/volume
- **Driver interface** — platform abstraction for WASAPI/CoreAudio/PulseAudio
- **NullDriver** — testing and headless mode (discards output)
- **WAV decoder** — Pure Go, 8/16/24/32-bit PCM + 32-bit float, stereo, duration
- **Mixer** — multi-channel mixing, per-source + master volume, clamping
- 42 tests, 94.8% coverage
