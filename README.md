# Disco Gopher

This is a fork of [golang-wasm-example](https://github.com/mattn/golang-wasm-example), with these added features:

* The gopher has a mustache.
* The gopher is moving randomly left and right.
* The background is blinking.

## Build


```sh
GOOS=js GOARCH=wasm go generate
GOOS=js GOARCH=wasm go build -o main.wasm main.go
```


## Build, run, serve and open in a browser

```sh
make
```

## License

MIT

## Original Author

Yasuhrio Matsumoto (a.k.a. mattn)
