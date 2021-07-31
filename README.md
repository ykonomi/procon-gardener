# procon-gardener

### 概要
提出したACしたコードを自動的に取得してローカルのgitレポジトリに保存してコミットするコマンドラインツールです。

### togatoga/procon-gardener との相違点

- 保存されるディレクトリパスの変更しました
  - 元: {repository_path}/atcoder.jp/abc133/abc133_d/Main.rs
  - これ: {repository_path}/{directory_path}/abc133/d.rs

- 典型90の問題に対応しました
  - 典型90の問題は元のまま実行すると typical90_cl と後半が26進数であるため, 10進数に直しています

- submission.json を一つのファイルにまとめました
  - submission.json が現状重複チェックにしか使われておらず, ファイルがかさばるためです
  - ファイルの中身は保存済みのプログラムIDのリストです
  - (未実装) 将来的には, なくす予定です

- 設定ファイルの項目の追加
  - directory_path: 保存先のディレクトリを指定できるようにしました
  - key_file_name: 重複チェック用のファイルのファイル名です. directory_path直下に保存されます
  - (未実装) github_user_id: 現状, atcoderのuserID でコミットされるため, 指定できるようにしたい

```
{
	"atcoder": {
		"repository_path": "", 
		"directory_path": "atcoder_contest",
		"key_file_name": "key_file.txt",
		"user_id": "",
		"user_email": ""
	}
}
```

- 試運転(Dry Run) の追加
  - 保存される問題のリストが表示されます
  - ABC, ARCが混在されて時代に対応して, 注意書きを表示するようにしています
    - 提出がABCであっても, プログラムIDがARCのものなら, ARCのディレクトリを作って保存されます


### インストール方法
インストールするには[Go](https://golang.org/doc/install)が必要です。

```
go get github.com/ykonomi/procon-gardener
```

### サポート環境
- Linux
- macOs

Windowsは動作確認してません。

### 使い方

1. 設定ファイルの初期化

必要な設定ファイルの作成を行います。`procon-gardener init`を実行してください。
```
% procon-gardener init  
2020/03/21 17:18:36 Initialize your config...
2020/03/21 17:18:36 Initialized your config at  /home/togatoga/.procon-gardener/config.json
```

2. 設定ファイルの編集

初期化した設定ファイルは以下のとおりです。設定ファイルを直接編集もしくは`procon-gardener edit`で編集することができます。  
`EDITOR`の環境変数が設定されていれば、`EDITOR`に設定されているエディタで開きます。そうでなければOS依存`open`のコマンドで開きます。
```
{
	"atcoder": {
		"repository_path": "",
		"user_id": "",
		"user_email": ""
	}
}
```

- `repository_path` アーカイブ先のディレクトリを指定してください
- `user_id` 保存したいユーザーIDを入力してください
- `user_email` `repository_path`が`Git`リポジトリの場合、`git commit`時のメールアドレスに指定されます

`user_email`をGitHubの登録メールアドレスに設定しないとGitHubのアクティビティには反映されません。  
今回は以下のように設定ファイルを編集しました。
```
{
	"atcoder": {
		"repository_path":"/home/togatoga/src/github.com/togatoga/procon-archive",
		"user_id": "togatoga",
		"user_email": "togasakitogatoga+github@gmail.com"
	}
}
```

3. ソースコードのアーカイブ

`procon-gardener archive`を実行すれば自動的にファイルがアーカイブされます。
AtCoderへの負荷対策のため1提出につき1.5秒sleepを行っています、AC数が多い人はしばらくお待ちぐださい。  

```
% procon-gardener archive                                 
2020/03/21 21:19:37 Archiving 1186 code...
2020/03/21 21:19:38 archived the code at  /home/togatoga/src/github.com/togatoga/procon-archive/atcoder.jp/abc133/abc133_d/Main.rs
Main.rs
2020/03/21 21:19:39 archived the code at  /home/togatoga/src/github.com/togatoga/procon-archive/atcoder.jp/abc148/abc148_e/Main.rs
Main.rs
2020/03/21 21:19:40 archived the code at  /home/togatoga/src/github.com/togatoga/procon-archive/atcoder.jp/abc134/abc134_d/Main.rs
Main.rs
2020/03/21 21:19:41 archived the code at  /home/togatoga/src/github.com/togatoga/procon-archive/atcoder.jp/abc115/abc115_d/Main.rs
Main.rs
2020/03/21 21:19:42 archived the code at  /home/togatoga/src/github.com/togatoga/procon-archive/atcoder.jp/agc033/agc033_a/Main.rs
Main.rs
2020/03/21 21:19:43 archived the code at  /home/togatoga/src/github.com/togatoga/procon-archive/atcoder.jp/abc141/abc141_d/Main.rs
Main.rs
2020/03/21 21:19:44 archived the code at  /home/togatoga/src/github.com/togatoga/procon-archive/atcoder.jp/ddcc2020-qual/ddcc2020_qual_d/Main.rs
Main.rs

```


```
$ cd /home/togatoga/src/github.com/togatoga/procon-archive/
$ git log
commit 412134182e09ab0e165e3499020bcebd80ecfe6d (HEAD -> master)
Author: togatoga <togasakitogatoga+github@gmail.com>
Date:   Sun Mar 15 15:08:28 2020 +0900

    [AC] abc141 abc141_d

commit d8d36f6cc5ca35ab433b5e6fbabe7ca4e4f7f8bd
Author: togatoga <togasakitogatoga+github@gmail.com>
Date:   Sun Mar 15 16:54:37 2020 +0900

    [AC] agc033 agc033_a

commit abf4779970804c3fd6fe8bf2d7b2ac02a15e3d34
Author: togatoga <togasakitogatoga+github@gmail.com>
Date:   Sun Mar 15 18:29:50 2020 +0900

    [AC] abc115 abc115_d

commit 2615058a482a7f7589d900fd5c84ff8a5ebfc871
Author: togatoga <togasakitogatoga+github@gmail.com>
Date:   Mon Mar 16 09:42:47 2020 +0900

    [AC] abc134 abc134_d

commit b84a716762fd4df6df19121b5599b526f2fdba89
Author: togatoga <togasakitogatoga+github@gmail.com>
Date:   Wed Mar 18 22:12:23 2020 +0900

    [AC] abc148 abc148_e

commit 7f905746a102190f054430e696da8ab742cffb5c
Author: togatoga <togasakitogatoga+github@gmail.com>
Date:   Fri Mar 20 06:30:19 2020 +0900

    [AC] abc133 abc133_d

```

### 不具合があったら

GitHubのissue へ

