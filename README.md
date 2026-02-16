Lem-in
------

Lem-in is a program that takes a file describing an ant farm -- defined by a number of ants, rooms, and tunnels -- and calculates the most efficient way to move all ants from the start room to the end room in the fewest possible turns. Each tunnel can be used only once per turn, and rooms (except the start and end) can hold only one ant at a time.

- The program is written in Go, with core logic located in the `pkg/` directory.
- The built-in visualizer runs as an HTTP server, with frontend assets and scripts located in the `visualizer/` directory.
- Sample files describing various ant farms are available in the `samples/` directory.

To build the project, ensure that you have `npm`, `make`, and Go installed, then simply run `make`.
