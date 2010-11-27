package timestep

import (
    "fmt"
)

const (
    IterTime = 30 // Minutes
)

type Tstep int64 // Iteration time since start

func (ts *Tstep) Print() string {
    return fmt.Sprint(*ts * IterTime / 60)
}

func init() {
    // Start the stepping goroutine
    go func() {
        //
        var Iter Tstep = 0
        for {
            master <- Iter
            for _, c := range chs {
                c <- Iter
            }
            Iter++
        }
        return
    }()
}

var firstPass bool = true
var master chan Tstep = make(chan Tstep)
var chs []chan Tstep = make([]chan Tstep, 0, 10)

func MasterTimeStep() chan Tstep {
    // If first time called, return master channel, else append new channel to slice
    var ch chan Tstep
    if firstPass {
        firstPass = false
        ch = master
    } else {
        ch = TimeStep()
    }
    return ch
}

func TimeStep() chan Tstep {
    // Append new channel to slice
    c := make(chan Tstep)
    chs = append(chs, c)
    return c
}
