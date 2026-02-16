import "./wasm_exec.js";

export async function initLemin() {
  const go = new Go();

  const result = await WebAssembly.instantiateStreaming(
    fetch(new URL("./lem-in.wasm", import.meta.url)),
    go.importObject
  );

  go.run(result.instance);

  return globalThis.Lemin;
}