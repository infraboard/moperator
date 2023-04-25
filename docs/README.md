# moperator 开发流程




## 本地调试

make run 会报如下错
```
1.646924212701068e+09 ERROR setup problem running manager {"error": "open /var/folders/67/375276sx6hv0nln1whwm5syh0000gq/T/k8s-webhook-server/serving-certs/tls.crt: no such file or directory"}
```

下面是一个自动创建客户端证书的脚本: [create-signed-cert.sh](https://raw.githubusercontent.com/kubernetes-sigs/windows-gmsa/master/admission-webhook/deploy/create-signed-cert.sh)
```sh
Generates certificate suitable for use with the GMSA webhook service.

This script uses k8s' CertificateSigningRequest API to a generate a
certificate signed by k8s CA suitable for use with the GMSA webhook
service. This requires permissions to create and approve CSR. See
https://kubernetes.io/docs/tasks/tls/managing-tls-in-a-cluster for
detailed explantion and additional instructions.

usage: ./create-signed-cert.sh --service SERVICE_NAME --namespace NAMESPACE_NAME --certs-dir PATH/TO/CERTS/DIR [--dry-run] [--overwrite]

If --dry-run is set, the script echoes what command it would perform
to stdout without actually affecting the k8s cluster.
If the files this script generates already exist and --overwrite is
not set, it will not regenerate the files.
```

## 线上部署



## 参考

+ [kubebuilder 官方文档](https://book.kubebuilder.io/introduction.html)
+ [Admission Webhook for Core Types](https://book.kubebuilder.io/reference/webhook-for-core-types.html)
+ [Webhook Example](https://github.com/kubernetes-sigs/controller-runtime/blob/main/examples/builtins/mutatingwebhook.go)
+ [给 K8s 中的 Operator 添加 Webhook 功能](https://blog.51cto.com/u_15773567/5671473)
+ [Kubebuilder 学习笔记之 Webhook Server](http://www.manongjc.com/detail/63-jkmbheqdeyizwra.html)
+ [windows-gmsa](https://github.com/kubernetes-sigs/windows-gmsa)
+ [k8s认证与鉴权](https://zhuanlan.zhihu.com/p/572600485)