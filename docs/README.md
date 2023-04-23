# moperator 开发流程




## 本地调试

make run 会报如下错
```
1.646924212701068e+09 ERROR setup problem running manager {"error": "open /var/folders/67/375276sx6hv0nln1whwm5syh0000gq/T/k8s-webhook-server/serving-certs/tls.crt: no such file or directory"}
```

### 准备本地证书

这是因为没有证书导致，[为WebHook手动签发证书](https://kubernetes.io/zh-cn/docs/tasks/tls/managing-tls-in-a-cluster/)

#### 签发请求

1. 创建证书签名请求
```sh
cat <<EOF | cfssl genkey - | cfssljson -bare server
{
  "hosts": [
    "172.22.109.184"
  ],
  "CN": "my-pod.default.pod.cluster.local",
  "key": {
    "algo": "ecdsa",
    "size": 256
  }
}
EOF
```

2. 创建证书签名请求（CSR）对象发送到 Kubernetes API
```sh
cat <<EOF | kubectl apply -f -
apiVersion: certificates.k8s.io/v1
kind: CertificateSigningRequest
metadata:
  name: moperator.default
spec:
  request: $(cat server.csr | base64 | tr -d '\n')
  signerName: example.com/serving
  usages:
  - digital signature
  - key encipherment
  - server auth
EOF
```

确认签发请求
```sh
kubectl describe csr moperator.default
```

#### 批准签发

1. 批准证书签名请求
```sh
kubectl certificate approve moperator.default
kubectl get csr
```


#### 签发证书

扮演CA 签发证书
1. 初始化CA
```sh
cat <<EOF | cfssl gencert -initca - | cfssljson -bare ca
{
  "CN": "My Example Signer",
  "key": {
    "algo": "rsa",
    "size": 2048
  }
}
EOF
```

2. 颁发证书, 办法请求: server-signing-config.json
```json
{
    "signing": {
        "default": {
            "usages": [
                "digital signature",
                "key encipherment",
                "server auth"
            ],
            "expiry": "876000h",
            "ca_constraint": {
                "is_ca": false
            }
        }
    }
}
```

执行颁发:
```sh
kubectl get csr moperator.default -o jsonpath='{.spec.request}' | \
  base64 --decode | \
  cfssl sign -ca ca/ca.pem -ca-key ca/ca-key.pem -config server-signing-config.json - | \
  cfssljson -bare ca-signed-server
```

这会生成一个签名的服务证书文件: ca-signed-server.pem
```sh
kubectl get csr moperator.default -o json | \
  jq '.status.certificate = "'$(base64 ca-signed-server.pem | tr -d '\n')'"' | \
  kubectl replace --raw /apis/certificates.k8s.io/v1/certificatesigningrequests/moperator.default/status -f -
```

####  上传签名证书

## 线上部署



## 参考

+ [kubebuilder 官方文档](https://book.kubebuilder.io/introduction.html)
+ [Admission Webhook for Core Types](https://book.kubebuilder.io/reference/webhook-for-core-types.html)
+ [Webhook Example](https://github.com/kubernetes-sigs/controller-runtime/blob/main/examples/builtins/mutatingwebhook.go)
+ [给 K8s 中的 Operator 添加 Webhook 功能](https://blog.51cto.com/u_15773567/5671473)
+ [Kubebuilder 学习笔记之 Webhook Server](http://www.manongjc.com/detail/63-jkmbheqdeyizwra.html)