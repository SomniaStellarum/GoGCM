package solar

import (
    timestep "github.com/Arrow/GoGCM/util/timestep"
    "rand"
)

const (
    S_Start = 1366.1
)

func Solar() (out chan float64) {
    out = make(chan float64)
    ts := timestep.TimeStep()
    go func() {
        S_Var := float64(S_Start)
        for {
            select {
            case <-ts:
                S_Var = S_Var + 5*rand.NormFloat64()
            case out <- S_Var:
            }
        }
        close(out)
    }()
    return out
}
