import { Vector2 } from "threejs-math";
import { Farm, FarmBounds, Turn } from "./farm";

export interface Settings {
  speed: number;
  zoom: number;
  antRadius: number;
  antColor: string;
  roomRadius: number;
  roomColor: string;
  roomStartColor: string;
  roomEndColor: string;
  tunnelThickness: number;
  tunnelColor: string;
}

interface Interant {
  a: Vector2;
  b: Vector2;
  t: number;
}

export class Lemin {
  settings: Settings;
  farm: Farm;
  canvas: HTMLCanvasElement;
  context: CanvasRenderingContext2D;

  private turnIndex: number;
  private animTurn: Turn;
  private scalingFactors: Vector2;
  private start: number;
  private isAnimating: boolean;

  constructor(
    farm: Farm,
    settings: Settings,
    canvas: HTMLCanvasElement,
    context: CanvasRenderingContext2D
  ) {
    settings.zoom = settings.zoom || 1.0;
    settings.speed = settings.speed || 1.0;

    this.farm = farm;
    this.settings = settings;
    this.canvas = canvas;
    this.context = context;
    this.turnIndex = 0;
    this.isAnimating = false;
    this.draw();

    document.addEventListener("keydown", this.keydown.bind(this));
  }

  private scaledBy(value: number, factor: "x" | "y") {
    return (
      (this.settings.zoom / 2 + value) *
      (factor == "x" ? this.scalingFactors.x : this.scalingFactors.y)
    );
  }

  private animate() {
    this.isAnimating = true;
    this.start = new Date().getTime();

    // Precalculations
    const turn = this.farm.turns[this.turnIndex];
    const longestDistance = turn.reduce((prev, current) => {
      return prev.tunnel.distance > current.tunnel.distance ? prev : current;
    }).tunnel.distance;

    const ants: Array<Interant> = turn.map(({ ant, from, to }) => ({
      a: new Vector2(from.x, from.y),
      b: new Vector2(to.x, to.y),
      t: 0.0,
    }));
  }

  private requestFrame() {
    return requestAnimationFrame(this.draw.bind(this));
  }

  private resize() {
    this.canvas.width = window.innerWidth;
    this.canvas.height = window.innerHeight;
    this.scalingFactors = new Vector2(
      this.canvas.clientWidth /
        (this.farm.bounds.right + this.farm.bounds.left + this.settings.zoom),
      this.canvas.clientHeight /
        (this.farm.bounds.bottom + this.farm.bounds.top + this.settings.zoom)
    );
  }

  private drawFarm() {
    const { rooms, tunnels } = this.farm;
    const { context: ctx, settings: sts } = this;

    tunnels.forEach(tunnel => {
      ctx.beginPath();
      ctx.moveTo(
        this.scaledBy(tunnel.from.x, "x"),
        this.scaledBy(tunnel.from.y, "y")
      );
      ctx.lineTo(
        this.scaledBy(tunnel.to.x, "x"),
        this.scaledBy(tunnel.to.y, "y")
      );
      ctx.lineWidth = sts.tunnelThickness;
      ctx.strokeStyle = sts.tunnelColor;
      ctx.stroke();
    });

    rooms.forEach(room => {
      ctx.beginPath();
      ctx.arc(
        this.scaledBy(room.x, "x"),
        this.scaledBy(room.y, "y"),
        sts.roomRadius,
        0,
        Math.PI * 2
      );
      ctx.fillStyle = room.isStart
        ? sts.roomStartColor
        : room.isEnd
        ? sts.roomEndColor
        : sts.roomColor;
      ctx.fill();
    });
  }

  private drawAnts() {
    const turn = this.farm.turns[this.turnIndex];
    const { context: ctx, settings: sts } = this;

    turn.forEach(move => {
      ctx.beginPath();
      ctx.arc(
        this.scaledBy(move.from.x, "x"),
        this.scaledBy(move.from.y, "y"),
        sts.antRadius,
        0,
        Math.PI * 2
      );
      ctx.fillStyle = sts.antColor;
      ctx.fill();
    });
  }

  private draw(time?: DOMHighResTimeStamp) {
    this.requestFrame();
    this.resize();
    this.drawFarm();
    this.drawAnts();
  }

  private keydown(ev: KeyboardEvent) {
    const { key } = ev;
    if (key === "ArrowLeft" && this.turnIndex - 1 >= 0) {
      this.turnIndex--;
    } else if (
      key === "ArrowRight" &&
      this.turnIndex + 1 < this.farm.turns.length
    ) {
      this.turnIndex++;
    } else if (key == "Space") {
      this.animate();
    }
  }
}
