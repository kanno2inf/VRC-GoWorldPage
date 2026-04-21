# VRC-GoWorldPage

VRChatの写真からワールドページを開くプログラム。(Windows用)

## 使い方

[Release](https://github.com/kanno2inf/VRC-GoWorldPage/releases)からzipファイルをダウンロードして展開してください。

**VRC-GoWorldPage.exe** に画像ファイルをドラッグアンドドロップするとワールドページを開きます。

※VRChatカメラの「メタデータの保存(Save Metadata)」設定を有効にした写真のみ対応しています。

## 仕組み
画像ファイルにメタデータとして保存されているWorldID情報を読み取ってワールドページを開きます。

## ビルド方法

`go build`してください。

### Windows

`build.bat`を実行すると`build`ディレクトリ内にビルドされます。
