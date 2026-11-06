# 2026-11-06 · sync.Once 與 OnceFunc / OnceValue / OnceValues

## 為什麼要學
「只初始化一次」到處都是：讀設定檔、建立 DB 連線池、編譯 regexp、載入憑證、lazy 建立 singleton client。
程式一啟動就做會拖慢啟動、而且不一定用得到；但「第一次用到才做」在多 goroutine 下就會有 race。
`sync.Once` 是標準庫給的答案；Go 1.21 加入的 `OnceFunc` / `OnceValue` / `OnceValues` 讓寫法更短、也修正了 `Once` 在 panic 時的坑。
接續上次的 `sync_cond`，這是 sync 套件系列的下一個。

## 核心觀念
1. `once.Do(f)`：不論多少 goroutine 同時呼叫，`f` 只執行一次；**而且所有呼叫者在 `Do` 返回時，`f` 都已經執行完畢**。
2. runtime / 標準庫實作（go1.26 `src/sync/once.go`）：
   - `Once` 只有兩個欄位：`done atomic.Bool` 和 `m Mutex`。
   - Fast path：`if !o.done.Load() { o.doSlow(f) }`，這段小到會被 inline 到呼叫端，初始化後每次呼叫只剩一次 atomic load。
   - Slow path：拿 mutex → 再檢查一次 `done`（double-checked locking）→ 執行 `f` → **`defer done.Store(true)`**。
   - 為什麼不用 `CompareAndSwap(false, true)` 再執行 `f`？因為搶輸 CAS 的 goroutine 會立刻返回，但 `f` 還沒跑完，它就會讀到半初始化的資料。原始碼註解特別寫了這點。
   - 因為 `Store(true)` 放在 defer，**`f` panic 也會被當成「做完了」**。
3. Go 1.21 新函式（https://pkg.go.dev/sync#OnceFunc）：
   - `OnceFunc(f) func()`：回傳一個只會執行 f 一次的函式。
   - `OnceValue(f func() T) func() T`：只算一次並快取回傳值。
   - `OnceValues(f func() (T1, T2)) func() (T1, T2)`：常用在 `(value, error)`。
   - 實作上是把 `f`、`Once`、`valid`、panic 值 `p`、結果包進同一個 struct（只有一次 heap 配置），執行後把 `f` 設為 nil 讓 closure 可被 GC。
   - **若 f panic，之後每次呼叫都會用同一個值再 panic**，不會像 `Once.Do` 一樣默默回傳零值。

## 完整範例
見本目錄 `main.go`（已在 go1.26.3 linux/amd64 實際 `go run`、`go run -race` 執行，並 `GOOS=windows go build` 確認可編譯）。
最低需求 Go 1.22（用了 `for range 100`）；`OnceValue` 系列需 1.21。

預期輸出：
```
== 1. sync.Once ==
100 個 goroutine 呼叫 GetConfig，實際載入次數 = 1
config = {DSN:file:app.db MaxConn:10}
== 2. OnceValue / OnceValues ==
  (計算平方表中…只會印一次)
table[3] = 9
table[9] = 81
  (讀取 secret 中…只會印一次)
第 1 次 loadSecret: "", err=secret file not found
第 2 次 loadSecret: "", err=secret file not found
== 3. panic 時的行為 ==
  Once.Do 第 1 次 -> panic: init failed
  Once.Do 第 2 次 -> 正常返回（f 沒有再被執行）
  OnceFunc 第 1 次 -> panic: init failed
  OnceFunc 第 2 次 -> panic: init failed
```

## 常見陷阱
1. **錯誤也會被永久快取**：`OnceValues` 第一次失敗（例如網路暫時斷線），之後永遠回傳同一個 error，不會重試。需要「失敗可重試」時要自己用 `Mutex` + 旗標實作。
2. **`Once.Do` 遇到 panic 會變成「已完成但沒初始化」**：第二次呼叫直接返回，拿到 nil 指標，錯誤在很遠的地方才爆出來。新程式建議優先用 `OnceValue(s)`。
3. **在 `f` 裡再呼叫同一個 `once.Do` 會 deadlock**：內部 `Mutex` 不可重入。另外 `Once` 不能被複製（含 `noCopy`，`go vet` 的 copylocks 會抓到），放在 struct 裡時記得用指標接收者。

## 練習題（不給答案）
寫一個 `RetryOnce[T any](f func() (T, error)) func() (T, error)`：
- 成功後的行為與 `OnceValues` 相同（只算一次、之後直接回傳快取值，且 fast path 不拿鎖）。
- 失敗時**不快取**，下一次呼叫會重新執行 `f`。
- 多個 goroutine 同時呼叫時，同一時間只有一個在執行 `f`。
- 用 `go run -race` 驗證，並寫一個「前兩次失敗、第三次成功」的 f 來測試。
提示方向：想想 `sync.Once` 原始碼裡 `done` 為什麼要在 `f` 之後才 `Store`。

## 建議目錄名稱
`sync_once`（與既有的 `sync_cond`、`sync_map` 命名一致，放在 README 的 🧵 Concurrency Patterns）

## 參考
- https://pkg.go.dev/sync#Once
- https://pkg.go.dev/sync#OnceFunc
- https://go.dev/doc/go1.21#sync
- https://go.dev/ref/mem#once （Go Memory Model 中 Once 的 happens-before 保證）
