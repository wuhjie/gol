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

#### parallel
```txt
gol 
    | distributor.go
    | event.go
    | gameLogic.go
    | gol.go
    | io.go
sdl
    | loop.go
    | window.go
util
    | cell.go
    | check.go
    | vissualise.go
test...

```

#### distribute
```txt
client
    | gol
    | sdl
    | util

remote
    | remoteutil
    | server
```