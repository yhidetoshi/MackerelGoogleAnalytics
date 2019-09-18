# MackerelGoogleAnalytics

■ デプロイ
```
export MKRKEY=XXX
export GOOGLE_APPLICATION_CREDENTIALS_JSON="{...}"

curl -X POST https://api.mackerelio.com/api/v0/services \
    -H "X-Api-Key: ${MKRKEY}" \
    -H "Content-Type: application/json" \
    -d '{"name": "GoogleAnalytics", "memo": "google analytics"}'

make build
sls deploy --aws-profile <PROFILE> --viewid <VIEW-ID> --mkrkey ${MKRKEY} --google-apikey ${GOOGLE_APPLICATION_CREDENTIALS_JSON}
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
