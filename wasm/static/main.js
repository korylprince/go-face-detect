window.WasmWorker.onmessage = event => {
	const {type, data} = event.data
	if (type === "RESULT") {
		const blob = new Blob([data], {type: "image/png"})
		document.getElementById("output").src = URL.createObjectURL(blob)
		document.getElementById("message").innerHTML = "Done!"
	} else if (type === "ERROR") {
		console.error(data)
		document.getElementById("message").innerHTML = data.toString()
	}
}

async function processDrop(items) {
	const files = []
	for (const item of items) {
		if (item.kind === "file" && (item.type === "image/jpeg" || item.type === "image/png")) {
			files.push(item)
		}
	}

	if (files.length < 1) {
		document.getElementById("message").innerHTML = "No valid images were dropped"
		return
	}

	const buffer = new Uint8Array(await files[0].getAsFile().arrayBuffer())

	document.getElementById("input").src = URL.createObjectURL(new Blob([buffer], {type: files[0].type}))
	document.getElementById("output").removeAttribute("src")

	document.getElementById("message").innerHTML = "Processing..."
	window.WasmWorker.postMessage(buffer)
}

document.body.addEventListener("dragover", (ev) => {
	ev.preventDefault()
})
document.body.addEventListener("drop", async (ev) => {
	ev.preventDefault()

	if (!ev.dataTransfer.items) {
		console.error("doesn't support dataTransfer.items")
		document.getElementById("message").innerHTML = "Browser doesn't support dataTransfer.items"
		return
	}

	processDrop(ev.dataTransfer.items)
})

