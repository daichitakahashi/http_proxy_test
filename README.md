# `HTTPS_PROXY`の検証
テストで接続先をテストダブル（フェイク/モック/スタブサーバー）に差し替えるのに使えるかどうかを確認。

## `proxy/main.go`
プロキシサーバーを実装。"google.com:443"へのアクセスのみをフェイク(?)サーバーに繋ぐようになっている。

## `main.go`
まず最初にプロキシサーバーを起動。
"https://example.com"と"https://google.com"のそれぞれにGETリクエストを行い、レスポンスボディを出力する。

## 確認できたこと
- `(tls.Config).InsecureSkipVerify = false`を設定する必要がある点を除けば、`HTTPS_PROXY`の環境変数をセットするだけで、プロキシを経由して接続をテストダブルに向けることができる

## 実行結果

```shell
$ make run
go build -o _proxy ./proxy
HTTPS_PROXY=http://localhost:9999 go run .
2024/10/03 13:45:12 GET https://example.com
2024/10/03 13:45:12 [001] INFO: Running 1 CONNECT handlers
2024/10/03 13:45:12 [001] INFO: Accepting CONNECT to example.com:443
2024/10/03 13:45:12 <!doctype html>
<html>
<head>
    <title>Example Domain</title>

    <meta charset="utf-8" />
    <meta http-equiv="Content-type" content="text/html; charset=utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <style type="text/css">
    body {
        background-color: #f0f0f2;
        margin: 0;
        padding: 0;
        font-family: -apple-system, system-ui, BlinkMacSystemFont, "Segoe UI", "Open Sans", "Helvetica Neue", Helvetica, Arial, sans-serif;
        
    }
    div {
        width: 600px;
        margin: 5em auto;
        padding: 2em;
        background-color: #fdfdff;
        border-radius: 0.5em;
        box-shadow: 2px 3px 7px 2px rgba(0,0,0,0.02);
    }
    a:link, a:visited {
        color: #38488f;
        text-decoration: none;
    }
    @media (max-width: 700px) {
        div {
            margin: 0 auto;
            width: auto;
        }
    }
    </style>    
</head>

<body>
<div>
    <h1>Example Domain</h1>
    <p>This domain is for use in illustrative examples in documents. You may use this
    domain in literature without prior coordination or asking for permission.</p>
    <p><a href="https://www.iana.org/domains/example">More information...</a></p>
</div>
</body>
</html>


2024/10/03 13:45:12 GET https://google.com/foo?search=bar
2024/10/03 13:45:12 [002] INFO: Running 1 CONNECT handlers
2024/10/03 13:45:12 proxy: HandleConnect google.com:443
2024/10/03 13:45:12 [002] INFO: on 0th handler: &{2 <nil> 0x129bbe0} google.com:443
2024/10/03 13:45:12 [002] INFO: Assuming CONNECT is TLS, mitm proxying it
2024/10/03 13:45:12 [002] INFO: signing for google.com
2024/10/03 13:45:13 [003] INFO: req google.com:443
2024/10/03 13:45:13 [003] INFO: Sending request GET http://localhost:8888?search=bar
2024/10/03 13:45:13 fake: received request: GET /?search=bar HTTP/1.1
Host: google.com
Accept-Encoding: gzip
User-Agent: Go-http-client/1.1


2024/10/03 13:45:13 [003] INFO: resp 200 OK
2024/10/03 13:45:13 Hello from fake server
```
