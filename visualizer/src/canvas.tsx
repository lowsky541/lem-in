import { useRef } from "react";

export function Canvas() {
  const canvasRef = useRef(null);
  return <canvas ref={canvasRef} />;
}
