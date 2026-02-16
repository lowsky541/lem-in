import React from "react";
import ReactDOM from "react-dom/client";
import Visualizer from './visualizer';

const rootElement = document.getElementById("root");

if (!rootElement) {
  throw new Error(
    "Root element not found."
  );
}

ReactDOM.createRoot(rootElement).render(
  <React.StrictMode>
    <Visualizer />
  </React.StrictMode>
);
