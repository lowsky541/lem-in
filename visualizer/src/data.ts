export type RoomId = number;

export interface Bounds {
    top: number;
    right: number;
    bottom: number;
    left: number;
}

export interface IFromTo {
    to: IRoom;
    from: IRoom;
}

export interface IRoom {
    id: number;
    name: string;
    isStart: boolean;
    isEnd: boolean;
    x: number, y: number;
}

export interface ITunnel {
    id: number;
    distance: number;
    fromId: RoomId, toId: RoomId;
}

export type Tunnel = ITunnel & IFromTo;

export interface IMove {
    ant: number;
    fromId: RoomId, toId: RoomId;
}

export type Move = IMove & IFromTo;

export interface ITurn {
    moves: Array<IMove>;
}

export interface Turn {
    moves: Array<Move>;
}

export interface ILeminContext {
    ants: number;
    startId: RoomId;
    endId: RoomId;
    rooms: Array<IRoom>;
    tunnels: Array<ITunnel>;
    turns: Array<ITurn>;
}

export interface LeminContext {
    ants: number;
    startId: RoomId;
    start: IRoom;
    endId: RoomId;
    end: IRoom;
    rooms: Array<IRoom>;
    tunnels: Array<Tunnel>;
    turns: Array<Turn>;
}

function findRoom(rooms: IRoom[], id: RoomId): IRoom {
    return rooms.find(room => room.id === id)!;
}

function concreteTunnel(rooms: IRoom[]) {
    return (tunnel: ITunnel): Tunnel => ({
        ...tunnel,
        from: findRoom(rooms, tunnel.fromId),
        to: findRoom(rooms, tunnel.toId),
    });
}

function concreteMove(rooms: IRoom[]) {
    return (move: IMove): Move => ({
        ...move,
        from: findRoom(rooms, move.fromId),
        to: findRoom(rooms, move.toId),
    });
}

function concreteTurn(rooms: IRoom[]) {
    return (turn: ITurn): Turn => ({
        moves: turn.moves.map(concreteMove(rooms)),
    });
}

export function transform(context: ILeminContext): LeminContext {
    const { rooms, startId, endId } = context;

    return {
        ...context,
        start: findRoom(rooms, startId),
        end: findRoom(rooms, endId),
        tunnels: context.tunnels.map(concreteTunnel(rooms)),
        turns: context.turns.map(concreteTurn(rooms)),
    };
}
