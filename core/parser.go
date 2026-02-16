package core

import (
	"bufio"
	"errors"
	"io"
	"lemin/util"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var ErrInvalidFormat = errors.New("invalid lem-in format")
var ErrTunnelLoop = errors.New("invalid tunnel, was linked to itself")
var ErrRoomDuplication = errors.New("room duplication was found in input")
var ErrUnknownRoom = errors.New("unknown room when parsing tunnels")
var ErrCommandIsNotAllowed = errors.New("invalid command, only 'start' and 'end' are allowed")
var ErrInvalidAntCount = errors.New("invalid ant count, count must be a positive non-zero integer")
var ErrInsaneFarm = errors.New("no connection from start to end room")

var regRoom = regexp.MustCompile(`^([^L]\w*)\s+(\d+)\s+(\d+)$`)
var regTunnel = regexp.MustCompile(`^(\w+)\-(\w+)$`)

func Parse(s string) (*Farm, error) {
	reader := strings.NewReader(s)
	return ParseFromReader(reader)
}

func ParseFromReader(reader io.Reader) (*Farm, error) {
	var p *Farm = &Farm{}

	scanner := bufio.NewScanner(reader)

	// Parse the ant count
	if scanner.Scan() {
		antsString := scanner.Text()
		ants, err := strconv.Atoi(antsString)
		if err != nil || !(ants > 0) {
			return nil, ErrInvalidAntCount
		}
		p.Ants = ants
	} else {
		return nil, ErrInvalidFormat
	}

	var curRoomId int = 0
	var curTunnelId int = 0

	// This is now the link section, no more rooms
	pastRooms := false

	// Next room is the start
	nextIsStart := false

	// Next room is the end
	nextIsEnd := false

	// Parser outputs
	var rooms = map[string]*Room{}
	var tunnels []*Tunnel = []*Tunnel{}

	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "##") {
			if line == "##start" {
				nextIsStart = true
				continue
			} else if line == "##end" {
				nextIsEnd = true
				continue
			} else {
				return nil, ErrCommandIsNotAllowed
			}
		} else if strings.HasPrefix(line, "#") || line == "" {
			// Allow blank lines without spaces and comments
			continue
		}

		// Is the current line defining a room ?
		roomMatch := regRoom.FindStringSubmatch(line)
		if roomMatch != nil && !pastRooms {
			x, err := strconv.Atoi(roomMatch[2])
			if err != nil {
				return nil, ErrInvalidFormat
			}

			y, err := strconv.Atoi(roomMatch[3])
			if err != nil {
				return nil, ErrInvalidFormat
			}

			roomName := roomMatch[1]
			room := &Room{
				Id:      curRoomId,
				Name:    roomName,
				X:       x,
				Y:       y,
				IsStart: nextIsStart,
				IsEnd:   nextIsEnd,
				Tunnels: nil,
			}

			// Check for room duplication
			if _, exist := rooms[roomName]; exist {
				return nil, ErrRoomDuplication
			}

			rooms[roomName] = room

			if nextIsStart {
				p.Start = room
			} else if nextIsEnd {
				p.End = room
			}

			nextIsStart = false
			nextIsEnd = false
			curRoomId++

			continue
		}

		// Is the current line defining a link ?
		tunnelMatch := regTunnel.FindStringSubmatch(line)
		if tunnelMatch != nil {
			p1 := tunnelMatch[1]
			p2 := tunnelMatch[2]

			room1, room1Exists := rooms[p1]
			room2, room2Exists := rooms[p2]
			if !room1Exists || !room2Exists {
				return nil, ErrUnknownRoom
			}

			distance := distance(room1, room2)
			tunnel := &Tunnel{
				Id:       curTunnelId,
				From:     room1,
				To:       room2,
				Distance: distance,
			}
			curTunnelId++

			tunnels = append(tunnels, tunnel)
			room1.Tunnels = append(room1.Tunnels, tunnel)
			room2.Tunnels = append(room2.Tunnels, tunnel)

			pastRooms = true
			continue
		}

		// This is neither a room nor a link
		// What the hell is this trash ? Don't wanna know.
		return nil, ErrInvalidFormat
	}

	// There was an error reading the file
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Only the ant count is set which doesn't make sense
	if len(rooms) == 0 || len(tunnels) == 0 {
		return nil, ErrInvalidFormat
	}

	// A ##start or a ##end hasn't a associated room
	if nextIsStart || nextIsEnd {
		return nil, ErrInvalidFormat
	}

	p.Rooms = util.Values(rooms)
	p.Tunnels = tunnels

	return p, nil
}

func distance(r1 *Room, r2 *Room) float64 {
	var x = math.Pow(float64(r1.X-r2.X), 2)
	var y = math.Pow(float64(r1.Y-r2.Y), 2)
	return math.Sqrt(x + y)
}
