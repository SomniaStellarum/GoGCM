package main

import (
    "fmt"
    "rand"
    "math"
)

const (
    A = 0.306
    Sigma = 5.6704e-8
    Sigma_Inv = 1 / 5.6704e-8
    S_Start = 1366.1
    T_Start = 305
    EarthR = 6371 // in km - Mean Radius
    AbsZero = 273.15
    CpConst = 10000000000000
    KConst = 100000
    KM = 1000
)

type Input struct {
    NumGridpoints int
    FluxTransfer []int // Contains a list of all flux transfers to be calculated.
    FluxIndex []int
}
// Flux Index lists the indexes that delimit a slice for each gridpoint. This slice indicates
// to which other gridpoints it calculates a flux to.
// eg Take gridpoint 1 of a cartesian grid with 4 gridpoints (North - West, North - East, ...etc)
// Each gridpoint is at a pole, thus there are only 3 gridpoints to be transfered to.
// Flux Index starts with [0,2, ...] and FluxTransfer starts with [2,3,4, ...]

type Gridpoint struct {
    Area float64
    BoundL []float64     // Boundary Lengths Same convention as Fluxes (0,1,2 maybe 3)
    Lat float64
    Long float64
    A float64
    Cp []float64         // 0: Surf, 1: Atmo1, 2: Atmo2
    K float64
    NPole bool
    SPole bool
}

type FluxComponent struct {
    Fc float64           // Calculated flux component
    Idx int              // Index of source gridpoint
}

type Datapoint struct {
    F []float64          // Contains all fluxes for each gridpoint
    Temp []float64       // Contains all temperatures for each gridpoint
    Idx int              // Index of source gridpoint
    *Gridpoint
}
// Temperature index, 0: Surface, 1: Atmosphere layer 1, 2: Atmosphere layer 2
// Flux index, 0: Solar, 1: Surf to Atmo1, 2: Atmo1 to Surf, 3: Atmo1 to Atmo2, 4: Atmo2 to Atmo1
// 5: Atmo2 out, 6-9: Conduction N-E-S-W (Skip N at North Pole and S at South Pole)
// Fluxes over Area are in W/m^2 (must divide by Area before calculating Temperature)
// (0, 1, 2, 3, 4 and 5)
// Fluxes over Boundary are in W/m (must divide by boundary length before calculating Temperature)
// (6, 7, 8 and possibly 9)

func NewDatapoint(Temp, Flux []float64) (dt *Datapoint) {
    return &Datapoint{Flux, Temp, 1,
        &Gridpoint{1,make([]float64,3,4),0,0,0.2,make([]float64,3),1,false,false}}
}

func (dt *Datapoint) String() string {
    s := fmt.Sprintf("%f",dt.Temp[0])
    for i := 1; i < len(dt.Temp); i++ {
        s = fmt.Sprintf("%s,%f", s, dt.Temp[i])
    }
    for i := 0; i < len(dt.F); i++ {
        s = fmt.Sprintf("%s,%f", s, dt.F[i])
    }
    return fmt.Sprintf("%s\n", s)
}

func Flux(ch1 chan *Datapoint, ch2 chan *Datapoint, in Input) (out chan *Datapoint) {
    out = make(chan *Datapoint)
    fc := make(chan *FluxComponent) // Channel to receive flux components from each gridpoint
    id := make(chan int)
    fc_ret := make(chan float64)
    go func() {
        // This goroutine is to accept flux components, store it in a slice
        // and return the transfer fluxes
        flux := make([]float64, in.NumGridpoints)
        for {
            select {
            case f := <-fc:
                flux[f.Idx] = f.Fc
            case idx := <-id:
                ft := in.FluxTransfer[in.FluxIndex[idx*2]:in.FluxIndex[idx*2 + 1]]
                for _, ftIdx := range ft {
                    fc_ret <- flux[ftIdx-1] - flux[idx]
                }
            }
        }
    } ()
    go func() {
        S_Var := float64(S_Start)
        for dt := range ch1 {
            fc <- &FluxComponent{dt.K * dt.Temp[0], dt.Idx}
            S_Var = S_Var + 5 * rand.NormFloat64()
            dt.F[0] = S_Var * 0.25 * (1 - dt.A) //solar effect
            dt.F[1] = Sigma * math.Pow(dt.Temp[0],4)
            dt.F[2] = Sigma * math.Pow(dt.Temp[1],4)
            dt.F[3] = dt.F[2]
            dt.F[4] = Sigma * math.Pow(dt.Temp[2],4)
            dt.F[5] = dt.F[4]
            out <- dt
        }
    } ()
    go func() {
        for dt := range ch2 {
            id <- dt.Idx
            dt.F[6] = <-fc_ret
            dt.F[7] = <-fc_ret
            dt.F[8] = <-fc_ret
            if (!(dt.NPole || dt.SPole)) {
                fmt.Printf("Shouldn't be here\n")
                dt.F[9] = <-fc_ret
            }
            out <- dt
        }
        close(out)
    } ()
    return out
}

func gcm(ch chan *Datapoint) (out chan *Datapoint) {
    out = make(chan *Datapoint)
    go func() {
        Q := float64(0)
        for dt := range ch {
            //Calculate Temperatures
            Q = dt.Area * (dt.F[0] - dt.F[1] + dt.F[2])
            Q += dt.BoundL[0] * dt.F[6]
            Q += dt.BoundL[1] * dt.F[7]
            Q += dt.BoundL[2] * dt.F[8]
            if (!(dt.NPole || dt.SPole)) {
                Q += dt.BoundL[3] * dt.F[9]
            }
            dt.Temp[0] += Q/dt.Cp[0]
            
            Q = dt.Area * (dt.F[1] - dt.F[2] - dt.F[3] + dt.F[4])
            dt.Temp[1] += Q/dt.Cp[1]
            
            Q = dt.Area * (dt.F[3] - dt.F[4] - dt.F[5])
            dt.Temp[2] += Q/dt.Cp[2]
            out <- dt
        }
        close(out)
    }()
    return out
}

func main() {
    dt := make([]Datapoint,4)
    
    for i := 0; i < 4; i++ {
        t := make([]float64,3)
        t[0] = 335//AbsZero + 60
        t[1] = 303//AbsZero + 30
        t[2] = 254//AbsZero - 20
        f := make([]float64,9)
        dt[i] = *NewDatapoint(t,f)
        fmt.Print(dt[i])
        dt[i].Idx = i
        dt[i].Area = KM * math.Pi * math.Pow(EarthR,2) // 4 pi r^2 / (4 gridpoints)
        dt[i].BoundL = []float64{KM*math.Pi*EarthR/2,KM*math.Pi*EarthR/2,KM*math.Pi*EarthR/2}
        dt[i].Cp = []float64{CpConst,CpConst,CpConst}
        dt[i].K = KConst
        dt[i].A = A
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
    
    in := &Input{4,[]int{2,3,4,1,3,4,1,2,4,1,2,3},[]int{0,3,3,6,6,9,9,12}}
    
    ch1 := make(chan *Datapoint)
    ch2 := make(chan *Datapoint)
    ch3 := make(chan *Datapoint)
    
    out1 := Flux(ch1, ch2, *in)
    out2 := gcm(ch3)
    for i := 0; i < 10; i++ {
        // Running the simulation
        fmt.Printf("Run %d!\n",i)
        for _, d := range dt {
            ch1 <- &d
            <-out1
        }
        for _, d := range dt {
            ch2 <- &d
            <-out1
        }
        for _, d := range dt {
            ch3 <- &d
            <-out2
        }
        for _, d := range dt {
            fmt.Print(&d)
        }
    }
}
