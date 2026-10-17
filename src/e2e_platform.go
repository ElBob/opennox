package opennox

import (
	"math/rand"
	"os"
	"sync"
	"time"
)

// e2eRealtime makes Sleep wait for real, so the client can talk to an external server (e.g. a dedicated one).
// Virtual time is not advanced by Sleep in this mode, so script steps are still scheduled by frames
// (use NOX_E2E_SLOW to pace frames). Network timing is not deterministic, so frames may differ between runs.
var e2eRealtime = os.Getenv("NOX_E2E_REALTIME") == "true"

func newPlayformE2E() *platformE2E {
	p := &platformE2E{real: e2eRealtime}
	p.RandSeed(1234)
	return p
}

type platformE2E struct {
	mu   sync.Mutex
	t    time.Duration
	rnd  *rand.Rand
	real bool
}

func (p *platformE2E) ticks() time.Duration {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.t
}

func (p *platformE2E) tick(dt time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.t += dt
}

func (p *platformE2E) Ticks() time.Duration {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.t
}

func (p *platformE2E) Sleep(dt time.Duration) {
	if p.real {
		time.Sleep(dt)
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.t += dt
}

func (p *platformE2E) TimeSeed() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return int64(p.t)
}

func (p *platformE2E) RandInt() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.rnd.Int()
}

func (p *platformE2E) RandSeed(seed int64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.rnd = rand.New(rand.NewSource(seed))
}

func (p *platformE2E) RandSeedTime() {
	p.RandSeed(p.TimeSeed())
}
