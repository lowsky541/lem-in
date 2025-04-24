export type RoomId = number;

export interface Room {
    id: number;
    name: string;
    x: number;
    y: number;
    isStart: boolean;
    isEnd: boolean;
}

export interface Tunnel {
    id: number;
    distance: number;
    from: RoomId;
    to: RoomId;
}

export interface Move {
    ant: number;
    from: RoomId;
    to: RoomId;
}

export interface Turn {
    moves: Array<Move>;
}

export interface LeminContext {
    ants: number;
    start: RoomId;
    end: RoomId;
    rooms: Array<Room>;
    tunnels: Array<Tunnel>;
    turns: Array<Turn>;
}
