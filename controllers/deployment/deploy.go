package deployment

import (
	"context"
	"fmt"

	"github.com/infraboard/mpaas/apps/deploy"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *Reconciler) HandleDeploy(ctx context.Context, obj appsv1.Deployment) error {
	// 获取日志对象
	l := log.FromContext(ctx)

	// 根据注解获取task id
	deployId := obj.Annotations[deploy.ANNOTATION_DEPLOY_ID]
	if deployId == "" {
		return nil
	}
	l.Info(fmt.Sprintf("get mpaas deploy: %s", deployId))

	return nil
}
