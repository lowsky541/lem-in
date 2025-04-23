package lemin

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Describe all the structures of a parsed input file.
type Context struct {
	Ants    int
	Start   *Room
	End     *Room
	Rooms   []*Room
	Tunnels []*Tunnel
}

///////////////////////////////////////////////////////////////////////////////////////
//                             Regular Expressions
///////////////////////////////////////////////////////////////////////////////////////

var RegRoom = regexp.MustCompile(`^([^L]\w*)\s+(\d+)\s+(\d+)$`)
var RegTunnel = regexp.MustCompile(`^(\w+)\-(\w+)$`)

///////////////////////////////////////////////////////////////////////////////////////
//                                 File parser
///////////////////////////////////////////////////////////////////////////////////////

// Open the file pointed by `filepath` scan and fill the context `p`.
// You can think of the context as the graph.
func (p *Context) Parse(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Parse the ant count
	if scanner.Scan() {
		antsString := scanner.Text()
		ants, err := strconv.Atoi(antsString)
		if err != nil || !(ants > 0) {
			return ErrInvalidAntCount
		}
		p.Ants = ants
	} else {
		return ErrInvalidFormat
	}

	var curRoomId uint = 0
	var curTunnelId uint = 0

	// This is now the link section, no more rooms
	pastRooms := false

	// Next room is the start
	nextIsStart := false

	// Next room is the end
	nextIsEnd := false

	// Parser outputs
	var rooms = map[string]*Room{}
	var tunnels []*Tunnel

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
				return ErrCommandIsNotAllowed
			}
		} else if strings.HasPrefix(line, "#") || line == "" {
			// Allow blank lines without spaces and comments
			continue
		}

		// Is the current line defining a room ?
		roomMatch := RegRoom.FindStringSubmatch(line)
		if roomMatch != nil && !pastRooms {

			x, err := strconv.Atoi(roomMatch[2])
			if err != nil {
				return ErrInvalidFormat
			}

			y, err := strconv.Atoi(roomMatch[3])
			if err != nil {
				return ErrInvalidFormat
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
				return ErrRoomDuplication
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
		tunnelMatch := RegTunnel.FindStringSubmatch(line)
		if tunnelMatch != nil {
			p1 := tunnelMatch[1]
			p2 := tunnelMatch[2]

			room1, room1Exists := rooms[p1]
			room2, room2Exists := rooms[p2]
			if !room1Exists || !room2Exists {
				return ErrUnknownRoom
			}

			distance := Distance(room1, room2)
			tunnel := &Tunnel{
				Id:       curTunnelId,
				From:     room1,
				To:       room2,
				FromId:   room1.Id,
				ToId:     room2.Id,
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
		return ErrInvalidFormat
	}

	// There was an error reading the file
	if err := scanner.Err(); err != nil {
		return err
	}

	// A ##start or a ##end hasn't a associated room
	if nextIsStart || nextIsEnd {
		return ErrInvalidFormat
	}

	p.Rooms = MapValues(rooms)
	p.Tunnels = tunnels

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////
//                                   Banner
///////////////////////////////////////////////////////////////////////////////////////

func (p *Context) PrintBanner() {
	fmt.Println(p.Ants)
	for _, r := range p.Rooms {
		fmt.Printf("%s %v %v\n", r.Name, r.X, r.Y)
	}

	for _, t := range p.Tunnels {
		fmt.Printf("%s-%s\n", t.From.Name, t.To.Name)
	}

	fmt.Println()
}

///////////////////////////////////////////////////////////////////////////////////////
//                                   Utilities
///////////////////////////////////////////////////////////////////////////////////////

// The distance between two rooms
func Distance(r1 *Room, r2 *Room) float64 {
	var x = math.Pow(float64(r1.X-r2.X), 2)
	var y = math.Pow(float64(r1.Y-r2.Y), 2)
	return math.Sqrt(x + y)
}
