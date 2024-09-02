package lemin

import "errors"

var ErrInvalidFormat = errors.New("invalid lem-in format")
var ErrTunnelLoop = errors.New("invalid tunnel, was linked to itself")
var ErrRoomDuplication = errors.New("room duplication was found in input")
var ErrUnknownRoom = errors.New("unknown room when parsing tunnels")
var ErrCommandIsNotAllowed = errors.New("invalid command, only 'start' and 'end' are allowed")
var ErrInsaneGraph = errors.New("invalid graph, no connection from start to end room")
var ErrInvalidAntCount = errors.New("invalid ant count, either count is negative or is garbage")
