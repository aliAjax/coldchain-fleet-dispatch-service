# Bug

并发读取仓储列表与并发更新同一仓储时发生 data race，读取路径未加读锁。

## 触发方式

```bash
go test -race ./internal/store -run '^TestStoreListUpdateNoRace$' -count=1
```

## 错误信息

```text
WARNING: DATA RACE
Read at 0x00c00009c450 by goroutine 7:
  github.com/example/coldchain-fleet-dispatch-service/internal/store.(*Store).ListShipments()
      internal/store/shipment_store.go:29
Previous write at 0x00c00009c450 by goroutine 8:
  runtime.mapaccessK()
  github.com/example/coldchain-fleet-dispatch-service/internal/store.(*Store).UpdateShipmentStatus()
      internal/store/shipment_store.go:44
```
