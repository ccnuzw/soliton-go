package transaction

import (
	"context"
	"fmt"
)

// Step represents a single step in a Saga transaction.
type Step struct {
	Name         string
	Action       func(ctx context.Context) error
	Compensation func(ctx context.Context) error
}

// SagaOrchestrator manages the execution of a Saga.
type SagaOrchestrator struct {
	steps []Step
}

// NewSaga creates a new SagaOrchestrator.
func NewSaga() *SagaOrchestrator {
	return &SagaOrchestrator{}
}

// AddStep adds a step to the Saga.
func (s *SagaOrchestrator) AddStep(name string, action func(ctx context.Context) error, compensation func(ctx context.Context) error) *SagaOrchestrator {
	s.steps = append(s.steps, Step{
		Name:         name,
		Action:       action,
		Compensation: compensation,
	})
	return s
}

// Execute runs the Saga steps in order.
// If any step fails, it executes the compensations of previously successful steps in reverse order.
func (s *SagaOrchestrator) Execute(ctx context.Context) error {
	var completedSteps []Step

	for _, step := range s.steps {
		if err := step.Action(ctx); err != nil {
			// Rollback
			fmt.Printf("Saga step '%s' failed: %v. Initiating rollback...\n", step.Name, err)
			s.rollback(ctx, completedSteps)
			return fmt.Errorf("saga failed at step '%s': %w", step.Name, err)
		}
		completedSteps = append(completedSteps, step)
	}

	return nil
}

func (s *SagaOrchestrator) rollback(ctx context.Context, completedSteps []Step) {
	// Execute compensations in reverse order
	for i := len(completedSteps) - 1; i >= 0; i-- {
		step := completedSteps[i]
		if step.Compensation != nil {
			if err := step.Compensation(ctx); err != nil {
				// We can't do much here other than logging.
				// In a real system, this might require manual intervention or a retry queue.
				fmt.Printf("Failed to compensate step '%s': %v\n", step.Name, err)
			} else {
				fmt.Printf("Compensated step '%s'\n", step.Name)
			}
		}
	}
}
