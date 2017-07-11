# Qiita::Team Images Downloader

## これは何？
Qiita::Teamの画像をローカルにダウンロードします。
Qiita::Teamから移行する際に役立ちます。

## Getting Start
```bash
$ go get github.com/sawadashota/qiita-images-donwloader
```

## Usage
```bash
$ qiita-images-donwloader -q [Qiita::Team name] -qToken [qiita token] [-dir [download directory] -restart-from [Process ID]]
```

|Flag|説明|備考|
|:--|:--|:--|
|-q|Qiita::Teamのチーム名|必須|
|-qToken|Qiitaのアクセストークン|必須|
|-eToken|esaのアクセストークン|必須|
|-dir|画像を保存するディレクトリ|任意。デフォルトは`~/Downloads/qiita-team-images/`|
|-restart-from|何個目の記事から処理を再開させるか|任意。数値のみ|