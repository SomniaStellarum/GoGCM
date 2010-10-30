package main

import (
    atm "github.com/Arrow/GoGCM/atmos"
    "fmt"
    "math"
)


func main() {
    dt := make([]atm.Datapoint, 4)

    for i := 0; i < 4; i++ {
        t := make([]float64, 3)
        t[0] = 335 //AbsZero + 60
        t[1] = 303 //AbsZero + 30
        t[2] = 254 //AbsZero - 20
        f := make([]float64, 9)
        dt[i] = *atm.NewDatapoint(t, f)
        fmt.Print(&dt[i])
        dt[i].Idx = i
        dt[i].Area = atm.KM * math.Pi * math.Pow(atm.EarthR, 2) // 4 pi r^2 / (4 gridpoints)
        dt[i].BoundL = []float64{atm.KM * math.Pi * atm.EarthR / 2, atm.KM * math.Pi * atm.EarthR / 2, atm.KM * math.Pi * atm.EarthR / 2}
        dt[i].Cp = []float64{atm.CpConst, atm.CpConst, atm.CpConst}
        dt[i].K = atm.KConst
        dt[i].A = atm.A
    }
    dt[0].Lat = 45.0
    dt[0].Long = -45.0
    dt[0].NPole = true
    dt[1].Lat = 45.0
    dt[1].Long = 45.0
    dt[1].NPole = true
    dt[2].Lat = -45.0
    dt[2].Long = -45.0
    dt[2].SPole = true
    dt[3].Lat = -45.0
    dt[3].Long = 45.0
    dt[3].SPole = true

    FluxIn := &atm.FluxInput{4, []int{2, 3, 4, 1, 3, 4, 1, 2, 4, 1, 2, 3}, []int{0, 3, 3, 6, 6, 9, 9, 12},
        make(chan *atm.Datapoint), make(chan *atm.FluxComponent)}

    ch1 := FluxIn.ChFlux

    GCMIn := &atm.GCMInput{4, make(chan *atm.Datapoint), FluxIn.ChFluxComp}
    ch2 := GCMIn.Ch

    out1 := atm.Flux(*FluxIn)
    out2 := atm.Gcm(*GCMIn)
    for i := 0; i < 10; i++ {
        // Running the simulation
        fmt.Printf("Run %d!\n", i)
        for _, d := range dt {
            ch1 <- &d
            <-out1
        }
        for _, d := range dt {
            ch2 <- &d
            <-out2
        }
        for _, d := range dt {
            fmt.Print(&d)
        }
    }
}
