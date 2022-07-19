package hw06pipelineexecution

import (
	"fmt"
	"time"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in

	fmt.Println("Before execute: ", time.Now().Format("15:04:05.000000"))
	for _, s := range stages {
		fmt.Println("Run stage: ", time.Now().Format("15:04:05.000000"))
		out = s(out)
	}
	return out
}
