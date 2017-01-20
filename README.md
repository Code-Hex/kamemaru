# kamemaru project
Imgur ライクなサービスを作成する。

# Public API Route
 - `/api/fetch`
   - Param: offset=`<num>`
 - `/register`
 - `/login`

# JWT Authenticated API
 - `/api/v1` (healthcheck)
 - `/api/v1/upload`
 
 # TODO
 - [ ] テストの作成
 - [ ] フィルターの実装
 - [ ] Streaming API の実装（リアルタイムでアップロードされたファイルを見る）

# How to build

## Frontend

    npm run dev

## Backend

- development: `make build-dev`
- staging: `make build-staging`

## DB migration

    make migrate
    
## Run

    make run
    make stop
    make restart
  
