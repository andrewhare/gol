# gol

A Game of Life simulator. The only dependency is `github.com/pterm/pterm` which is used to draw the board.

### Instructions

To run the simulator:

```
$ go run main.go
```

By default, the simulator will create a 25x25 board, a glider, and create a new generation every second.

These values are configurable:

```terminal
$ go run main.go -h
Usage of /var/folders/hw/c45n61nx7sncm735hxv2rgtm0000gp/T/go-build3791553864/b001/exe/main:
  -dimensions int
    	board dimensions (default 25)
  -interval duration
    	the interval between generations (default 1s)
  -pattern string
    	initial pattern on the board (default "glider")
```

