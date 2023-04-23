# 




## 本地调试

make run 会报如下错
```
1.646924212701068e+09 ERROR setup problem running manager {"error": "open /var/folders/67/375276sx6hv0nln1whwm5syh0000gq/T/k8s-webhook-server/serving-certs/tls.crt: no such file or directory"}
```

这是因为没有证书导致，[为WebHook手动签发证书](https://kubernetes.io/zh-cn/docs/tasks/tls/managing-tls-in-a-cluster/)

```sh
cat <<EOF | cfssl genkey - | cfssljson -bare server
{
  "hosts": [
    "172.22.109.184"
  ],
  "CN": "my-pod.my-namespace.pod.cluster.local",
  "key": {
    "algo": "ecdsa",
    "size": 256
  }
}
EOF
```



## 参考

+ [kubebuilder 官方文档](https://book.kubebuilder.io/introduction.html)
+ [Admission Webhook for Core Types](https://book.kubebuilder.io/reference/webhook-for-core-types.html)
+ [Webhook Example](https://github.com/kubernetes-sigs/controller-runtime/blob/main/examples/builtins/mutatingwebhook.go)
+ [给 K8s 中的 Operator 添加 Webhook 功能](https://blog.51cto.com/u_15773567/5671473)
+ [Kubebuilder 学习笔记之 Webhook Server](http://www.manongjc.com/detail/63-jkmbheqdeyizwra.html)