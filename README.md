# 基于 k8s operator 的部署同步

选择开发的工具: https://book.kubebuilder.io/


## 准备 

```sh
# download kubebuilder and install locally.
curl -L -o kubebuilder "https://go.kubebuilder.io/dl/latest/$(go env GOOS)/$(go env GOARCH)"
chmod +x kubebuilder && sudo mv kubebuilder /usr/local/bin/
```

```sh
echo "source <(kubebuilder completion zsh)" >> ~/.zshrc
source ~/.zshrc
```


## 创建项目

```sh
kubebuilder init --domain devops.mcloud --repo github.com/infraboard/devops/moperator
```

## 创建资源

```sh
kubebuilder create api --group devops --version v1beta1 --kind Mcloud
```