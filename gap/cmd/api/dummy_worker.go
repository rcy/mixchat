package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/riverqueue/river"
)

type DummyArgs struct {
	GoogleProfileID uuid.UUID
	SyncToken       string
}

func (DummyArgs) Kind() string { return "Dummy" }

type DummyWorker struct {
	river.WorkerDefaults[DummyArgs]
}

func (w *DummyWorker) Work(ctx context.Context, job *river.Job[DummyArgs]) error {
	fmt.Println("WORK")
	panic("die")
}
