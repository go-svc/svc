# 範例

* `grpc`：基本的 gRPC 客戶端、伺服端範例。
* `grpc-database`：gRPC 伺服端與和資料庫互動的範例。
* `grpc-lb`：搭載基於 Consul 的服務探索、負載平衡 gRPC 範例。
* `grpc-fixed-lb`：最基本、簡單，無需依賴其他額外服務的負載平衡 gRPC 範例。
* `grpc-ver-lb`：基於 Consul 並帶有版本控制的服務探索 gRPC 範例。
* `nsq-pubsub`：透過 NSQ 分佈式訊息佇列系統在不同服務之間異步傳遞訊息。
* `pb-pubsub`：以 NSQ 在不同服務之間異步傳遞 Protobuf 訊息並編碼、解碼。
* `es-pubsub`：基於 EventStore 事件存儲中心廣播還有處理事件。
* `grpc-opentracing`：基於負載平衡還有 OpenTracing 分布式追蹤系統的 gRPC 範例。
* `grpc-metadata`：在 gRPC 內傳遞額外的中繼資料供追蹤、紀錄用途。
* `grpc-jwt`：在 gRPC 內透過中繼資料傳遞、並簽發與解析 JSON Web Token。