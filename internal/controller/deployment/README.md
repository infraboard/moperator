# deployment 扩展控制器

如果 job里面的执行逻辑是异步的, 比如执行 kubectl apply, 那么该task的状态并不能表示任务真正的状态

需要处理的Annotation:
+ task.mpaas.inforboar.io/id