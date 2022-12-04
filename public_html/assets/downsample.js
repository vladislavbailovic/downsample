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

const getPalette = () => {
	return [...document.querySelectorAll(".palette input")]
		.map(x => parseInt(x.value.substr(1), 16));
}

const render = algo => {
	img = getSource();
	if ("average" == algo) {
		palette = getPalette()
		raw = average(img.data, img.width, img.height, palette)
	} else if ("normalize" == algo) {
		raw = normalize(img.data, img.width, img.height)
	} else {
		raw = pixelate(img.data, img.width, img.height)
	}
	renderData(new ImageData(raw, img.width, img.height));
	updateInterface(algo)
}

const updateInterface = algo => {
	const tile = document.querySelector(".tile-size input");
	tile.value = getTileSize();

	const palette = document.querySelector(".palette");
	if (algo != "average") {
		palette.style.display = "none";
		return;
	}

	palette.style.display = "flex";
	[...palette.querySelectorAll(".color")].map(c => {
		const x = c.querySelector("input");
		if (!x) return;

		c.style.backgroundColor = x.value;
	});
}

const rerender = () => {
	const algo = document.getElementById("algo")
	render(algo.value)
}

const init = () => {
	const algo = document.getElementById("algo")
	algo.addEventListener("change", e => {
		render(algo.value)
	});

	const add = document.querySelector(".palette .add");
	add.addEventListener("click", e => {
		const palette = document.querySelector(".palette");
		const clr = palette.querySelector(".color")
			.cloneNode(true);
		add.before(clr);
		render(algo.value);
	});

	const tile = document.querySelector(".tile-size input");
	tile.addEventListener("change", e => {
		const tileSize = parseInt(tile.value, 10);
		if (!tileSize) {
			return;
		}
		setTileSize(tileSize)
		rerender()
	});

	document.addEventListener("change", e => {
		if (e.target.nodeName == "INPUT" && e.target.closest(".color")) {
			return rerender();
		}
	});

	document.addEventListener("click", e => {
		if (e.target.nodeName == "BUTTON" && e.target.closest(".color")) {
			e.target.closest(".color").remove();
			rerender();
		}
	});

	render();
}

window.addEventListener("load", init);
