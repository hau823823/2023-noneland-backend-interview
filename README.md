# 2023 noneland backend interview

## 目錄結構

```plaintext
myapp/
|-- cmd/
|   `-- main.go
|-- configs/
|   |-- app.env
|   `-- config.go
|-- internal/
|   |-- api/
|   |   |-- api_test.go
|   |   |-- api.go
|   |   `-- server.go
|   |-- db/
|   |   |-- db_test.go
|   |   |-- db.go
|   |   `-- mock_db.go
|   |-- entity/
|   |   `-- entity.go
|   `-- pkg/
|       |-- cache.go
|       |-- client_test.go
|       |-- client.go
|       |-- mock_client.go
|-- go.mod
`-- go.sum
```
## API 說明
1. GetBalances
   
   獲取現貨和合約餘額
   
   **Endpoint**
   ```plaintext
   GET /api/v1/balances
   ```
   **Response**
   ```json
   {
     "spot_balance": {
       "asset": "USDT",
       "balance": 1000.0
     },
     "contract_balance": {
       "asset": "USDT",
       "balance": 500.0
     }
   }
   ```
   
2. FetchAndSaveSpotTransferRecords

   獲取第三方交易所現貨帳戶轉入轉出紀錄並保存到 DB

   **Endpoint**
   ```plaintext
   POST /api/v1/spot_transfer_records
   ```
   **Response**
   ```json
   {
     "message": "Transactions have been saved to the database."
   }
   ```

3. GetTransactions

   從DB中獲取交易紀錄
   
   **Endpoint**
   ```plaintext
   POST /api/v1/spot_transfer_records
   ```
   **Guery Parameters**
   - startTime (optional): 開始時間
   - endTime (optional): 結束時間

   **Response**
   ```json
   {
    "id": "1",
    "type": "IN",
    "amount": 0.1,
    "asset": "BNB",
    "status": "CONFIRMED",
    "timestamp": "2019-08-27T15:36:57Z",
    "txId": 5240372201
     },
     {
       "id": "2",
       "type": "OUT",
       "amount": 5.0,
       "asset": "USDT",
       "status": "CONFIRMED",
       "timestamp": "2019-08-27T13:13:56Z",
       "txId": 5239810406
     }
   ```
## 流程圖
```mermaid
graph TD
    A[Client Request] -->|GET /api/v1/balances| B[API Handler]
    B --> C[Check Cache]
    C -->|Cache Hit| D[Return Cached Response]
    C -->|Cache Miss| E[Call GetSpotBalance API]
    E --> F[Call GetContractBalance API]
    F --> G[Combine Results]
    G --> H[Save to Cache]
    H --> I[Return Response]
