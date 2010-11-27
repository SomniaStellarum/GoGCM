package solar

import (
	timestep "github.com/Arrow/GoGCM/util/timestep"
    "testing"
    "math"
)

var ch chan timestep.Tstep = timestep.MasterTimeStep()
var ch_sol chan float64 = Solar()

func notNearEqual(a, b, closeEnough float64) bool {
    absDiff := math.Fabs(a - b)
    if absDiff < closeEnough {
        return false
    }
    return true
}

func TestSolar(t *testing.T) {
    const n = 100
    for i := 0; i < n; i++ {
        s := <-ch_sol
        if notNearEqual(s, S_Start, float64((i+1)*5)) {
            t.Errorf("Not close enough! %d (Expected: %d)", s, S_Start)
        }
    }
}

func BenchmarkSolar(b *testing.B) {
    const n = 1
    for i := 0; i < b.N; i++ {
        <-ch
        for j := 0; j < n; j++ {
            <-ch_sol
        }
    }
}
