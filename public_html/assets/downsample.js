const renderData = img => {
	const output = document.getElementById("output")
	output.width = img.width
	output.height = img.height

	const otx = output.getContext("2d")
	otx.putImageData(img, 0, 0, 0, 0, img.width, img.height);
}

const getSource = () => {
	const data = document.createElement('canvas');
	const ctx = data.getContext("2d");
	const input = document.getElementById("input-file");
	data.width = input.width
	data.height = input.height

	ctx.drawImage(input, 0, 0);
	return ctx.getImageData(0, 0, input.width, input.height);
}

const render = algo => {
	img = getSource();
	if ("average" == algo) {
		raw = average(img.data, img.width, img.height)
	} else if ("normalize" == algo) {
		raw = normalize(img.data, img.width, img.height)
	} else {
		raw = pixelate(img.data, img.width, img.height)
	}
	renderData(new ImageData(raw, img.width, img.height));
}

const init = () => {
	const algo = document.getElementById("algo")
	algo.addEventListener("change", e => {
		render(algo.value)
	});
	render()
}

window.addEventListener("load", init);
