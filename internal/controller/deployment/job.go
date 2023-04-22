package deployment

import (
	"context"
	"fmt"

	"github.com/infraboard/mpaas/apps/task"
	"github.com/infraboard/mpaas/common/format"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *Reconciler) HandleJobTask(ctx context.Context, obj appsv1.Deployment) error {
	// 获取日志对象
	l := log.FromContext(ctx)

	// 根据注解获取task id
	taskId := obj.Annotations[task.ANNOTATION_TASK_ID]
	if taskId == "" {
		return nil
	}
	l.Info(fmt.Sprintf("get mpaas job: %s", taskId))

	// 查询Task
	t, err := r.mpaas.JobTask().DescribeJobTask(ctx, task.NewDescribeJobTaskRequest(taskId))
	if err != nil {
		return fmt.Errorf("get task error, %s", err)
	}

	// 判断job当前状态
	updateReq := task.NewUpdateJobTaskStatusRequest(taskId)
	for _, cond := range obj.Status.Conditions {
		switch cond.Type {
		case appsv1.DeploymentReplicaFailure:
			if cond.Status == corev1.ConditionTrue {
				updateReq.Stage = task.STAGE_FAILED
				if cond.Message != "" {
					updateReq.Message = fmt.Sprintf("%s, %s", cond.Reason, cond.Message)
				}
			}
		case appsv1.DeploymentAvailable:
			if cond.Status == corev1.ConditionTrue {
				updateReq.Stage = task.STAGE_SUCCEEDED
				updateReq.Message = "执行成功"
			}
		case appsv1.DeploymentProgressing:
			if cond.Status == corev1.ConditionTrue {
				updateReq.Stage = task.STAGE_ACTIVE
			}
		}
	}

	if updateReq.Stage.Equal(task.STAGE_PENDDING) {
		return fmt.Errorf("task status is pendding, skip update")
	}

	// 比对状态, 状态没变化不更新
	if t.Status.Stage.Equal(updateReq.Stage) {
		return fmt.Errorf("task status is %s, not changed", updateReq.Stage)
	}

	// 状态变化更新
	updateReq.UpdateToken = t.Spec.UpdateToken
	updateReq.Detail = format.MustToYaml(obj)
	_, err = r.mpaas.JobTask().UpdateJobTaskStatus(ctx, updateReq)
	if err != nil {
		return fmt.Errorf("update jot task status error, %s", err)
	}
	return nil
}
