package pod

import (
	"context"
	"fmt"

	"github.com/infraboard/mpaas/apps/deploy"
	"github.com/infraboard/mpaas/provider/k8s/workload"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *PodWebHook) Default(ctx context.Context, obj runtime.Object) error {
	l := log.FromContext(ctx)
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return fmt.Errorf("expected a Pod but got a %T", obj)
	}

	l.Info("get pod", "pod", pod)

	if pod.Annotations == nil {
		pod.Annotations = map[string]string{}
	}

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
	l.Info(fmt.Sprintf("get mpaas deploy annotation: %s", deployId))

	// 查询Pod需要注入的Env变量
	queryEnv := deploy.NewQueryDeploymentInjectEnvRequest(deployId)
	set, err := r.mpaas.Deploy().QueryDeploymentInjectEnv(ctx, queryEnv)
	if err != nil {
		return fmt.Errorf("get deploy error, %s", err)
	}

	// 注入变量
	for i := range set.EnvGroups {
		group := set.EnvGroups[i]
		// 符合匹配条件的才进行注入
		if group.Enabled && group.IsLabelMatched(obj.Labels) {
			workload.InjectPodEnvVars(&obj.Spec, group.ToContainerEnvVars())
		}
	}

	return nil
}
