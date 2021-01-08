
.ONESHELL:
buildBasicExample:
	cd examples/basic
	GOOS=js GOARCH=wasm go build -o ../../web/main.wasm
	cp index.html ../../web

.ONESHELL:
buildImageExample:
	cd examples/gifs
	GOOS=js GOARCH=wasm go build -o ../../web/main.wasm
	cp index.html ../../web