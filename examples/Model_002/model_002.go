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
    T_Start = 305
)

type datapoint struct {
    count int
    S float64
    F float64
    temp []float64
}

func (dt *datapoint) String() string {
    return fmt.Sprintf("%d,%f,%f,%f,%f\n", dt.count, dt.S, dt.F, dt.temp[0], dt.temp[1])
}

func solar(ch chan *datapoint) (out chan *datapoint) {
    out = make(chan *datapoint)
    go func() {
        S_Var := float64(S_Start)
        for dt := range ch {
            switch dt.F {
            case 0.0: 
                S_Var = S_Var + 5 * rand.NormFloat64()
                dt.S, dt.F = S_Var, S_Var * 0.25 * (1 - A)
            default:
                S_Var = dt.F * 4 / (1 - A)
                dt.S = S_Var
            }
            out <- dt
        }
        close(out)
    }()
    return out
}            

func gcm(ch chan *datapoint) (out chan *datapoint) {
    out = make(chan *datapoint)
    go func() {
        TempS := float64(T_Start)
        TempA := float64(T_Start)
        for dt := range ch {
            switch dt.F {
            case 0.0:
                TempS = TempS + 0.05 + 0.08 * rand.NormFloat64()
                TempA = TempS / math.Pow(2,0.25)
                dt.temp[0], dt.temp[1], dt.F = TempS, TempA, math.Pow(TempA, 4) / Sigma_Inv
            default:
                TempA = math.Pow(dt.F * Sigma_Inv, 0.25)
                TempS = math.Pow(2,0.25) * TempA
                dt.temp[0], dt.temp[1] = TempS, TempA
            }
            out <- dt
        }
        close(out)
    }()
    return out
}

func Solar_Drive() {
    dt := &datapoint{0,0,0,[]float64{0,0}}
    ch1 := make(chan *datapoint)
    ch2 := solar(ch1)
    ch3 := gcm(ch2)
    ch1 <- dt
    for dt = range ch3 {
        fmt.Print(dt)
        if (dt.count > 99) {
            return
        }
        dt.count++
        dt.S, dt.F, dt.temp[0], dt.temp[1] = 0,0,0,0
        ch1 <- dt
    }
    close(ch1)
}

func Temp_Drive() {
    dt := &datapoint{0,0,0,[]float64{0,0}}
    ch1 := make(chan *datapoint)
    ch2 := gcm(ch1)
    ch3 := solar(ch2)
    ch1 <- dt
    for dt = range ch3 {
        fmt.Print(dt)
        if (dt.count > 99) {
            return
        }
        dt.count++
        dt.S, dt.F, dt.temp[0], dt.temp[1] = 0,0,0,0
        ch1 <- dt
    }
    close(ch1)
}

func IntMod(in, div int) (out int) {
    tmp := in
    for {
        if(tmp < div) {
            out = tmp
            return out
        }
        tmp = tmp - div 
    }
    return out
}

func Hybrid_Drive() {
    dt := &datapoint{0,0,0,[]float64{0,0}}
    Sch_In := make(chan *datapoint)
    Gch_In := make(chan *datapoint)
    Sch_Out := solar(Sch_In)
    Gch_Out := gcm(Gch_In)
    Ch_Out := make(chan *datapoint)
    go func () {
        for {
            select {
            case in := <-Sch_Out:
                Ch_Out <- in
            case in := <-Gch_Out:
                Ch_Out <- in
            }
            if (closed(Sch_Out) && closed(Gch_Out)) {
                close(Ch_Out)
                return
            }
        }
    }()
    Gch_In <- dt
    for dt = range Ch_Out {
        switch {
        case (dt.temp[0] == 0): Gch_In <- dt
        case (dt.S == 0): Sch_In <- dt
        default:
            fmt.Print(dt)
            dt.count++
            if (dt.count > 100) {
                return
            }
            dt.S, dt.F, dt.temp[0], dt.temp[1] = 0,0,0,0
            if (IntMod(dt.count,5) == 0) {
                Gch_In <- dt
            } else {
                Sch_In <- dt
            }
        }
    }
}

func main() {
    Solar_Drive()
    Temp_Drive()
    Hybrid_Drive()
}
