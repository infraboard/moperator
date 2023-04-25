package pod

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// validate admits a pod if a specific annotation exists.
func (v *PodWebHook) validate(ctx context.Context, obj runtime.Object) error {
	log := log.FromContext(ctx)
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return fmt.Errorf("expected a Pod but got a %T", obj)
	}

	log.Info("Validating Pod")
	key := "example-mutating-admission-webhook"
	anno, found := pod.Annotations[key]
	if !found {
		return fmt.Errorf("missing annotation %s", key)
	}
	if anno != "foo" {
		return fmt.Errorf("annotation %s did not have value %q", key, "foo")
	}

	return nil
}

func (v *PodWebHook) ValidateCreate(ctx context.Context, obj runtime.Object) error {
	return v.validate(ctx, obj)
}

func (v *PodWebHook) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) error {
	return v.validate(ctx, newObj)
}

func (v *PodWebHook) ValidateDelete(ctx context.Context, obj runtime.Object) error {
	return v.validate(ctx, obj)
}
