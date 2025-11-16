export type RoomId = number;
export type TunnelId = number;

export interface FarmBounds {
  top: number;
  right: number;
  bottom: number;
  left: number;
}

export interface Room {
  id: number;
  name: string;
  isStart: boolean;
  isEnd: boolean;
  x: number;
  y: number;
}

export interface Tunnel {
  id: number;
  distance: number;
  from: Room;
  fromId: RoomId;
  to: Room;
  toId: RoomId;
}

export interface Move {
  ant: number;
  from: Room;
  fromId: RoomId;
  to: Room;
  toId: RoomId;
  tunnel: Tunnel;
  tunnelId: number;
}

export type Turn = Array<Move>;

export interface Farm {
  ants: number;
  startId: RoomId;
  start: Room;
  endId: RoomId;
  end: Room;
  rooms: Array<Room>;
  tunnels: Array<Tunnel>;
  turns: Array<Turn>;
  bounds: FarmBounds;
}

function findRoom(rooms: Room[], id: RoomId): Room {
  return rooms.find(room => room.id === id)!;
}

function findTunnel(tunnels: Tunnel[], id: TunnelId): Tunnel {
  return tunnels.find(tunnel => tunnel.id === id)!;
}

function resolvedTunnel(rooms: Room[]) {
  return (tunnel: Tunnel): Tunnel => ({
    ...tunnel,
    from: findRoom(rooms, tunnel.fromId),
    to: findRoom(rooms, tunnel.toId),
  });
}

function resolvedMove(rooms: Room[], tunnels: Tunnel[]) {
  return (move: Move): Move => ({
    ...move,
    from: findRoom(rooms, move.fromId),
    to: findRoom(rooms, move.toId),
    tunnel: findTunnel(tunnels, move.tunnelId),
  });
}

function resolvedTurn(rooms: Room[], tunnels: Tunnel[]) {
  return (turn: Turn): Turn => turn.map(resolvedMove(rooms, tunnels));
}

function getBounds(rooms: Array<Room>): FarmBounds {
  const [{ x: initialX, y: initialY }] = rooms;
  const initial = {
    left: initialX,
    right: initialX,
    top: initialY,
    bottom: initialY,
  };

  return rooms.reduce((prev, curr) => {
    return {
      top: Math.min(prev.top, curr.y),
      right: Math.max(prev.right, curr.x),
      bottom: Math.max(prev.bottom, curr.y),
      left: Math.min(prev.left, curr.x),
    };
  }, initial);
}

export function build(farm: Farm): Farm {
  const { rooms, tunnels, startId, endId } = farm;

  console.log("build");

  return {
    ...farm,
    start: findRoom(rooms, startId),
    end: findRoom(rooms, endId),
    tunnels: farm.tunnels.map(resolvedTunnel(rooms)),
    turns: farm.turns.map(resolvedTurn(rooms, tunnels)),
    bounds: getBounds(rooms),
  };
}
