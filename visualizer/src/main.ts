import Graph from "graphology";
import Sigma from 'sigma';
import "./main.scss";

const DEFAULT_VERTEX_SIZE: number = 15;
const DEFAULT_EDGE_SIZE: number = 5;

type RoomId = number;

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
    from: RoomId;
    to: RoomId;
}

interface Move {
    ant: number;
    from: RoomId;
    to: RoomId;
}

interface Turn {
    moves: Array<Move>;
}

interface Data {
    ants: number;
    start: RoomId;
    end: RoomId;
    rooms: Array<Room>;
    tunnels: Array<Tunnel>;
    turns: Array<Turn>;
}

window.onload = async () => {
    let currentTurn: number = 0;

    const graph = new Graph({});
    const sigmaContainer = document.querySelector("#sigma")! as HTMLDivElement;

    const data: Data = await fetch("/data", {
        method: "get",
    }).then(res => res.json(), console.error);

    const rooms = data.rooms;
    const tunnels = data.tunnels;
    const turns = data.turns;

    rooms.forEach(room => {
        graph.addNode(room.id, {
            x: room.x,
            y: room.y,
            size: DEFAULT_VERTEX_SIZE,
            label: room.name,
            color: room.isStart ? "green" : room.isEnd ? "red" : "gray",
        });
    });

    tunnels.forEach(tunnel => {
        graph.addEdge(tunnel.from, tunnel.to, {
            size: DEFAULT_EDGE_SIZE,
            color: "magenta",
            label: tunnel.distance.toFixed(2),
        });
    });

    const sigma = new Sigma(graph, sigmaContainer, {
        renderEdgeLabels: true,
        renderLabels: true,
    });

    const camera = sigma.getCamera();
    camera.setState({ ratio: 1.5, x: camera.x, y: camera.y });

    ///////////////////////////////////////////////////////////

    const currentTurnIndicator = document.querySelector(
        "#controller-current-turn",
    ) as HTMLSpanElement;

    const nextTurnButton = document.querySelector(
        "#controller-next-turn",
    ) as HTMLAnchorElement;

    const history = document.querySelector("#history") as HTMLOListElement;

    const colorStack = [
        "yellow",
        "aqua",
        "fuchsia",
        "goldenrod",
        "greenyellow",
        "blue",
    ];
    let colorStackIndex = 0;

    const colors = new Map<number, string>();

    nextTurnButton.addEventListener("click", () => {
        if (currentTurn + 1 === turns.length)
            currentTurnIndicator.innerText =
                (currentTurn + 1).toString() + " (end)";
        else if (currentTurn + 1 < turns.length)
            currentTurnIndicator.innerText = (currentTurn + 1).toString();
        else return;

        const turn = turns[currentTurn++];
        let turnHistoryList = document.createElement("ul");

        const turnColor = colorStack[colorStackIndex++ % colorStack.length];
        turn.moves.forEach(m => {
            let color = colors.get(m.ant);
            if (color === undefined) {
                colors.set(m.ant, turnColor);
                color = turnColor;
            }

            if (m.to !== data.end) {
                graph.setNodeAttribute(m.to, "color", color);
            }

            if (m.from !== data.start) {
                graph.setNodeAttribute(m.from, "color", "gray");
            }

            const turnHistoryItem = document.createElement("li");
            turnHistoryItem.innerText = `${m.ant} moved from ${m.from} to ${
                m.to
            }${m.to === data.end ? " (out)" : ""}`;
            turnHistoryList.appendChild(turnHistoryItem);
        });

        const turnHistoryWrapper = document.createElement("li");
        turnHistoryWrapper.style.color = turnColor;
        turnHistoryWrapper.appendChild(turnHistoryList);
        history.append(turnHistoryWrapper);
    });
};
