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


### tests and todo related

#### stage 1
1. benchmark modification
        
#### stage 2
1. 12/03
    - stage 2 started

#### report related
>Report (30 marks)
 
>You need to submit a CONCISE (strictly max 6 pages) report which should cover the following topics:
 
>Functionality and Design: Outline what functionality you have implemented, which problems you have solved with your implementations and how your program is designed to solve the problems efficiently and effectively.
 
>Critical Analysis: Describe the experiments and analysis you carried out. Provide a selection of appropriate results. Keep a history of your implementations and provide benchmark results from various stages. Explain and analyse the benchmark results obtained. Analyse the important factors responsible for the virtues and limitations of your implementations.
 
>Make sure your team member’s names and user names appear on page 1 of the report. Do not include a cover page.


### notes

#### 11/24
1. the `ioInput` and `ioOutput` is linked together to share the world, in this way, we don't need an extra `tempworld` channel to pass the changes in world in each round.
2. we need to input and output the world in each round, either before the calculation or after the modification.

#### 11/26
1. stage one done with single threads

#### 11/30
1. benchmark tests
```linux
go get github.com/ChrisGora/benchgraph
go test -bench . | benchgraph
```

#### 12/02
1. report starts

#### 12/07
1. benchmark done








