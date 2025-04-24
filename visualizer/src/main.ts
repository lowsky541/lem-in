import { LeminContext } from "./data";
import "./main.scss";

function draw(canvas: HTMLCanvasElement, context: CanvasRenderingContext2D, { rooms, tunnels, turns }: LeminContext) {
  canvas.width = window.innerWidth, canvas.height = window.innerHeight;

  const firstRoom = rooms[0];
  const highest = rooms.reduce((prev, curr) => {
    return {
      x1: Math.min(prev.x1, curr.x), x2: Math.max(prev.x2, curr.x),
      y1: Math.min(prev.y1, curr.y), y2: Math.max(prev.y2, curr.y),
    };
  }, { x1: firstRoom.x, x2: firstRoom.x, y1: firstRoom.y, y2: firstRoom.y });
  const scalingFactor = { x: canvas.clientWidth / (highest.x2 + highest.x1), y: canvas.clientHeight / (highest.y2 + highest.y1) };

  tunnels.forEach(tunnel => {
    const from = rooms.find(r => r.id === tunnel.from)!;
    const to = rooms.find(r => r.id === tunnel.to)!;

    context.beginPath();
    context.moveTo(from.x * scalingFactor.x, from.y * scalingFactor.y);
    context.lineTo(to.x * scalingFactor.x, to.y * scalingFactor.y);
    context.strokeStyle = 'grey';
    context.stroke();
  });

  rooms.forEach(room => {
    context.beginPath();
    context.arc(room.x * scalingFactor.x, room.y * scalingFactor.y, 10, 0, Math.PI * 2);
    context.fillStyle = room.isStart ? "green" : room.isEnd ? "red" : "magenta",
      context.fill();
  });
}

function main(lemin: LeminContext) {
  const canvas = document.createElement('canvas');
  document.body.appendChild(canvas);

  const context = canvas?.getContext('2d');
  if (!context) {
    alert("Your browser doesn't seem to support canvas.");
    return;
  }

  draw(canvas, context, lemin);
}

window.onload = async () => fetch("/context")
  .then(r => r.json())
  .then(r => main(r))
  .catch(alert);
