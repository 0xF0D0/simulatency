package leader

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/0xf0d0/simulatency-controller/client/k8s"
	"github.com/0xf0d0/simulatency-controller/client/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/klog/v2"
)

const (
	leaseLockName = "simulatency-lock"
)

func InsertElectableRunloop(leaseLockNamespace string, loop LeaderElectableRunloop) error {
	id := uuid.NewUUIDGenerator().GenerateString()
	ctx, cancel := context.WithCancel(context.Background())
	k8sClient, err := k8s.NewClient("")

	if err != nil {
		klog.Fatal(err)
	}
	
	defer cancel()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		klog.Info("Received termination, signaling shutdown")
		cancel()
	}()

	lock := &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      leaseLockName,
			Namespace: leaseLockNamespace,
		},
		Client: k8sClient.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: id,
		},
	}

	// start the leader election code loop
	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock:            lock,
		LeaseDuration:   60 * time.Second,
		RenewDeadline:   15 * time.Second,
		RetryPeriod:     5 * time.Second,
		Callbacks:       leaderelection.LeaderCallbacks{
			OnStartedLeading: loop.Run,
			OnStoppedLeading: func() {
			},
			OnNewLeader: func(identity string) {
			},
		},
		ReleaseOnCancel: true,
		Name:            leaseLockName,
	})

	//Probably wouldn't reach here
	return nil
}

