### CSA Coursework: Game of Life

#### Go, one of the reasons that I can't sleep well until Christmas.

```txt
_________                                                 
__  ____/______ ______ ______ ______ ______ ______ ______ 
_  / __  _  __ \_  __ \_  __ \_  __ \_  __ \_  __ \_  __ \
/ /_/ /  / /_/ // /_/ // /_/ // /_/ // /_/ // /_/ // /_/ /
\\____/   \____/ \____/ \____/ \____/ \____/ \____/ \____/ 

```

#### outline

##### main
1. main.go
    call functions to run the whole project

##### gol
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


##### util
1. cell.go
    - struct
        - `Cell`


#### tests and todo related

##### 11/20







