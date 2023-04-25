package pod

import (
	mpaas "github.com/infraboard/mpaas/client/rpc"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

// +kubebuilder:webhook:path=/mutate-v1-pod,mutating=true,failurePolicy=fail,groups="",resources=pods,verbs=create;update,versions=v1,name=mpod.kb.io
type PodWebHook struct {
	mpaas *mpaas.ClientSet
}

func (r *PodWebHook) SetupWebhookWithManager(mgr ctrl.Manager) error {
	r.mpaas = mpaas.C()
	return ctrl.NewWebhookManagedBy(mgr).
		For(&corev1.Pod{}).
		WithDefaulter(&PodWebHook{}).
		WithValidator(&PodWebHook{}).
		Complete()
}
