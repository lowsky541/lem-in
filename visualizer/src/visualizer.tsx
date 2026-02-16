import { useEffect, useState } from "react";
import { Farm } from "./farm";
import { Layer, Rect, Stage, Text } from "react-konva";


export default function Visualizer() {
  const [farm, setFarm] = useState<Farm | null>();

  useEffect(() => {
    fetch("/farm")
      .then(res => res.json())
      .then(data => new Farm(data))
      .then(instance => setFarm(instance));
  }, []);

  // Get bounds
  const bounds = farm?.bounds;
  const rooms = farm?.rooms.map(room => {
    const fill = room.isStart ? "blue" : "red";
    return (
      <Rect fill={fill} x={room.x} y={room.y} width={20} height={20}></Rect>
    );
  });

  return (
    <Stage width={window.innerWidth} height={window.innerHeight}>
      <Layer>{rooms}</Layer>
    </Stage>
  );
}
