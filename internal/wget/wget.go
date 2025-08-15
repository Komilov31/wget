package wget

import (
	"github.com/Komilov31/wget/internal/worker"
)

type Wget struct {
	worker *worker.Worker
}

func New(worker *worker.Worker) *Wget {
	return &Wget{
		worker: worker,
	}
}

func (w *Wget) Run() {
	w.worker.Start()
}
