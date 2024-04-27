#

版定系统角色
```sh
# 创建只读

# 创建帐号
kubectl create serviceaccount mcloud

# 绑定
kubectl create clusterrolebinding mcloud-cluster-admin --clusterrole=cluster-admin --serviceaccount=default:mcloud
```