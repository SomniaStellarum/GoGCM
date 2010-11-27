package timestep

import (
    "testing"
)

var ch chan Tstep = MasterTimeStep()
var ch1 chan Tstep = TimeStep()

func TestTimeStep(t *testing.T) {
    const n = 1000
    for i := 0; i < n; i++ {
        ts := <-ch
        ts2 := <-ch1
        if ts != Tstep(i) {
            t.Errorf("Master timestep wrong! %d (Expected: %d)", ts, i)
        }
        if ts2 != Tstep(i) {
            t.Errorf("Master timestep wrong! %d (Expected: %d)", ts2, i)
        }
    }
}

func BenchmarkTimeStep(b *testing.B) {
    for i := 0; i < b.N; i++ {
        <-ch
        <-ch1
    }
}
