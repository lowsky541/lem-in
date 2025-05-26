import { Farm, build } from "./farm";
import { Lemin, Settings } from "./lemin";

const settings = {
  zoom: 1.0,
  antRadius: 8,
  antColor: "blue",
  roomRadius: 12,
  roomColor: "magenta",
  roomStartColor: "green",
  roomEndColor: "red",
  tunnelThickness: 2,
  tunnelColor: "grey",
} as Settings;

function main(farm: Farm) {
  const canvas = document.querySelector<HTMLCanvasElement>("#canvas")!;
  const context = canvas.getContext("2d")!;
  new Lemin(settings, farm, canvas, context);
}

addEventListener(
  "load",
  async () =>
    await fetch("/farm")
      .then(response => {
        if (!response.ok) throw new Error(response.statusText);
        else return response.json();
      })
      .then(build)
      .then(main)
      .catch(alert),
  { once: true }
);
