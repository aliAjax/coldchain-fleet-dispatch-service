# Bug

查询不存在的交接清单时，错误链断裂导致接口返回 500 而不是 404。

## 触发方式

```bash
go test ./internal/api -run '^TestChecklistNotFoundReturns404$' -count=1
```

## 错误信息

```text
--- FAIL: TestChecklistNotFoundReturns404
    ops_handlers_test.go:31: expected 404, got 500
```
