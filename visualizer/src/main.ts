import { Bounds, LeminFarm as Farm, Room, transform } from "./farm";
import "./main.scss";

interface AnimationState {
  playing: boolean;
  antId: number;
}

interface Visuals {
  tunnels: {
    thickness: number;
    color: string;
  };
  rooms: {
    radius: number;
    color: string;
    startColor: string;
    endColor: string;
  };
}

interface State {
  farm: Farm;
  canvas: HTMLCanvasElement;
  context: CanvasRenderingContext2D;
  zoom: number;
  bounds: Bounds;
  animation: AnimationState;
  visuals: Visuals;
}

function getBounds(rooms: Array<Room>): Bounds {
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

function run(
  time: DOMHighResTimeStamp,
  farm: Farm,
  cvs: HTMLCanvasElement,
  ctx: CanvasRenderingContext2D,
  bounds: Bounds
) {
  const { rooms, tunnels } = farm;

  // Resize the canvas
  cvs.width = window.innerWidth;
  cvs.height = window.innerHeight;

  // Get the scaling factor
  const factors = {
    x: cvs.clientWidth / (bounds.right + bounds.left + state.zoom),
    y: cvs.clientHeight / (bounds.bottom + bounds.top + state.zoom),
  };

  // Scale helper
  const scale = (v: number, factor: number) => (state.zoom / 2 + v) * factor;

  tunnels.forEach(tunnel => {
    ctx.beginPath();
    ctx.moveTo(
      scale(tunnel.from.x, factors.x),
      scale(tunnel.from.y, factors.y)
    );
    ctx.lineTo(scale(tunnel.to.x, factors.x), scale(tunnel.to.y, factors.y));
    ctx.lineWidth = lineThickness;
    ctx.strokeStyle = colors.lines;
    ctx.stroke();
  });

  rooms.forEach(room => {
    ctx.beginPath();
    ctx.arc(
      scale(room.x, factors.x),
      scale(room.y, factors.y),
      nodeRadius,
      0,
      Math.PI * 2
    );
    ctx.fillStyle = roomColor(room);
    ctx.fill();
    ctx.save();
  });

  window.requestAnimationFrame(time => run(time, farm, cvs, ctx, bounds));
}

function install(farm: Farm) {
  const canvas = document.querySelector<HTMLCanvasElement>("#canvas")!;
  const context = canvas.getContext("2d")!;
  const bounds = getBounds(farm.rooms);
  window.requestAnimationFrame(time => {
    run(time, farm, canvas, context, bounds);
  });
}

window.onload = async () =>
  fetch("/farm")
    .then(response => {
      if (!response.ok) throw new Error(response.statusText);
      else return response.json();
    })
    .then(transform)
    .then(install);
