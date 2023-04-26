# moperator 开发流程




## 本地调试

make run 会报如下错
```
1.646924212701068e+09 ERROR setup problem running manager {"error": "open /var/folders/67/375276sx6hv0nln1whwm5syh0000gq/T/k8s-webhook-server/serving-certs/tls.crt: no such file or directory"}
```

### 准备证书

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

### 服务转发配置

[没有选择算符的 Service](https://kubernetes.io/zh-cn/docs/concepts/services-networking/service/#services-without-selectors)
```yaml
apiVersion: v1
kind: Service
metadata:
  name: moperator
spec:
  ports:
    - port: 443
      protocol: TCP
      targetPort: 9443
```

直接使用yml创建:
```sh
kubectl create -f docs/deploy/local/service.yml
kubectl get svc moperator -o yaml
```

手动添加 EndpointSlice
```yaml
apiVersion: discovery.k8s.io/v1
kind: EndpointSlice
metadata:
  name: moperator-1 # 按惯例将服务的名称用作 EndpointSlice 名称的前缀
  labels:
    # 你应设置 "kubernetes.io/service-name" 标签。
    # 设置其值以匹配服务的名称
    kubernetes.io/service-name: moperator
addressType: IPv4
ports:
  - name: '' # 留空，因为 port 9376 未被 IANA 分配为已注册端口
    appProtocol: http
    protocol: TCP
    port: 9443
endpoints:
  - addresses:
      - "10.4.5.6" 修改为你本地ip
```

直接使用yml创建:
```sh
kubectl create -f docs/deploy/local/endpoint.yml
kubectl get endpointslice moperator-1 -o yaml
```

验证转发是否成功, 测试是否成正常转发到本地服务
```sh
curl https://cluster_ip:443/mutate--v1-pod
```

### Webhook配置

```sh
kubectl create -f docs/deploy/local/webhook_config.yaml
kubectl get mutatingWebhookConfiguration
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