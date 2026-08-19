# Bug

批量派车走一遍后，三个待派托运单都还停在 `pending`，没有被分配。

## 触发方式

```bash
go test -race ./internal/worker -run '^TestDispatchPendingAssignsAllShards$' -count=1
```

## 错误信息

```text
--- FAIL: TestDispatchPendingAssignsAllShards
    dispatcher_test.go:78: shipment ship-1 not assigned: <nil>
```
