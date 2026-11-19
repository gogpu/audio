# apigolden — the public-API golden of audio

`apigolden` prints every exported declaration of `github.com/gogpu/audio` — types, fields, interface method
sets, functions, methods, constants and variables, with full signatures — as deterministic text. The
committed files under `testdata/` are the frozen public shape; a change to them is a change to a contract
consumers depend on.

## Regenerating and checking

```bash
GOWORK=off go -C tools run ./apigolden -out apigolden/testdata/api_windows.txt
GOWORK=off go -C tools run ./apigolden -goos linux -out apigolden/testdata/api_unix.txt
git diff --exit-code tools/apigolden/testdata/   # must be empty
```

## The two contexts

The package has two driver build contexts: the WASAPI driver on Windows and the null driver everywhere
else (`driver_default_windows.go`, `driver_default_other.go`). Both are generated and committed; today
their exported surface is identical (1 package / 14 identifiers), and the two files keep that a checked
fact rather than an assumption. The `linux` golden stands for every `!windows` platform.

## The floor

`contextFloors` pins the package and identifier counts the loader must reach before a golden is written
(1 package / 14 identifiers in both contexts). A load that loses the package or drops identifiers fails
instead of producing a short file that looks reviewed.

## How the loader works

The tool uses only the standard library: `go list -export` produces export data with the running
toolchain and `go/importer` decodes it with that same toolchain, so the tool cannot fall out of step with
the Go version that invokes it. The package under analysis is parsed from source and type-checked in the
tool, so parameter names survive into the golden.

The formatter is covered by unit tests: `GOWORK=off go -C tools test ./apigolden`.
