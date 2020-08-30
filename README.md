# samを使ったLambda(go言語)+APIGateway+dynamoDB

## sam インストール

```
brew tap aws/tap
brew install aws-sam-cli
sam --version
#SAM CLI, version 1.0.0
```

## ファイル構成

- lambda(lambdaファイル)
    - build.sh ビルド
    - go.mod パッケージ管理
    - go.sum パッケージ記録
    - main.go エントリーポイント
    - main_test.go テストファイル
- event_api.json apigatewayのリクエストのサンプルファイル
- Makefile ビルドやデバッグなどのファイル
- template.yaml samのサンプルファイル

## sam コマンド

プロジェクトフォルダ&&テンプレート作成
```
sam init --runtime go1.x(言語をいろいろ選べる) --name プロジェクト名
```


apigatewayのダミーリクエストファイル作成
```
sam local generate-event apigateway aws-proxy > event_api.json
```

lambdaローカル実行
```
sam local invoke LambdaSampleFunction --event event_api.json

#実行結果 dockerのイメージをダウンロードして実行しているようで4〜5秒かかる。初回はもっと。
updating: main (deflated 49%)
Invoking main (go1.x)
・・・・・・・・・・・・・・
START RequestId: Version: $LATEST
2020/07/23 04:23:17 <nil>
morigami
Function 'LambdaSampleFunction' timed out after 5 seconds

```


## ローカル開発


### lambda
*ソースかえても`make build`しないと結果が変わらない
```
make debug
```

- 自動でビルドファイル構築& sam local invoke 

### apiGetway
ビルド&ローカルAPI
```
make build
sam local start-api

Mounting LambdaSampleFunction at http://127.0.0.1:3000/sample [GET]
You can now browse to the above endpoints to invoke your functions. You do not need to restart/reload SAM CLI while working on your functions, changes will be reflected instantly/automatically. You only need to restart SAM CLI if you update your AWS SAM template
2020-07-23 16:56:18  * Running on http://127.0.0.1:3000/ (Press CTRL+C to quit)

#別ウィンドウで下記コマンドを叩く
curl http://127.0.0.1:3000/sample
#その後、api開いた画面でログが出ていればOK
```

## ビルド　&& デプロイ
```
make build 

sam deploy --guided
#全てyesでOK(ここでおそらくdeploy前準備のようなもの)
#その後、再度
sam deploy 

#Deply this changeset? yで実際にデプロイ
```

## dynamoDB

dockerで仮装環境を構築

GUI
http://dev-host.local:8001/


コマンド実行できる(あまり使わないかも・・・)
http://dev-host.local:8000/shell/

テーブル確認
```
aws dynamodb \
list-tables \
--endpoint-url http://localhost:8000 
```

テーブル作成

```
aws dynamodb \
create-table \
--endpoint-url http://localhost:8000 \
--cli-input-json file://dynamoDB/script/area.json

#--cli-input-json ファイルパスを表示　dynamoDB/script/area.json
```

#全件取得
```
aws dynamodb scan \
--endpoint-url http://localhost:8000 \
--table-name テーブル名
```

#特定キーで取得
```
aws dynamodb get-item \
--endpoint-url http://localhost:8000 \
--table-name テーブル名 \
--key '{"ID":{"N":"1"}}'
```

#保存
```
aws dynamodb put-item \
    --endpoint-url http://localhost:8000 \
    --table-name テーブル名 \
    --item '{
        "product_id": {"N": "3"},
        "product_name": {"S": "本"} ,
        "price_min": {"N": "111"},
        "price_max": {"N": "222"}
      }'
```

boto3

https://boto3.amazonaws.com/v1/documentation/api/latest/index.html

## postgres 
CSV→DBへのインサート
csv \copy area (zip,pref_kana,city_kana,town_kana,pref,city,town) from '/docker-entrypoint-initdb.d/KEN_ALL.CSV' with csv header