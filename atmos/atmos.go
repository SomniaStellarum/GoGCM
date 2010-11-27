package atmos

import (
	solar "github.com/Arrow/GoGCM/solar"
	"fmt"
	//"rand"
	"math"
)

const (
	A         = 0.306
	Sigma     = 5.6704e-8
	Sigma_Inv = 1 / 5.6704e-8
	T_Start   = 305
	EarthR    = 6371 // in km - Mean Radius
	AbsZero   = 273.15
	CpConst   = 10000000000000
	KConst    = 100000
	KM        = 1000
)

type FluxInput struct {
	NumGridpoints int
	FluxTransfer  []int // Contains a list of all flux transfers to be calculated.
	FluxIndex     []int
	ChFlux        chan *Datapoint
	ChFluxComp    chan *FluxComponent
}

type GCMInput struct {
	NumGridpoints int
	Ch            chan *Datapoint
	ChFluxComp    chan *FluxComponent
}
// Flux Index lists the indexes that delimit a slice for each gridpoint. This slice indicates
// to which other gridpoints it calculates a flux to.
// eg Take gridpoint 1 of a cartesian grid with 4 gridpoints (North - West, North - East, ...etc)
// Each gridpoint is at a pole, thus there are only 3 gridpoints to be transfered to.
// Flux Index starts with [0,2, ...] and FluxTransfer starts with [2,3,4, ...]

type Gridpoint struct {
	Area   float64
	BoundL []float64 // Boundary Lengths Same convention as Fluxes (0,1,2 maybe 3)
	Lat    float64
	Long   float64
	A      float64
	Cp     []float64 // 0: Surf, 1: Atmo1, 2: Atmo2
	K      float64
	NPole  bool
	SPole  bool
}

type FluxComponent struct {
	Fc  float64 // Calculated flux component
	Idx int     // Index of source gridpoint
}

type Datapoint struct {
	F    []float64 // Contains all fluxes for each gridpoint
	Temp []float64 // Contains all temperatures for each gridpoint
	Idx  int       // Index of source gridpoint
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
		&Gridpoint{1, make([]float64, 3, 4), 0, 0, 0.2, make([]float64, 3), 1, false, false}}
}

func (dt *Datapoint) String() string {
	s := fmt.Sprintf("%f", dt.Temp[0])
	for i := 1; i < len(dt.Temp); i++ {
		s = fmt.Sprintf("%s,%f", s, dt.Temp[i])
	}
	for i := 0; i < len(dt.F); i++ {
		s = fmt.Sprintf("%s,%f", s, dt.F[i])
	}
	return fmt.Sprintf("%s\n", s)
}

func Flux(in FluxInput) (out chan *Datapoint) {
	out = make(chan *Datapoint)
	id := make(chan int)
	fc_ret := make(chan float64)
	go func() {
		// This goroutine is to accept flux components, store it in a slice
		// and return the transfer fluxes
		flux := make([]float64, in.NumGridpoints)
		for {
			select {
			case f := <-in.ChFluxComp:
				flux[f.Idx] = f.Fc
			case idx := <-id:
				ft := in.FluxTransfer[in.FluxIndex[idx*2]:in.FluxIndex[idx*2+1]]
				for _, ftIdx := range ft {
					fc_ret <- flux[ftIdx-1]-flux[idx]
				}
			}
		}
	}()
	go func() {
		ch_sol := solar.Solar()
		for dt := range in.ChFlux {
			s := <-ch_sol
			dt.F[0] = s * 0.25 * (1 - dt.A) //solar effect
			dt.F[1] = Sigma * math.Pow(dt.Temp[0], 4)
			dt.F[2] = Sigma * math.Pow(dt.Temp[1], 4)
			dt.F[3] = dt.F[2]
			dt.F[4] = Sigma * math.Pow(dt.Temp[2], 4)
			dt.F[5] = dt.F[4]
			
			id <- dt.Idx
			dt.F[6] = <-fc_ret
			dt.F[7] = <-fc_ret
			dt.F[8] = <-fc_ret
			if !(dt.NPole || dt.SPole) {
				dt.F[9] = <-fc_ret
			}
			out <- dt
		}
		close(out)
	}()
	return out
}

func Gcm(in GCMInput) (out chan *Datapoint) {
	out = make(chan *Datapoint)
	go func() {
		Q := float64(0)
		for dt := range in.Ch {
			//Calculate Temperatures
			Q = dt.Area * (dt.F[0] - dt.F[1] + dt.F[2])
			Q += dt.BoundL[0] * dt.F[6]
			Q += dt.BoundL[1] * dt.F[7]
			Q += dt.BoundL[2] * dt.F[8]
			if !(dt.NPole || dt.SPole) {
				Q += dt.BoundL[3] * dt.F[9]
			}
			dt.Temp[0] += Q / dt.Cp[0]

			Q = dt.Area * (dt.F[1] - dt.F[2] - dt.F[3] + dt.F[4])
			dt.Temp[1] += Q / dt.Cp[1]

			Q = dt.Area * (dt.F[3] - dt.F[4] - dt.F[5])
			dt.Temp[2] += Q / dt.Cp[2]
			
			in.ChFluxComp <- &FluxComponent{dt.K * dt.Temp[0], dt.Idx}
			
			out <- dt
		}
		close(out)
	}()
	return out
}
