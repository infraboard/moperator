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

package job

import (
	"context"
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/infraboard/mpaas/apps/task"
	mpaas "github.com/infraboard/mpaas/clients/rpc"
	"github.com/infraboard/mpaas/common/format"
)

// Reconciler reconciles a Job object
type Reconciler struct {
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
func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// 获取日志对象
	l := log.FromContext(ctx)

	// 通过名称获取Pod对象, 并打印
	var obj batchv1.Job
	if err := r.Get(ctx, req.NamespacedName, &obj); err != nil {
		// 如果Pod对象不存在就删除该Pod
		if apierrors.IsNotFound(err) {
			l.Info(err.Error())
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 根据注解获取task id
	taskId := obj.Annotations[task.ANNOTATION_TASK_ID]
	if taskId == "" {
		return ctrl.Result{}, nil
	}
	l.Info(fmt.Sprintf("get mpaas job: %s", taskId))

	// 查询Task
	t, err := r.mpaas.JobTask().DescribeJobTask(ctx, task.NewDescribeJobTaskRequest(taskId))
	if err != nil {
		l.Error(err, "get task error")
		return ctrl.Result{}, nil
	}

	// 判断job当前状态
	updateReq := task.NewUpdateJobTaskStatusRequest(taskId)
	for _, cond := range obj.Status.Conditions {
		switch cond.Type {
		case batchv1.JobFailed, batchv1.JobFailureTarget:
			if cond.Status == corev1.ConditionTrue {
				updateReq.Stage = task.STAGE_FAILED
				if cond.Message != "" {
					updateReq.Message = fmt.Sprintf("%s, %s", cond.Reason, cond.Message)
				}
			}
		case batchv1.JobComplete:
			if cond.Status == corev1.ConditionTrue && !updateReq.Stage.Equal(task.STAGE_FAILED) {
				updateReq.Stage = task.STAGE_SUCCEEDED
				updateReq.Message = "执行成功"
			}
		case batchv1.JobSuspended:
			if cond.Status == corev1.ConditionTrue {
				updateReq.Stage = task.STAGE_CREATING
			}
		}
	}

	// 状态变化更新, 运行一种状态多次更新, 获取更多细节信息
	updateReq.UpdateToken = t.Spec.UpdateToken
	updateReq.Detail = format.MustToYaml(obj)
	_, err = r.mpaas.JobTask().UpdateJobTaskStatus(ctx, updateReq)
	if err != nil {
		l.Error(err, "update failed")
		return ctrl.Result{}, nil
	}

	l.Info("update success")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.mpaas = mpaas.C()
	return ctrl.NewControllerManagedBy(mgr).
		For(&batchv1.Job{}).
		Complete(r)
}
