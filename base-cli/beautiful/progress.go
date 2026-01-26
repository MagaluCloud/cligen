package beautiful

import (
	"sync/atomic"

	"github.com/pterm/pterm"
)

type PTermProgressAdapter struct {
	bar    *pterm.ProgressbarPrinter
	closed atomic.Bool
}

func NewPTermProgress(
	total *int64,
	title *string,
) *PTermProgressAdapter {

	bar, _ := pterm.DefaultProgressbar.
		WithShowCount(true).
		WithRemoveWhenDone(false).
		Start()

	if total != nil {
		bar.WithTotal(int(*total))
	}

	if title != nil {
		bar.UpdateTitle(*title)
	}

	return &PTermProgressAdapter{bar: bar}
}

func (p *PTermProgressAdapter) Start(total int64) {
	if p.closed.Load() {
		return
	}

	if total == 0 {
		p.bar.Total = 1
		_ = p.bar.Add(1)
		p.Finish()
		return
	}

	p.bar.Total = int(total)
	_ = p.bar.UpdateTitle(p.bar.Title)
}

func (p *PTermProgressAdapter) Add(n int64) {
	if p.closed.Load() {
		return
	}

	_ = p.bar.Add(int(n))
}

func (p *PTermProgressAdapter) Finish() {
	if p.closed.Swap(true) {
		return
	}

	_, _ = p.bar.Stop()
}
