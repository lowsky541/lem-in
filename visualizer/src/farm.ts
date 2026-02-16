interface Bounds {
  left: number;
  right: number;
  top: number;
  bottom: number;
}

interface Room {
  id: number;
  name: string;
  x: number;
  y: number;
  isStart: boolean;
  isEnd: boolean;
}

interface Tunnel {
  id: number;
  distance: number;
}

type Rooms = { [key: string]: Room };
type Tunnels = { [key: string]: Tunnel };

interface Response {
  ants: number;
  start: string;
  end: string;
  rooms: Rooms;
  tunnels: Tunnels;
}

export class Farm {
  private _ants: number;
  private _rooms: Rooms;
  private _start: Room;
  private _end: Room;
  private _tunnels: Tunnels;

  constructor(response: Response) {
    this._ants = response.ants;
    this._rooms = response.rooms;
    this._start = response.rooms[response.start];
    this._end = response.rooms[response.end];
    this._tunnels = response.tunnels;
  }

  public get bounds(): Bounds {
    // Get x and y coordinates of the first room
    const [{ x: firstRoomX, y: firstRoomY }] = this.rooms;

    let bounds = {
      left: firstRoomX,
      right: firstRoomX,
      top: firstRoomY,
      bottom: firstRoomY,
    };

    this.rooms.forEach(room => {
      bounds = {
        top: Math.min(bounds.top, room.y),
        right: Math.max(bounds.right, room.x),
        bottom: Math.max(bounds.bottom, room.y),
        left: Math.min(bounds.left, room.x),
      };
    });

    return bounds;
  }

  public get ants(): number {
    return this._ants;
  }

  public get start(): Room {
    return this._start;
  }

  public get end(): Room {
    return this._end;
  }

  public get rooms(): Array<Room> {
    return Object.values(this._rooms);
  }

  public get tunnels(): Array<Tunnel> {
    return Object.values(this._tunnels);
  }
}

function getFarmData() {}
