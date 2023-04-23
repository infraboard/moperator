package pod

import (
	"context"
	"fmt"

	"github.com/infraboard/mpaas/apps/deploy"
	mpaas "github.com/infraboard/mpaas/client/rpc"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type PodWebHook struct {
	mpaas *mpaas.ClientSet
}

func (r *PodWebHook) SetupWebhookWithManager(mgr ctrl.Manager) error {
	r.mpaas = mpaas.C()
	return ctrl.NewWebhookManagedBy(mgr).
		For(&corev1.Pod{}).
		WithDefaulter(&PodWebHook{}).
		Complete()
}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *PodWebHook) Default(ctx context.Context, obj runtime.Object) error {
	l := log.FromContext(ctx)
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return fmt.Errorf("expected a Pod but got a %T", obj)
	}

	if pod.Annotations == nil {
		pod.Annotations = map[string]string{}
	}
	pod.Annotations["example-mutating-admission-webhook"] = "foo"
	l.Info("Annotated Pod")

	// 处理deploy注解: deploy.mpaas.inforboard.io/id
	err := r.HandlePodEnvInject(ctx, pod)
	if err != nil {
		l.Error(err, "hanle deploy error")
	}

	return nil
}

func (r *PodWebHook) HandlePodEnvInject(ctx context.Context, obj *corev1.Pod) error {
	// 获取日志对象
	l := log.FromContext(ctx)

	// 根据注解获取task id
	deployId := obj.Annotations[deploy.ANNOTATION_DEPLOY_ID]
	if deployId == "" {
		return nil
	}
	l.Info(fmt.Sprintf("get mpaas deploy: %s", deployId))

	// 查询Pod需要注入的Env变量
	queryEnv := deploy.NewQueryDeploymentInjectEnvRequest(deployId)
	set, err := r.mpaas.Deploy().QueryDeploymentInjectEnv(ctx, queryEnv)
	if err != nil {
		return fmt.Errorf("get deploy error, %s", err)
	}

	// 注入变量
	for i := range set.EnvGroups {
		group := set.EnvGroups[i]
		if group.Enabled {

		}
	}

	return nil
}
