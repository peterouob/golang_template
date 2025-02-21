[//]: # (## 實作類似抖音功能，並且能處理高併發大流量任務)

[//]: # (- [x] Logger)

[//]: # (- [ ] mq)

[//]: # (- [ ] middleware)

[//]: # (  - [ ] 令牌桶)

[//]: # (  - [ ] 熔斷)

[//]: # (  - [ ] 降級)

[//]: # (  - [x] jwt)

[//]: # (- [ ] message queue)

[//]: # (  - [ ] kafka or rebbit mq)

[//]: # (- [x] grpc)

[//]: # (  - [x] 攔截器)

[//]: # (  - [x] gateway)

[//]: # (  - [x] pool)

[//]: # (- [x] etcd)

[//]: # (- [ ] prometheus)

[//]: # (  - [ ] Grafana)

[//]: # (- [ ] Jagger)

[//]: # (- [ ] docker)

[//]: # (- [ ] load balance)

[//]: # (- [ ] k8s)

[//]: # (  - [ ] minikub)

[//]: # (- [ ] Helm)

[//]: # (- [ ] aws in localstack)

[//]: # (  - [ ] ec2)

[//]: # (- [ ] CICD)

[//]: # (  - [x] github flow or jenkins)

[//]: # (  - [ ] argocd)

[//]: # ()
[//]: # (## Token 邏輯)

[//]: # (### 滿足black list條件,偵測當前流量|使用普羅米修斯|,假設在高流量狀態將id推送到mq上再批量處理blacklist邏輯,反之則直接寫入redis  )

## V1 版本 (本地個人)
1. 影片留言
2. 登入註冊
3. 上傳影片
4. 顯示個人資料中上傳的影片

## V2 版本 (預計接入多人)
1. websocket 多人通訊
2. 查看别人檔案
3. 訊息互動
4. Redis 儲存影片資料(有可能V1處理完)