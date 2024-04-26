async function FetchWASM(url) {
	let response = await fetch(url);
	let wasm = await WebAssembly.instantiateStreaming(response, go.importObject);

	go.run(wasm.instance);
	return new Promise((r) => r(true));
};

