package main

import (
    "fmt"
    "rand"
    "math"
)

const (
    A = 0.306
    Sigma_Inv = 1 / 5.6704e-8
    S_Start = 1366.1
)

type datapoint struct {
    count int
    S float64
    F float64
    temp  float64
}

func (dt *datapoint) String() string {
    return fmt.Sprintf("%d,%f,%f,%f\n", dt.count, dt.S, dt.F, dt.temp)
}

func counter() (ch chan int) {
    ch = make(chan int)
    go func() {
        for i := 0; i <= 100; i++ {
            ch <- i
        }
        close(ch)
    }()
    return ch
}

func solar(ch chan int) (out chan *datapoint) {
    out = make(chan *datapoint)
    go func() {
        S_Var := float64(S_Start)
        for c := range ch {
            S_Var = S_Var + 5 * rand.NormFloat64()
            dt := datapoint{c, S_Var, S_Var * 0.25 * (1 - A), 0}
            out <- &dt
        }
        close(out)
    }()
    return out
}            

func gcm(ch chan *datapoint) (out chan *datapoint) {
    out = make(chan *datapoint)
    go func() {
        for dt := range ch {
            dt.temp = math.Pow(dt.F * Sigma_Inv, 0.25)
            out <- dt
        }
        close(out)
    }()
    return out
}

func main() {
    ch_i := counter()
    ch_f := solar(ch_i)
    ch_f = gcm(ch_f)
    for in := range ch_f {
        fmt.Print(in)
    }
}
