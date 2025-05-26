import { Vector2 } from "threejs-math";
import { Farm, FarmBounds } from "./farm";

export interface Settings {
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

export class Lemin {
  settings: Settings;
  farm: Farm;
  farmBounds: FarmBounds;
  canvas: HTMLCanvasElement;
  context: CanvasRenderingContext2D;
  turnIndex: number;
  factors: Vector2;

  constructor(
    settings: Settings,
    farm: Farm,
    canvas: HTMLCanvasElement,
    context: CanvasRenderingContext2D
  ) {
    this.settings = settings;
    this.farm = farm;
    this.canvas = canvas;
    this.context = context;
    this.turnIndex = 0;
    this.farmBounds = this.farm.bounds;

    this.resize();
    this.factors = this.updateFactors();
    this.draw(0);

    document.addEventListener("keydown", this.keydown.bind(this));
    window.addEventListener("resize", this.resize.bind(this));
  }

  private scaledBy(value: number, factor: "x" | "y") {
    return (
      (this.settings.zoom / 2 + value) *
      (factor == "x" ? this.factors.x : this.factors.y)
    );
  }

  private updateFactors() {
    this.factors = new Vector2(
      this.canvas.clientWidth /
        (this.farmBounds.right + this.farmBounds.left + this.settings.zoom),
      this.canvas.clientHeight /
        (this.farmBounds.bottom + this.farmBounds.top + this.settings.zoom)
    );
    return this.factors;
  }

  private resize() {
    this.canvas.width = window.innerWidth;
    this.canvas.height = window.innerHeight;
    this.updateFactors();
  }

  animate() {
    requestAnimationFrame(this.draw.bind(this));
  }

  private draw(time: DOMHighResTimeStamp) {
    requestAnimationFrame(this.draw.bind(this));

    const { rooms, tunnels } = this.farm;
    const { context: ctx } = this;

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
      ctx.lineWidth = this.settings.tunnelThickness;
      ctx.strokeStyle = this.settings.tunnelColor;
      ctx.stroke();
    });

    rooms.forEach(room => {
      ctx.beginPath();
      ctx.arc(
        this.scaledBy(room.x, "x"),
        this.scaledBy(room.y, "y"),
        this.settings.roomRadius,
        0,
        Math.PI * 2
      );
      ctx.fillStyle = room.isStart
        ? this.settings.roomStartColor
        : room.isEnd
        ? this.settings.roomEndColor
        : this.settings.roomColor;
      ctx.fill();
    });
  }

  keydown(ev: KeyboardEvent) {
    const { key } = ev;
    if (key === "KeyLeft") {
      this.turnIndex--;
      this.animate();
    } else if (key === "KeyRight") {
      this.turnIndex++;
      this.animate();
    }
  }
}
