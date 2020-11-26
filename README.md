## CSA Coursework: Game of Life

### Go, one of the reasons that I can't sleep well until Christmas.

```txt
_________                                                 
__  ____/______ ______ ______ ______ ______ ______ ______ 
_  / __  _  __ \_  __ \_  __ \_  __ \_  __ \_  __ \_  __ \
/ /_/ /  / /_/ // /_/ // /_/ // /_/ // /_/ // /_/ // /_/ /
\\____/   \____/ \____/ \____/ \____/ \____/ \____/ \____/ 

```

### outline

#### main
1. main.go
    call functions to run the whole project

#### gol
1. distributor.go
    - struct
        - `distributorChannels`
    - func
        - `calculateNeighbors (p Params, x, y int, world [][]byte) int`
        - `calculateNextStage (p Params, world [][]byte) [][]byte`
        - `calculateAliveCells (p Params, world [][]byte) []util.Cell`
        - `distributor(p Params, c distributorChannels)`

2. gol.go
    - struct
        - `Params`
    - func
        - `Run(p Params, events chan<- Event, keyPresses <-chan rune)`

3. io.go
    - struct
        - `ioChannels`
        - `ioState`
        - `ioCommand`
    - func
        - `(io *ioState) writePgmImage()`
        - `(io *ioState) readPgmImage()`
        - `startIo(p Params, c ioChannels)`


#### util
1. cell.go
    - struct
        - `Cell`


### tests and todo related

#### stage 1
1. 11/23
    - step 1 & step 2 done
    - step 3
        - the `alivecells` event can not be received

### notes

#### 11/24
1. the `ioInput` and `ioOutput` is linked together to share the world, in this way, we don't need an extra `tempworld` channel to pass the changes in world in each round.
2. we need to input and output the world in each round, either before the calculation or after the modification.

#### 11/26
1. stage one done with single threads








