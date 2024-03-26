package pod

import (
	"context"
	"fmt"
	"strings"

	"github.com/infraboard/mflow/apps/task"
	"github.com/infraboard/mpaas/common/format"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *PodReconciler) HandleJobTask(ctx context.Context, obj v1.Pod) error {
	// 获取日志对象
	l := log.FromContext(ctx)

	taskId := obj.Labels["job-name"]
	if taskId != "" && strings.HasPrefix(taskId, "task-") {
		l.Info(fmt.Sprintf("get mflow task: %s, start sync", taskId))

		// 查询Task, 获取更新Token
		t, err := r.mflow.JobTask().DescribeJobTask(ctx, task.NewDescribeJobTaskRequest(taskId))
		if err != nil {
			return fmt.Errorf("get mflow task error, %s", err)
		}

		// 当Pod处于Running时
		updateReq := task.NewUpdateJobTaskStatusRequest(taskId)
		updateReq.UpdateToken = t.Spec.UpdateToken
		updateReq.Stage = task.STAGE_ACTIVE
		updateReq.Message = fmt.Sprintf("Pod %s Status: %s", obj.Name, obj.Status.Phase)
		updateReq.Extension[t.Status.GetOrNewPodKey(obj.Name)] = format.MustToYaml(obj)

		_, err = r.mflow.JobTask().UpdateJobTaskStatus(ctx, updateReq)
		if err != nil {
			return fmt.Errorf("update mflow task failed, %s", err)
		}
		l.Info(fmt.Sprintf("sync mflow task success: %s ", taskId))
	}
	return nil
}
