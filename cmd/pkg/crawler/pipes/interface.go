package pipes

import "context"

// generic data input for a pipeline
type Payload interface {
	//deep-copy for concurrency safety
	Clone() Payload

	//MarkasProcessed is invoked to either signal that the payload reached the end
	//or has been discarded mid flight
	MarkAsProcessed()
}

// generic data processing stage, process payloads in a set order
type Processor interface {
	//Process operates on the input payload and returns back a new payload
	//to be forwarded to the next stage.
	Process(context.Context, Payload) (Payload, error)
}

// closure for using a simple function as a processor, it must have the
// signature that is written down below
type ProcessorFunc func(context.Context, Payload) (Payload, error)

// Calls the Process() method for the close f.
func (f ProcessorFunc) Process(ctx context.Context, P Payload) (Payload, error) {
	return f(ctx, P)
}

// StageParameters encapsulates the information required to structure the pipeline
// It can be seen as the interface that defines the ins and outs of a certain stage.
type StageParameters interface {
	StageIndex() int
	Input() <-chan Payload
	Output() chan<- Payload
	Error() chan<- error
}

// Implements a way to define a multi-stage pipeline
type StageRunner interface {
	// Run implements the logic for the stage in terms of processing inputs
	// and writing outputs to an output channel.
	Run(context.Context, StageParameters)
}

// Source implements the types that generate Payload instances which can be used as inputs to a
// pipeline instance
type Source interface {
	// Is true is there is more data to be read and processed.
	Next(context.Context) bool

	//Payload returns the next payload to be processed.
	Payload() Payload

	//Error return the last error observed by the source
	Error() error
}

// Implements the types that act as the tail of the pipeline.
type Sink interface {
	// Consumes a Payload instance that has been emitted successfully from the
	// pipeline.
	Consume(context.Context, Payload) error
}
