package lemin

type Farm struct {
	Ants    int
	Start   *Room
	End     *Room
	Rooms   []*Room
	Tunnels []*Tunnel
}

type FarmResponse struct {
	Ants    int       `json:"ants"`
	Start   int       `json:"startId"`
	End     int       `json:"endId"`
	Rooms   []*Room   `json:"rooms"`
	Tunnels []*Tunnel `json:"tunnels"`
	Turns   []Turn    `json:"turns"`
}

type Turn = []Move

type Move struct {
	Ant    int   `json:"ant"`
	FromId uint  `json:"fromId"`
	From   *Room `json:"-"`
	ToId   uint  `json:"toId"`
	To     *Room `json:"-"`
}
