# MackerelGoogleAnalytics

■ デプロイ
```
AWS Lambdaに GoogleAnalytics の名前でfunctionを作成 （ランタイムはGo1.x）
環境変数に以下をセットする
- GOOGLE_APPLICATION_CREDENTIALS_JSON: GCPのクレデンシャル(json)
- MKRKEY: mackerelのapikey
- TZ: Asia/Tokyo
- VIEW_ID: GA

export MKRKEY=XXX

curl -X POST https://api.mackerelio.com/api/v0/services \
    -H "X-Api-Key: ${MKRKEY}" \
    -H "Content-Type: application/json" \
    -d '{"name": "GoogleAnalytics", "memo": "google analytics"}'

make build
zip main.zip main

aws lambda --profile <PROFILE> update-function-code --function-name GoogleAnalytics --zip-file fileb://main.zip --region ap-northeast-1
```


`$ make help`
```
build:             Build binaries
build-deps:        Setup build
deps:              Install dependencies
devel-deps:        Setup development
help:              Show help
lint:              Lint
```
