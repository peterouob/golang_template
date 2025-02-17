## 統一開啟prometheus然後report給Grafana
- 重構prometheus register讓他能吃不同的tag and help
- 依照服務需求選擇histogram or counterVec
- 想辦法拿到prometheus 的 response
    - 主要是拿到流量即可
- 將資料喂給Grafana

## Token 編寫成rpc service
- interceptors 寫在Token Valid server中只有要調用的去掉用Valid
- Redis 儲存Refresh Token
- Mq 傳送資料 or 直接從Redis
