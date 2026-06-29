# go-perry 專案主題總覽

Go 學習筆記與主題式範例集合。每個資料夾通常代表一個技術主題、實作練習，或單一機制的示範專案。

## 根目錄各資料夾說明

| 資料夾 | 主題說明 |
| --- | --- |
| alogrithm | 演算法與資料結構練習，底下再分成 design、search、sort、structure 等子主題。 |
| async | 非同步處理範例，包含 producer、consumer、tasks 等背景任務處理概念。 |
| background_faktory | 使用 Faktory 建立背景工作佇列的範例。 |
| background_gocraft | 使用 gocraft/work 之類的工作佇列框架處理背景任務。 |
| benchmark | Go 效能測試、CPU profiling、memory profiling 的示範。 |
| boomfilter | Bloom Filter 實作範例，用來示範機率型集合查詢。資料夾名稱應是 Bloom Filter 的變體拼法。 |
| breakcircut | Circuit Breaker 熔斷機制範例。資料夾名稱是 break circuit / circuit breaker 的變體拼法。 |
| build_tag | Go build tags 的使用方式，例如不同環境或不同作業系統的編譯條件。 |
| cache | 快取機制範例，例如簡易快取與 cached HTTP client。 |
| channel | Go channel 基本用法與多個小範例。 |
| channel_buffered | buffered channel 的使用示範。 |
| channel_unbuffered | unbuffered channel 的使用示範。 |
| context | Go context 的使用方式，例如取消、逾時、值傳遞等。 |
| convert | 型別轉換、資料格式轉換相關範例。 |
| decorator | Decorator 設計模式示範。 |
| dependency | 依賴注入或模組依賴設計的範例。 |
| errgroup | 使用 errgroup 管理 goroutine 與錯誤彙整。 |
| error, handling | Go 錯誤處理模式與基本實作示範。 |
| function, slice, array | Go 函式、slice、array 的語法與行為練習。 |
| generator | 生成器模式或用 goroutine 模擬 generator 的範例。 |
| gocron | 使用 gocron 做排程工作。 |
| gorilla_websocket | 使用 gorilla/websocket 建立 WebSocket 服務與前端測試頁。 |
| gorm | GORM ORM 範例集合，包含 migration、model 與多個 example。 |
| goroutine_confinement | goroutine confinement 模式，強調共享資料的封裝與限制存取。 |
| goroutine_mutex | 以 mutex 保護共享資源的 goroutine 同步範例。 |
| graphql | GraphQL 伺服器範例，含 schema 與 graph 設定。 |
| grpc | gRPC client/server 與 proto 定義範例。 |
| hard | 偏實戰型的小題目或進階練習；目前內容看起來是 worker pool 實作。 |
| ipfs | IPFS 或去中心化儲存相關範例。 |
| iter | iterator 模式或迭代器風格程式設計練習。 |
| jwt | JWT 驗證、授權與簡單 Web 範例。 |
| kafka | Kafka producer / worker 訊息處理範例。 |
| logger | 記錄日誌與輸出到檔案的範例。 |
| lsm | LSM Tree 儲存引擎概念實作，包含 wal、mem、ssl、lsm 等子模組。 |
| memcache | Memcache 的使用範例，可能包含簡單 API 或資料存取示範。 |
| module | Go module 使用方式與模組拆分練習。 |
| multipartload | multipart 檔案上傳的 client/server 範例。 |
| orDone&T | Go pipeline 常見模式示範，包含 orDone、tee 等 channel 組合技巧。 |
| p2p | P2P 網路相關範例，可能包含 launchpad tutorial 類型練習。 |
| paseto | PASETO token 驗證與安全性相關範例。 |
| pin_thread | 將 goroutine 綁定 OS thread 或執行緒控制的實驗。 |
| pipeline | pipeline 串流處理與併發資料流範例。 |
| pool | 資源池或 worker pool 類型的基本範例。 |
| protobuf | Protocol Buffers 序列化與程式碼生成範例。 |
| rabbitmq | RabbitMQ 訊息佇列與 consumer / producer 範例。 |
| ratelimit | 限流機制實作或套件使用示範。 |
| redis_lock | Redis 分散式鎖範例。 |
| reflect | Go reflect 套件與反射操作練習。 |
| rss | RSS feed 產生或解析範例。 |
| runtime | Go runtime 相關觀念或執行期行為觀察。 |
| schnorr | Schnorr signature 或相關密碼學演算法範例。 |
| singleflight | singleflight 去重請求機制範例。 |
| sort | 排序演算法或排序操作範例。 |
| struct, method | struct 與 method 定義方式的基礎練習。 |
| sync_map | sync.Map 的使用範例。 |
| test | Go 測試主題範例，底下有多個 func_xx 子資料夾。 |
| unsafe | unsafe 套件與底層記憶體操作實驗。 |
| variable, loop | 變數宣告與迴圈語法練習。 |
| worker_pool | worker pool 模式範例。 |
| workerpool | 另一組 worker pool 或工作池實作練習，可能與 worker_pool 分別示範不同寫法。 |

## Commit 自動化

- 已加入可版本控制的 Git hooks，安裝方式：執行 `pwsh -File scripts/install-git-hooks.ps1`。
- `prepare-commit-msg` 會在你沒有手寫 commit message 時，自動依 staged 的頂層資料夾產生預設訊息。
- `pre-commit` 會檢查 staged 檔案所屬的新頂層資料夾；如果 `info.md` 還沒有該條目，就自動補上一列 `TODO` 說明並重新 stage。
- 這表示你新增新的主題資料夾後，commit 時會自動把該主題補進這份文件，但主題說明文字仍建議後續手動補完整。

| falseshare | TODO: 補上此資料夾的專案主題說明。 |

| barrier | TODO: 補上此資料夾的專案主題說明。 |

| escape | TODO: 補上此資料夾的專案主題說明。 |
