# `HTTPS_PROXY`の検証
テストで接続先をテストダブル（フェイク/モック/スタブサーバー）に差し替えるのに使えるかどうかを確認。

## `proxy/main.go`
プロキシサーバーを実装。`google.com:443`へのアクセスのみをフェイク(?)サーバーに繋ぐようになっている。

## `main.go`
まず最初にプロキシサーバーを起動。
`https://example.com`と`https://google.com`のそれぞれにGETリクエストを行い、レスポンスボディを出力する。

## 確認できたこと
- `(tls.Config).InsecureSkipVerify = false`を設定する必要がある点を除けば、`HTTPS_PROXY`の環境変数をセットするだけで、プロキシを経由して接続をテストダブルに向けることができる
    - Go 1.20からサポートされた、カバレッジ出力可能なテスト用バイナリを使いたいときに嬉しい（テスト用の細工を最低限にすることで、本番用とテスト用で同じコードを使ったビルドがやりやすい）
    - https://future-architect.github.io/articles/20230203a/
- `DoFunc`を複数仕掛ける場合、それぞれのハンドラーで`*http.Response`を返さないと意図通りに動かない
- S3で`BUCKET_NAME.s3.REGION.amazonaws.com`というスタイルのリクエストが来た場合、`s3.REGION.amazonaws.com/BUCKET_NAME`のスタイルに変換してあげると`github.com/johannesboyne/gofakes3`がうまく動く

## 実行結果

```shell
$ make run
go build -o _proxy ./proxy
HTTPS_PROXY=http://localhost:9999 go run .
2024/10/07 22:58:44 GET https://google.com/foo?search=bar
2024/10/07 22:58:44 [001] INFO: Running 2 CONNECT handlers
2024/10/07 22:58:44 [001] INFO: on 0th handler: &{2 <nil> 0x12b46e0} google.com:443
2024/10/07 22:58:44 [001] INFO: Assuming CONNECT is TLS, mitm proxying it
2024/10/07 22:58:44 [001] INFO: signing for google.com
2024/10/07 22:58:45 [002] INFO: req google.com:443
2024/10/07 22:58:45 [002] INFO: original: https://google.com:443/foo?search=bar
2024/10/07 22:58:45 [002] INFO: proxied: http://localhost:8888/foo?search=bar
2024/10/07 22:58:45 fake: received request: GET /foo?search=bar HTTP/1.1
Host: google.com
Accept-Encoding: gzip
User-Agent: Go-http-client/1.1


2024/10/07 22:58:45 Hello from fake server
2024/10/07 22:58:45 [001] INFO: Exiting on EOF
2024/10/07 22:58:45 [003] INFO: Running 2 CONNECT handlers
2024/10/07 22:58:45 [003] INFO: on 1th handler: &{2 <nil> 0x12b46e0} hey-0.s3.ap-northeast-1.amazonaws.com:443
2024/10/07 22:58:45 [003] INFO: Assuming CONNECT is TLS, mitm proxying it
2024/10/07 22:58:45 [003] INFO: signing for hey-0.s3.ap-northeast-1.amazonaws.com
2024/10/07 22:58:45 [004] INFO: req hey-0.s3.ap-northeast-1.amazonaws.com:443
2024/10/07 22:58:45 [004] INFO: original: https://hey-0.s3.ap-northeast-1.amazonaws.com:443/
2024/10/07 22:58:45 [004] INFO: proxied: http://localhost:7777/hey-0?
2024/10/07 22:58:45 [003] INFO: Exiting on EOF
2024/10/07 22:58:45 [005] INFO: Running 2 CONNECT handlers
2024/10/07 22:58:45 [005] INFO: on 1th handler: &{2 <nil> 0x12b46e0} s3.ap-northeast-1.amazonaws.com:443
2024/10/07 22:58:45 [005] INFO: Assuming CONNECT is TLS, mitm proxying it
2024/10/07 22:58:45 [005] INFO: signing for s3.ap-northeast-1.amazonaws.com
2024/10/07 22:58:45 [006] INFO: req s3.ap-northeast-1.amazonaws.com:443
2024/10/07 22:58:45 [006] INFO: original: https://s3.ap-northeast-1.amazonaws.com:443/?x-id=ListBuckets
2024/10/07 22:58:45 [006] INFO: proxied: http://localhost:7777/?x-id=ListBuckets
2024/10/07 22:58:45 [005] INFO: Exiting on EOF
2024/10/07 22:58:45 [007] INFO: Running 2 CONNECT handlers
2024/10/07 22:58:45 [007] INFO: on 1th handler: &{2 <nil> 0x12b46e0} hey-0.s3.ap-northeast-1.amazonaws.com:443
2024/10/07 22:58:45 [007] INFO: Assuming CONNECT is TLS, mitm proxying it
2024/10/07 22:58:45 [007] INFO: signing for hey-0.s3.ap-northeast-1.amazonaws.com
2024/10/07 22:58:45 [008] INFO: req hey-0.s3.ap-northeast-1.amazonaws.com:443
2024/10/07 22:58:45 [008] INFO: original: https://hey-0.s3.ap-northeast-1.amazonaws.com:443/object1?x-id=PutObject
2024/10/07 22:58:45 [008] INFO: proxied: http://localhost:7777/hey-0/object1?x-id=PutObject
2024/10/07 22:58:45 [007] INFO: Exiting on EOF
2024/10/07 22:58:45 [009] INFO: Running 2 CONNECT handlers
2024/10/07 22:58:45 [009] INFO: on 1th handler: &{2 <nil> 0x12b46e0} hey-0.s3.ap-northeast-1.amazonaws.com:443
2024/10/07 22:58:45 [009] INFO: Assuming CONNECT is TLS, mitm proxying it
2024/10/07 22:58:45 [009] INFO: signing for hey-0.s3.ap-northeast-1.amazonaws.com
2024/10/07 22:58:46 [010] INFO: req hey-0.s3.ap-northeast-1.amazonaws.com:443
2024/10/07 22:58:46 [010] INFO: original: https://hey-0.s3.ap-northeast-1.amazonaws.com:443/object1?x-id=GetObject
2024/10/07 22:58:46 [010] INFO: proxied: http://localhost:7777/hey-0/object1?x-id=GetObject
2024/10/07 22:58:46 value1
2024/10/07 22:58:46 kill proxy: <nil>
```
