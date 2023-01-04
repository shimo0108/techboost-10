# ca-backend-boost最終課題

- ca-backend-boost最終課題内容
TOKYO OPEN DATA(https://portal.data.metro.tokyo.lg.jp/)
ここから好きなデータを選んで下記の機能を実装してください。またこれ以外の
機能を実装しても構いません。
AWSにデプロイして誰もが使えるサービスを目指してください。
1) 1日一回 データの更新チェックを行う。更新があった場合CSV等のデータから
PostgreSQLにデータを格納する
2) 検索APIを定義してAPIで検索した結果のデータを返す

# 都内病院検索API
[診療・検査医療機関の一覧リスト
](https://catalog.data.metro.tokyo.lg.jp/dataset/t000010d0000000095/resource/176ddddc-7297-48e3-960e-fc79711c445e)から都内で診療可能な病院のCSVデータ参照し、DBへ保存。cronを用いて毎日10時に再更新できるように実装。
さらに市区町村名で該当する病院を検索できるAPI機能を実装しました。

## APIドキュメント
https://shimo0108.github.io/techboost-10/
