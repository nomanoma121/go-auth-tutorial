# Go React Auth Template

GoとReactで、認証つきWebアプリを作るためのテンプレートです。

JWTを使った認証を実装しています。

## 初期設定

JWTの署名のため、SECRETという環境変数を設定する必要があります。

```sh
echo "SECRET=$(openssl rand -base64 32)" > ./server/.env
```

これでランダムな文字列がSECRETという環境変数に設定されます。

手動で設定する場合は、`./server/.env`に以下のように書いてください。

```sh
SECRET=ひみつの文字列
```

## フォルダ構成

```txt
.
├── README.md
├── scripts
│   ├── prod.sh
│   ├── dev.sh
│   └── reset.sh
├── client
│   └── ...
└── server
    ├── ...
    └── db.sqlite3
```

### scripts

サーバーを起動や停止、データベースのリセットを行うためのシェルスクリプトが入っています。
中身はただのコマンド群なので、いつも通りcdして手打ちでコマンドを叩いても構いません。

#### 開発環境での起動

```sh
./scripts/dev.sh
```

#### 本番（を再現した）環境での起動

```sh
./scripts/prod.sh
```

#### データベースのリセット

```sh
./scripts/reset.sh
```

### client

フロントエンドのソースコードが入っています。
詳しくは[client/README.md](./client/README.md)を参照してください。

### server

バックエンドのソースコードが入っています。
詳しくは[server/README.md](./server/README.md)を参照してください。

## デプロイについて

デプロイとはアプリをサーバーに載せて全世界に公開することです。
なるべく無料の範囲でデプロイできた方がいいので、以下のサービスを使うことをおすすめします。

### バックエンド

- [Fly.io](https://fly.io/)

### フロントエンド

- [Vercel](https://vercel.com/)
- [Cloudflare Pages](https://pages.cloudflare.com/)

ただ、ここは調べても難しかったりするので、もしわからなければ @sor4chi に聞いてください。

## ライセンス

MIT License

## 作者

[sor4chi](https://github.com/sor4chi)
