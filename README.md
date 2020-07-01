# zengin

Goで全銀総合振込フォーマットのデータを作成するためのpackageです。

もう令和だというのに...

# Usage

```go

// 振込元口座情報を設定する
sender := NewSender(
	"1234567890",
	"ジッケン（カ",
	"0001",
	"ミズホ",
	"211",
	"アオヤマ",
	"1",
	"0000001",
)

// 振込内容
data := []*Transfer{
	NewTransfer("0001", "ミズホ", "211", "アオヤマ", "1", "0000001", "テストコウザ", 100),
	NewTransfer("0005", "ﾐﾂﾋﾞｼﾕｰｴﾌｼﾞｪｲ", "191", "ｱｲﾁｹﾝﾁｮｳ", "2", "0000001", "ﾊﾝｶｸｺﾓｼﾞﾃｽﾄ", 100),
}

// 振込日と振込内容を元に全銀フォーマットデータを生成する
res, err := sender.BuildZenginData("0701", data)
if err != nil {
	t.Error(err)
}

fmt.Println(string(res))
```

# Author

yami20

# License

MIT
