import { Bounds, ILeminContext, IRoom, LeminContext, transform, Vec2 } from "./data";
import "./main.scss";

const settings = {
  zoom: 1.0,
  nodeRadius: 5,
  lineThickness: 3,
  colors: {
    lines: 'grey',
    normal: 'blue',
    start: 'green',
    end: 'red',
  }
};

function getBounds(rooms: Array<IRoom>): Bounds {
  const [{ x: initialX, y: initialY }] = rooms;
  const initial = { left: initialX, right: initialX, top: initialY, bottom: initialY };

  return rooms.reduce((prev, curr) => {
    return {
      top: Math.min(prev.top, curr.y),
      right: Math.max(prev.right, curr.x),
      bottom: Math.max(prev.bottom, curr.y),
      left: Math.min(prev.left, curr.x),
    };
  }, initial);
}

function draw(canvas: HTMLCanvasElement, context: CanvasRenderingContext2D, { rooms, tunnels }: LeminContext) {
  canvas.width = window.innerWidth, canvas.height = window.innerHeight;

  const bounds = getBounds(rooms);
  const scalingFactors = {
    x: canvas.clientWidth / (bounds.right + bounds.left + settings.zoom),
    y: canvas.clientHeight / (bounds.bottom + bounds.top + settings.zoom),
  } as Vec2;

  tunnels.forEach(tunnel => {
    context.beginPath();
    context.moveTo((settings.zoom / 2 + tunnel.from.x) * scalingFactors.x, (settings.zoom / 2 + tunnel.from.y) * scalingFactors.y);
    context.lineTo((settings.zoom / 2 + tunnel.to.x) * scalingFactors.x, (settings.zoom / 2 + tunnel.to.y) * scalingFactors.y);
    context.lineWidth = settings.lineThickness;
    context.strokeStyle = settings.colors.lines;
    context.stroke();
  });

  rooms.forEach(room => {
    context.beginPath();
    context.arc((settings.zoom / 2 + room.x) * scalingFactors.x, (settings.zoom / 2 + room.y) * scalingFactors.y, settings.nodeRadius, 0, Math.PI * 2);
    context.fillStyle = room.isStart ? settings.colors.start : room.isEnd ? settings.colors.end : settings.colors.normal;
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
  .then((response): Promise<ILeminContext> => {
    if (!response.ok) throw new Error(response.statusText);
    else return response.json();
  })
  .then(transform)
  .then(main)
  .catch(alert);
