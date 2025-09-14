package controller

import (
	"context"
	"time"

	"github.com/infraboard/devops/api/mpaas/cluster"
	"github.com/infraboard/devops/api/mpaas/deploy"
	"github.com/infraboard/mcube/v2/tools/pretty"
	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type DeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *DeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	// TODO(user): your logic here
	// 1. 获取Deployment对象
	var obj appsv1.Deployment
	if err := r.Get(ctx, req.NamespacedName, &obj); err != nil {
		// 如果Deployment对象不存在就删除该Deployment
		if apierrors.IsNotFound(err) {
			l.Info(err.Error())
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if obj.Annotations == nil {
		l.Info("deployment not annotations")
		return ctrl.Result{}, nil
	}

	clusterId := obj.Annotations[deploy.ANNOTATION_CLUSTER_ID]
	if clusterId == "" {
		l.Info("deployment not annotations cluster id")
		return ctrl.Result{}, nil
	}

	// 2. 把这个deployment更新到mpaas服务
	l.Info("sync deployment to mpaas", "deployment", obj.Name, "cluster_id", clusterId)
	updateReq := cluster.NewUpdateDeployStatusRequest(clusterId)
	updateReq.SetSyncAt(time.Now())
	updateReq.SetUpdateAt(time.Now())
	updateReq.Replicas = obj.Status.Replicas
	updateReq.Detail = pretty.MustToYaml(obj)
	updateReq.Status, updateReq.Message = cluster.ComputeK8sDeploymentStatus(&obj)

	// TODO 这里可以把更多的状态更新过去
	_, err := cluster.NewClusterServiceClient().RPCUpdateDeployStatus(ctx, updateReq)
	if err != nil {
		l.Error(err, "update mcloud deploy status error")
	}

	return ctrl.Result{}, nil
}

func (r *DeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: 5,
		}).
		Named("devops_deployment_syncer").
		Complete(r)
}
