package main

import (
    "fmt"
    "rand"
)

const (
    Mean = 17.5
)

type datapoint struct {
    count int
    temp  float64
}

func (dt datapoint) String() string {
    return fmt.Sprintf("%d,%f\n", dt.count, dt.temp)
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

func gcm(ch chan int) (out chan datapoint) {
    out = make(chan datapoint)
    go func() {
        dt := new(datapoint)
        for c := range ch {
            dt.count = c
            dt.temp = Mean + rand.NormFloat64()
            out <- *dt
        }
        close(out)
    }()
    return out
}

func main() {
    ch := counter()
    inChan := gcm(ch)
    for in := range inChan {
        fmt.Print(in)
    }
}