```
    
## 前情提要

假定我們自己的 `中心化` 交易所，已有錢包 app 完成初步現貨交易

現在正在處理『風險控制後台』管理介面開發，並提供 api 給『風險控制後台』前端開發人員串接

『風險控制後台』主要處理邏輯為，用戶下單現貨時，會在第三方開合約來做避險動作，因此需要有個後台讓風險控制人員觀看

我們會需要呼叫 `第三方XX交易所` 報價商後台所需要的參數

以下是 PM 提出的規格

## 基本規格

1. :white_check_mark:第一隻 api，同時回傳兩種餘額 
   1. 需要顯示 `XX交易所` 的 `現貨` 帳戶 USDT 餘額 
   1. 需要顯示 `XX交易所` 的 `合約` 帳戶 USDT 餘額
2. :white_check_mark:需要顯示 `現貨` 帳戶轉出轉入紀錄（`第三方XX交易所僅提供一個月內的交易紀錄查詢`）
   1. 根據法律遵循，我們應該保存 `6年` 內的所有交易紀錄

## 加分題

1. :white_check_mark:目前前端 app 的現貨買賣幣種報價使用的是 `現貨` 相關的 api，並且 api 存在呼叫限制，後台的呼叫不應該影響報價邏輯
   使用 cache 來處理
2. :white_check_mark:請撰寫可被測試的程式碼，或是直接附上測試程式
   預期你可以使用 mock 的方式來處理第三方，並通過 di 的方式注入，但使用其他方式也行
3. 架構也能調整，假設你覺得有更好的改法

## 待完善

1. :small_orange_diamond:cache 部分僅進行了簡單的實作，實務上可以搭配 redis 來完善
2. :small_orange_diamond:di 部分沒有使用 wire 框架實作

## 備註

可使用任何第三方工具、服務、套件來設計

補充一下，簡單來說目前後台會需要呼叫 `第三方XX交易` 來`取得`資料並於後台`顯示`

而 `第三方XX交易` 有 `現貨` 與 `合約` 兩個 api 的 endpoint

有任何題目上的問題都可以寫信發問


## XX交易所 api 文件

### 現貨帳戶 api

#### 前端 app 現貨下單的 api

api 權重：5

這邊為加分題，前端 app 使用的 endpoint 與下方 `取得現貨帳戶餘額` 一致

因此會受限於下方 `{{ url1 }}/exchangeInfo` 的限制

#### api 使用限制

- REQUEST_WEIGHT 單位時間請求權重之和的上限
- RAW_REQUESTS 單位時間請求次數上限

GET `{{ url1 }}/exchangeInfo`

api 權重：0

request:

Header: `Authorization: Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==`

response:

```json
{
  "timezone": "UTC",
  "serverTime": 1565246363776,
  "rateLimits": [
     {
        "rateLimitType": "REQUEST_WEIGHT",
        "interval": "MINUTE",
        "intervalNum": 1,
        "limit": 1200
     },
     {
        "rateLimitType": "RAW_REQUESTS",
        "interval": "MINUTE",
        "intervalNum": 5,
        "limit": 6100
     }
  ]
}
```

#### 取得現貨帳戶餘額

GET `{{ url1 }}/spot/balance`

api 權重：5

request:

Header: `Authorization: Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==`

response:

```json
{
  "free": "10.12345"
}
```

#### 現貨帳戶轉入轉出紀錄

GET `{{ url1 }}/spot/transfer/records`

api 權重：5

request 參數：

名稱 | 類型 | 是否必填 | 描述
----|------|--------|---------
startTime | LONG | NO |
endTime | LONG | NO |
endTime | LONG | NO |
current | LONG | NO | 當前回傳頁數，預設為 1
size | LONG | NO | 回傳筆數，預設 10，最大 100

Header: `Authorization: Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==`

response:

- status: PENDING (等待), CONFIRMED (成功), FAILED (失敗);

```json
{
   "rows": [
      {
         "amount": "0.10000000",
         "asset": "BNB",
         "status": "CONFIRMED",
         "timestamp": 1566898617,
         "txId": 5240372201,
         "type": "IN"
      },
      {
         "amount": "5.00000000",
         "asset": "USDT",
         "status": "CONFIRMED",
         "timestamp": 1566888436,
         "txId": 5239810406,
         "type": "OUT"
      },
      {
         "amount": "1.00000000",
         "asset": "EOS",
         "status": "CONFIRMED",
         "timestamp": 1566888403,
         "txId": 5239808703,
         "type": "IN"
      }
   ],
   "total": 3
}
```

### 合約帳戶 api

#### api 使用限制

- REQUEST_WEIGHT 單位時間請求權重之和的上限
- RAW_REQUESTS 單位時間請求次數上限

GET `{{ url2 }}/exchangeInfo`

api 權重：0

request:

Header: `Authorization: Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==`

response:

```json
{
  "timezone": "UTC",
  "serverTime": 1565246363776,
  "rateLimits": [
     {
        "rateLimitType": "REQUEST_WEIGHT",
        "interval": "MINUTE",
        "intervalNum": 1,
        "limit": 1200
     },
     {
        "rateLimitType": "RAW_REQUESTS",
        "interval": "MINUTE",
        "intervalNum": 5,
        "limit": 6100
     }
  ]
}
```


#### 取得合約帳戶餘額

GET `{{ url2 }}/futures/balance`

api 權重：5

request:

Header: `Authorization: Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==`

response:

```json
{
  "free": "10.12345"
}
```



