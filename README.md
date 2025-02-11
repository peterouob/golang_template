## Go Gin template

### Use Todo List to show full fundamental template

- [x] Logger
- [ ] mq
- [ ] middleware
  - [ ] 令牌桶
  - [ ] 熔斷
  - [ ] 降級
  - [x] jwt
- [ ] message queue
  - [ ] kafka or rebbit mq
- [x] grpc
  - [x] 攔截器
  - [x] gateway
  - [x] pool
- [x] etcd
- [ ] prometheus
  - [ ] Grafana
- [ ] Jagger
- [ ] docker
- [ ] load balance
- [ ] k8s
  - [ ] minikub
- [ ] Helm
- [ ] aws in localstack
  - [ ] ec2
- [ ] CICD
  - [x] github flow or jenkins
  - [ ] argocd

[//]: # (Token 邏輯)
[//]: # (滿足black list條件,偵測當前流量|使用普羅米修斯|,假設在高流量狀態將id推送到mq上再批量處理blacklist邏輯,反之則直接寫入redis ) 