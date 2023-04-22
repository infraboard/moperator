package pod

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type PodWebHook struct {
}

func (r *PodWebHook) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&corev1.Pod{}).
		WithDefaulter(&PodWebHook{}).
		Complete()
}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *PodWebHook) Default(ctx context.Context, obj runtime.Object) error {
	log := logf.FromContext(ctx)
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return fmt.Errorf("expected a Pod but got a %T", obj)
	}

	if pod.Annotations == nil {
		pod.Annotations = map[string]string{}
	}
	pod.Annotations["example-mutating-admission-webhook"] = "foo"
	log.Info("Annotated Pod")
	return nil
}
