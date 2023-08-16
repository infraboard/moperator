/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pod

import (
	"context"
	"fmt"
	"strings"

	"github.com/infraboard/mpaas/apps/task"
	mpaas "github.com/infraboard/mpaas/clients/rpc"
	"github.com/infraboard/mpaas/common/format"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// PodReconciler reconciles a Pod object
type PodReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	mpaas *mpaas.ClientSet
}

//+kubebuilder:rbac:groups=mpaas.mdevcloud.com,resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=mpaas.mdevcloud.com,resources=pods/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=mpaas.mdevcloud.com,resources=pods/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Pod object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *PodReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// 获取日志对象
	l := log.FromContext(ctx)

	// TODO(user): your logic here

	// 1.通过名称获取Pod对象, 并打印
	var obj v1.Pod
	if err := r.Get(ctx, req.NamespacedName, &obj); err != nil {
		// 如果Pod对象不存在就删除该Pod
		if apierrors.IsNotFound(err) {
			l.Info(err.Error())
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if obj.Labels == nil {
		obj.Labels = map[string]string{}
	}

	// 根据注解获取task id, 更新Task状态
	taskId := obj.Labels["job-name"]
	if taskId != "" && strings.HasPrefix(taskId, "task-") {
		l.Info(fmt.Sprintf("get mpaas task: %s", taskId))

		// 查询Task, 获取更新Token
		t, err := r.mpaas.JobTask().DescribeJobTask(ctx, task.NewDescribeJobTaskRequest(taskId))
		if err != nil {
			l.Error(err, "get task error")
			return ctrl.Result{}, nil
		}

		// 当Pod处于Running时
		updateReq := task.NewUpdateJobTaskStatusRequest(taskId)
		updateReq.UpdateToken = t.Spec.UpdateToken
		updateReq.Stage = task.STAGE_ACTIVE
		updateReq.Message = fmt.Sprintf("Pod Status: %s", obj.Status.Phase)
		updateReq.Extension[task.EXTENSION_FOR_TASK_POD_DETAIL] = format.MustToYaml(obj)
		updateReq.Extension[task.EXTENSION_FOR_TASK_POD_STATUS] = string(obj.Status.Phase)
		_, err = r.mpaas.JobTask().UpdateJobTaskStatus(ctx, updateReq)
		if err != nil {
			l.Error(err, "update failed")
			return ctrl.Result{}, nil
		}

		l.Info("update success")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.mpaas = mpaas.C()
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Pod{}).
		Complete(r)
}
