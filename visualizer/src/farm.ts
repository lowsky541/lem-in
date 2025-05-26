export type RoomId = number;

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

function resolvedTunnel(rooms: Room[]) {
  return (tunnel: Tunnel): Tunnel => ({
    ...tunnel,
    from: findRoom(rooms, tunnel.fromId),
    to: findRoom(rooms, tunnel.toId),
  });
}

function resolvedMove(rooms: Room[]) {
  return (move: Move): Move => ({
    ...move,
    from: findRoom(rooms, move.fromId),
    to: findRoom(rooms, move.toId),
  });
}

function resolvedTurn(rooms: Room[]) {
  return (turn: Turn): Turn => turn.map(resolvedMove(rooms));
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

export function build(input: Farm): Farm {
  const { rooms, startId, endId } = input;

  return {
    ...input,
    start: findRoom(rooms, startId),
    end: findRoom(rooms, endId),
    tunnels: input.tunnels.map(resolvedTunnel(rooms)),
    turns: input.turns.map(resolvedTurn(rooms)),
    bounds: getBounds(rooms),
  };
}
