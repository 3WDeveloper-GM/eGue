package pipes

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var _ StageParameters = (*workerParameters)(nil)

type workerParameters struct {
	stage int

	// Channels for the propagation of inputs, outputs, and
	// errors
	inCh  <-chan Payload
	outCh chan<- Payload
	errch chan<- error
}

func (p *workerParameters) StageIndex() int        { return p.stage }
func (p *workerParameters) Input() <-chan Payload  { return p.inCh }
func (p *workerParameters) Output() chan<- Payload { return p.outCh }
func (p *workerParameters) Error() chan<- error    { return p.errch }

// Implements a modular multi-stage pipeline. Each pipeline is constructed
// out of an input source, an output sink and zero or more
// processing stages
type Pipeline struct {
	stages []StageRunner
}

func New(stages ...StageRunner) *Pipeline {
	return &Pipeline{
		stages: stages,
	}
}

func (p *Pipeline) Process(ctx context.Context, source Source, sink Sink) error {
	var wg sync.WaitGroup

	pCtx, ctxCancelFunc := context.WithCancel(ctx)

	// Channel allocation for wiring the source, stages, and sink together.
	// The output of the n-th stage is fed to the n+1-th stage and so on and
	// so forth.

	// Allocation of an extra channel for the wiring of the source/sink
	stageCh := make([]chan Payload, len(p.stages)+1)

	// Allocation for an extra two channels in case something happens at the
	// source or sink level.

	errCh := make(chan error, len(p.stages)+2)
	for i := 0; i < len(stageCh); i++ {
		// populate each of the stages with a payload channel
		stageCh[i] = make(chan Payload)
	}

	// Start a worker for each stage, this is equivalent to spin a goroutine for
	// each stage.

	for i := 0; i < len(p.stages); i++ {
		wg.Add(1)
		go func(stageIndex int) {
			p.stages[stageIndex].Run(ctx, &workerParameters{
				stage: stageIndex,
				inCh:  stageCh[stageIndex],
				outCh: stageCh[stageIndex+1],
				errch: errCh,
			})

			// Signal next stage that we've run out of data.
			close(stageCh[stageIndex+1])
			wg.Done()
		}(i)
	}

	wg.Add(2)
	go func() {
		// spin up the source worker at the first stage of the pipe
		sourceWorker(pCtx, source, stageCh[0], errCh)
		close(stageCh[0])
		wg.Done()
	}()

	go func() {
		// spin up the new go routine for the
		sinkWorker(pCtx, sink, stageCh[len(stageCh)-1], errCh)
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(errCh)
		ctxCancelFunc()
	}()

	var err error
	for pErr := range errCh {
		err = errors.Join(err, pErr)
		ctxCancelFunc()
	}

	return err

}

func sinkWorker(ctx context.Context, sink Sink, inCh <-chan Payload, errCh chan<- error) {
	for {
		// non blocking operation ahead
		select {
		case payload, ok := <-inCh:
			if !ok {
				return
			}

			if err := sink.Consume(ctx, payload); err != nil {
				wrapper := "pipeline sink : %w"
				wrappedErr := fmt.Errorf(wrapper, err)
				maybeEmitError(wrappedErr, errCh)
				return
			}
			payload.MarkAsProcessed()
		case <-ctx.Done():
			// shutdown due to a context timeout
			return
		}
	}
}

func sourceWorker(ctx context.Context, source Source, outCh chan<- Payload, errCh chan<- error) {
	for source.Next(ctx) {
		Payload := source.Payload()
		select {
		case outCh <- Payload:
		case <-ctx.Done():
			// Shutdown due to a context timeout.
			return
		}
	}

	if err := source.Error(); err != nil {
		wrapper := "pipeline source : %w"
		wrappedErr := fmt.Errorf(wrapper, err)
		maybeEmitError(wrappedErr, errCh)
	}
}

func maybeEmitError(err error, errCh chan<- error) {
	// select for non blocking operation
	select {
	case errCh <- err: // error emitted
	default: //error channel full with other errors
	}
}
