Lem-in
------

Lem-in is a program/library that takes a file describing an ant farm -- defined by a number of ants, rooms, and tunnels -- and calculates the most efficient way to move all ants from the start room to the end room in the fewest possible turns. Each tunnel can be used only once per turn, and rooms (except the start and end) can hold only one ant at a time.

The library can either be run in native platform-agnostic Go or in Wasm. 

- Program is written in Go, with core logic located in the `core/` directory.
- Sample files describing various ant farms are available in the `samples/` directory.
- Native command-line application in the `cli/` directory.
- WebAssembly bridge in the `wasm/` directory.
- Some utilities in the `utils/` directory.

To build the project, ensure that you have `npm`, `make`, and `go` installed.
