package pipes

import (
	"context"
	"sync"

	"golang.org/x/xerrors"
)

// first in - first out, simple implementation for something that needs to
// be processed in a linear fashion
type fifo struct {
	proc Processor
}

func FIFO(proc Processor) StageRunner {
	return fifo{proc: proc}
}

func (r fifo) Run(ctx context.Context, parameters StageParameters) {
	for {
		// non blocking operation ahead
		select {
		case <-ctx.Done():
			return
		case payloadIn, ok := <-parameters.Input():
			if !ok {
				//asked to return if no data
				return
			}

			payloadOut, err := r.proc.Process(ctx, payloadIn)
			if err != nil {
				wrappedErr := xerrors.Errorf("pipeline stage %d: %w", parameters.StageIndex(), err)
				maybeEmitError(wrappedErr, parameters.Error())
				return
			}

			// this one is useful for when we want to discard payloads for some
			// reason
			if payloadOut == nil {
				payloadIn.MarkAsProcessed()
			}

			select {
			case parameters.Output() <- payloadOut:
			case <-ctx.Done():
				// early return if context timeouts
				return
			}
		}
	}
}

type fixedWorkerPool struct {
	fifos []StageRunner
}

func FixedWorkerPool(proc Processor, numWorkers int) StageRunner {
	if numWorkers <= 0 {
		panic("FixedWorkerPool: numWorkers must be > 0")
	}

	fifos := make([]StageRunner, numWorkers)
	for i := 0; i < numWorkers; i++ {
		fifos[i] = FIFO(proc)
	}

	return &fixedWorkerPool{fifos: fifos}
}

func (p *fixedWorkerPool) Run(ctx context.Context, parameters StageParameters) {
	var wg sync.WaitGroup

	// Spin each worker (that is just a FIFO stage) and wait for them to
	// exit

	for i := 0; i < len(p.fifos); i++ {
		wg.Add(1)
		go func(fifoIndex int) {
			p.fifos[fifoIndex].Run(ctx, parameters)
			wg.Done()
		}(i)
	}

	wg.Wait()
}

type dynamicWorkerPool struct {
	proc      Processor
	tokenPool chan struct{}
}

func DynamicWorkerPool(proc Processor, maxWorkers int) StageRunner {
	if maxWorkers <= 0 {
		panic("DynamicWorkerPool: numWorkers must be > 0")
	}

	tokenPool := make(chan struct{}, maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		tokenPool <- struct{}{}
	}

	return &dynamicWorkerPool{proc: proc, tokenPool: tokenPool}
}

func (d *dynamicWorkerPool) Run(ctx context.Context, parameters StageParameters) {
stop:
	for {
		// non blocking operation ahead
		select {
		case <-ctx.Done():
			// fail early
			break stop
		case payloadIn, ok := <-parameters.Input():
			if !ok {
				// fail early
				break stop
			}

			var token struct{}
			select {
			case token = <-d.tokenPool: // get token from the pool
			case <-ctx.Done():
				// fail early
				break stop
			}

			go func(payloadIn Payload, token struct{}) {
				defer func() {
					// retrieve token
					// after execution
					d.tokenPool <- token
				}()
				payloadOut, err := d.proc.Process(ctx, payloadIn)
				if err != nil {
					wrappedErr := xerrors.Errorf("pipeline stage %d: %w", parameters.StageIndex(), err)
					maybeEmitError(wrappedErr, parameters.Error())
					return
				}

				// for removing payloads that we don't want to process on the
				// fly, this removes them from the pipeline altogether.
				if payloadOut == nil {
					payloadIn.MarkAsProcessed()
					return
				}

				select {
				case parameters.Output() <- payloadOut:
				case <-ctx.Done():
					return
				}

			}(payloadIn, token)
		}
	}

	// emptying the  token pool to exit
	for i := 0; i < cap(d.tokenPool); i++ {
		<-d.tokenPool
	}
}

type broadcast struct {
	fifos []StageRunner
}

func Broadcast(procs ...Processor) StageRunner {
	if len(procs) == 0 {
		panic("Broadcast: at least one processor must be specified")
	}

	fifos := make([]StageRunner, len(procs))
	for i, p := range procs {
		fifos[i] = FIFO(p)
	}

	return &broadcast{fifos: fifos}
}

func (b *broadcast) Run(ctx context.Context, parameter StageParameters) {
	var (
		wg   sync.WaitGroup
		inCh = make([]chan Payload, len(b.fifos))
	)

	// start each FIFO in a goroutine, allocate an input channel for
	// every subworker and poss the output channel to send further down
	// the payloads
	for i := 0; i < len(b.fifos); i++ {
		wg.Add(1)
		go func(fifoIndex int) {
			fifoparams := &workerParameters{
				stage: parameter.StageIndex(),
				inCh:  inCh[fifoIndex],
				outCh: parameter.Output(),
				errch: parameter.Error(),
			}
			b.fifos[fifoIndex].Run(ctx, fifoparams)
			wg.Done()
		}(i)
	}

done:
	for {
		// non blocking operation incoming
		select {
		case <-ctx.Done():
			// fail early
			return
		case payload, ok := <-parameter.Input():
			if !ok {
				break done
			}
			for i := len(b.fifos) - 1; i >= 0; i-- {
				// To avoid a data race situation
				// we clone each payload in order to paralellize
				// the processing
				var fifoPayload = payload
				if i != 0 {
					fifoPayload = payload.Clone()
				}
				select {
				case <-ctx.Done():
					break done
				case inCh[i] <- fifoPayload:
					// Sent the payload to the i-th FIFO
				}
			}
		}
	}

	// closing each and every channel to signal the next stage
	// that we're done.
	for _, ch := range inCh {
		close(ch)
	}

	// wait for all the goroutines to end and exit
	wg.Wait()
}
