importScripts("wasm_exec.js") // "$(go env GOROOT)/misc/wasm/wasm_exec.js"

const go = new self.Go()
WebAssembly.instantiateStreaming(fetch("convert.wasm"), go.importObject).then((result) => {
    go.run(result.instance)
})

self.addEventListener("message", async (event) => {
    try {
        const converted = await convertPortrait(event.data)
        self.postMessage({type: "RESULT", data: converted})
    } catch (err) {
        self.postMessage({type: "ERROR", data: err})
    }
})
