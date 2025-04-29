export type RoomId = number;

export interface Bounds {
  top: number;
  right: number;
  bottom: number;
  left: number;
}

export interface IFromTo {
  to: Room;
  from: Room;
}

export interface Room {
  id: number;
  name: string;
  isStart: boolean;
  isEnd: boolean;
  x: number;
  y: number;
}

export interface ITunnel {
  id: number;
  distance: number;
  fromId: RoomId;
  toId: RoomId;
}

export type Tunnel = ITunnel & IFromTo;

export interface IMove {
  ant: number;
  fromId: RoomId;
  toId: RoomId;
}

export type Move = IMove & IFromTo;
export type Turn = Array<Move>;

export interface ILeminFarm {
  ants: number;
  startId: RoomId;
  endId: RoomId;
  rooms: Array<Room>;
  tunnels: Array<ITunnel>;
  turns: Array<Turn>;
}

export interface LeminFarm {
  ants: number;
  startId: RoomId;
  start: Room;
  endId: RoomId;
  end: Room;
  rooms: Array<Room>;
  tunnels: Array<Tunnel>;
  turns: Array<Turn>;
}

function findRoom(rooms: Room[], id: RoomId): Room {
  return rooms.find(room => room.id === id)!;
}

function concreteTunnel(rooms: Room[]) {
  return (tunnel: ITunnel): Tunnel => ({
    ...tunnel,
    from: findRoom(rooms, tunnel.fromId),
    to: findRoom(rooms, tunnel.toId),
  });
}

function concreteMove(rooms: Room[]) {
  return (move: IMove): Move => ({
    ...move,
    from: findRoom(rooms, move.fromId),
    to: findRoom(rooms, move.toId),
  });
}

function concreteTurn(rooms: Room[]) {
  return (turn: Turn): Turn => turn.map(concreteMove(rooms));
}

export function transform(context: ILeminFarm): LeminFarm {
  const { rooms, startId, endId } = context;

  return {
    ...context,
    start: findRoom(rooms, startId),
    end: findRoom(rooms, endId),
    tunnels: context.tunnels.map(concreteTunnel(rooms)),
    turns: context.turns.map(concreteTurn(rooms)),
  };
}
