package deployment

import (
	"context"
	"fmt"

	"github.com/infraboard/mpaas/apps/deploy"
	"github.com/infraboard/mpaas/common/format"
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

	// 查询Deploy
	ins, err := r.mpaas.Deploy().DescribeDeployment(ctx, deploy.NewDescribeDeploymentRequest(deployId))
	if err != nil {
		return fmt.Errorf("get deploy error, %s", err)
	}

	// 更新Depoy
	updateReq := deploy.NewUpdateDeploymentStatusRequest(deployId)
	updateReq.UpdateToken = ins.Credential.Token
	obj.Kind = "Deployment"
	updateReq.UpdatedK8SConfig.WorkloadConfig = format.MustToYaml(obj)
	updateReq.UpdateBy = r.name
	_, err = r.mpaas.Deploy().UpdateDeploymentStatus(ctx, updateReq)
	if err != nil {
		return err
	}

	return nil
}
