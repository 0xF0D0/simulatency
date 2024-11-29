package leader

import "context"

type LeaderElectableRunloop interface {
	CleanUp()
	Run(ctx context.Context)
}


