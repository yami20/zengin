package zengin

import (
	"testing"
)

// TestBuildZenginData 全銀データ生成(sjis化)
func TestBuildZenginData(t *testing.T) {
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

	data := []*Transfer{
		NewTransfer("0001", "ミズホ", "211", "アオヤマ", "1", "0000001", "テストコウザ", 100),
		NewTransfer("0005", "ﾐﾂﾋﾞｼﾕｰｴﾌｼﾞｪｲ", "191", "ｱｲﾁｹﾝﾁｮｳ", "2", "0000001", "ﾊﾝｶｸｺﾓｼﾞﾃｽﾄ", 100),
	}

	date := "0701"

	res, err := sender.BuildZenginData(date, data)
	if err != nil {
		t.Error(err)
	}

	if len(res) != 608 { // 120 x 5 + CRLFx4
		t.Errorf("expected 608 bytes, but got %d bytes", len(res))
	}
}

// TestBuildZenginData 全銀データ生成(utf8 slice)
func TestBuildUTFZenginData(t *testing.T) {
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

	data := []*Transfer{
		NewTransfer("0001", "ミズホ", "211", "アオヤマ", "1", "0000001", "テストコウザ", 100),
		NewTransfer("0005", "ﾐﾂﾋﾞｼﾕｰｴﾌｼﾞｪｲ", "191", "ｱｲﾁｹﾝﾁｮｳ", "2", "0000001", "ﾊﾝｶｸｺﾓｼﾞﾃｽﾄ", 100),
	}

	date := "0701"

	res, err := sender.buildUTFZenginData(date, data)
	if err != nil {
		t.Error(err)
	}

	if len(res) != 5 {
		t.Errorf("expected 5 lines, but got %d lines", len(res))
	}

	headerPostfix := "                 "
	expectedL0 := "12101234567890                                 ｼﾞﾂｹﾝ(ｶ07010001           ﾐｽﾞﾎ211           ｱｵﾔﾏ10000001" + headerPostfix
	if res[0] != expectedL0 {
		t.Errorf("expected line0 is [%s], but got [%s]", expectedL0, res[0])
	}

	dataPostfix := "                             "
	expectedL1 := "20001           ﾐｽﾞﾎ211           ｱｵﾔﾏ    10000001                       ﾃｽﾄｺｳｻﾞ00000001001" + dataPostfix
	if res[1] != expectedL1 {
		t.Errorf("expected line1 is [%s], but got [%s]", expectedL1, res[1])
	}

	expectedL2 := "20005  ﾐﾂﾋﾞｼﾕｰｴﾌｼﾞｴｲ191       ｱｲﾁｹﾝﾁﾖｳ    20000001                   ﾊﾝｶｸｺﾓｼﾞﾃｽﾄ00000001001" + dataPostfix
	if res[2] != expectedL2 {
		t.Errorf("expected line2 is [%s], but got [%s]", expectedL2, res[2])
	}

	expectedL3 := "8000002000000000200                                                                                                     "
	if res[3] != expectedL3 {
		t.Errorf("expected line3 is [%s], but got [%s]", expectedL3, res[3])
	}

	expectedL4 := "9                                                                                                                       "
	if res[4] != expectedL4 {
		t.Errorf("expected line4 is [%s], but got [%s]", expectedL4, res[4])
	}

}

// TestBuildHeader ヘッダレコードのチェック
func TestBuildHeader(t *testing.T) {
	type in struct {
		clientCode    string
		clientName    string
		date          string
		bankCode      string
		bankName      string
		branchCode    string
		branchName    string
		accountType   string
		accountNumber string
	}
	type out struct {
		err  error
		line string
	}
	fixedBrank := "                 "

	ins := []in{
		{"1234567890", "ジッケン（カ", "1231", "0001", "ミズホ", "211", "アオヤマ", "1", "0000001"},
		{"1234567890", "ｼﾞﾂｹﾝ(ｶ", "1231", "0001", "ﾐｽﾞﾎ", "211", "ｱオヤマ", "1", "0000001"},
		{"1234567890", "ｼﾞｯｹﾝ（カ", "1231", "000１", "ﾐｽﾞﾎ", "２１１", "ｱオヤマ", "1", "0000001"},
		{"1234567890", "カ)ケンショウ", "0101", "0005", "ミツビシユーエフジエイ", "191", "アイチケンチヨウ", "1", "1000001"},
		{"1234567890", "カ)ケンシヨウ", "0101", "0005", "ミツビシユーエフジェイ", "191", "アイチケンチョウ", "1", "1000001"},
		{"1234567890", "ｶ)ｹﾝｼﾖｳ", "0101", "0005", "ﾐﾂﾋﾞｼﾕｰｴﾌｼﾞｪｲ", "191", "ｱｲﾁｹﾝﾁｮｳ", "1", "1000001"},
		{"1234567890", "ｶ)ｹﾝｼｮｳ", "0101", "0009", "ミツイスミトモ", "969", "アオイ", "2", "1000001"},
	}
	outs := []out{
		{nil, "12101234567890                                 ｼﾞﾂｹﾝ(ｶ12310001           ﾐｽﾞﾎ211           ｱｵﾔﾏ10000001" + fixedBrank},
		{nil, "12101234567890                                 ｼﾞﾂｹﾝ(ｶ12310001           ﾐｽﾞﾎ211           ｱｵﾔﾏ10000001" + fixedBrank},
		{nil, "12101234567890                                 ｼﾞﾂｹﾝ(ｶ12310001           ﾐｽﾞﾎ211           ｱｵﾔﾏ10000001" + fixedBrank},
		{nil, "12101234567890                                 ｶ)ｹﾝｼﾖｳ01010005  ﾐﾂﾋﾞｼﾕｰｴﾌｼﾞｴｲ191       ｱｲﾁｹﾝﾁﾖｳ11000001" + fixedBrank},
		{nil, "12101234567890                                 ｶ)ｹﾝｼﾖｳ01010005  ﾐﾂﾋﾞｼﾕｰｴﾌｼﾞｴｲ191       ｱｲﾁｹﾝﾁﾖｳ11000001" + fixedBrank},
		{nil, "12101234567890                                 ｶ)ｹﾝｼﾖｳ01010005  ﾐﾂﾋﾞｼﾕｰｴﾌｼﾞｴｲ191       ｱｲﾁｹﾝﾁﾖｳ11000001" + fixedBrank},
		{nil, "12101234567890                                 ｶ)ｹﾝｼﾖｳ01010009        ﾐﾂｲｽﾐﾄﾓ969            ｱｵｲ21000001" + fixedBrank},
	}

	for k, in := range ins {
		sender := NewSender(
			in.clientCode,
			in.clientName,
			in.bankCode,
			in.bankName,
			in.branchCode,
			in.branchName,
			accountType(in.accountType),
			in.accountNumber,
		)
		res, err := sender.header(in.date)

		if err != outs[k].err {
			t.Errorf("cases[%d] expected err is [%v], but got [%v]", k, outs[k].err, err)
		}
		if len([]rune(res)) != 120 {
			t.Errorf("cases[%d] zengin record length must be 120, but got %d", k, len([]rune(res)))
		}
		if res != outs[k].line {
			t.Errorf("cases[%d] expected res is [%v], but got [%v]", k, outs[k].line, res)
		}
	}
}

// TestBuildData データレコードのチェック
func TestBuildData(t *testing.T) {
	type in struct {
		bankCode         string
		bankName         string
		branchCode       string
		branchName       string
		accountType      string
		accountNumber    string
		accountOwnerName string
		amount           int
	}
	type out struct {
		err  error
		line string
	}
	fixedBrank := "                             "

	ins := []in{
		{"0001", "ミズホ", "211", "アオヤマ", "1", "0000001", "テストコウザ", 100},
		{"0001", "ﾐｽﾞﾎ", "211", "ｱオヤマ", "1", "0000001", "テストコウザ", 100},
		{"000１", "ﾐｽﾞﾎ", "２１１", "ｱオヤマ", "1", "0000001", "テストコウザ", 100},
		{"0005", "ミツビシユーエフジエイ", "191", "アイチケンチヨウ", "1", "0000001", "ハイフンテスト", 100},
		{"0005", "ミツビシユーエフジェイ", "191", "アイチケンチョウ", "1", "0000001", "コモジテスト", 100},
		{"0005", "ﾐﾂﾋﾞｼﾕｰｴﾌｼﾞｪｲ", "191", "ｱｲﾁｹﾝﾁｮｳ", "1", "0000001", "ﾊﾝｶｸｺﾓｼﾞﾃｽﾄ", 100},
		{"0009", "ミツイスミトモ", "969", "アオイ", "2", "1000001", "ィロンﾅ（(モジ)）ｼｭノパターンｯﾎﾟｲ　ョ- .", 99998},
	}
	outs := []out{
		{nil, "20001           ﾐｽﾞﾎ211           ｱｵﾔﾏ    10000001                       ﾃｽﾄｺｳｻﾞ00000001001" + fixedBrank},
		{nil, "20001           ﾐｽﾞﾎ211           ｱｵﾔﾏ    10000001                       ﾃｽﾄｺｳｻﾞ00000001001" + fixedBrank},
		{nil, "20001           ﾐｽﾞﾎ211           ｱｵﾔﾏ    10000001                       ﾃｽﾄｺｳｻﾞ00000001001" + fixedBrank},
		{nil, "20005  ﾐﾂﾋﾞｼﾕｰｴﾌｼﾞｴｲ191       ｱｲﾁｹﾝﾁﾖｳ    10000001                       ﾊｲﾌﾝﾃｽﾄ00000001001" + fixedBrank},
		{nil, "20005  ﾐﾂﾋﾞｼﾕｰｴﾌｼﾞｴｲ191       ｱｲﾁｹﾝﾁﾖｳ    10000001                       ｺﾓｼﾞﾃｽﾄ00000001001" + fixedBrank},
		{nil, "20005  ﾐﾂﾋﾞｼﾕｰｴﾌｼﾞｴｲ191       ｱｲﾁｹﾝﾁﾖｳ    10000001                   ﾊﾝｶｸｺﾓｼﾞﾃｽﾄ00000001001" + fixedBrank},
		{nil, "20009        ﾐﾂｲｽﾐﾄﾓ969            ｱｵｲ    21000001  ｲﾛﾝﾅ((ﾓｼﾞ))ｼﾕﾉﾊﾟﾀｰﾝﾂﾎﾟｲ ﾖ- .00000999981" + fixedBrank},
	}

	for k, in := range ins {
		transfer := NewTransfer(
			in.bankCode,
			in.bankName,
			in.branchCode,
			in.branchName,
			accountType(in.accountType),
			in.accountNumber,
			in.accountOwnerName,
			in.amount,
		)

		res, err := transfer.toLine()

		if err != outs[k].err {
			t.Errorf("cases[%d] expected err is [%v], but got [%v]", k, outs[k].err, err)
		}
		if len([]rune(res)) != 120 {
			t.Errorf("cases[%d] zengin record length must be 120, but got %d", k, len([]rune(res)))
		}
		if res != outs[k].line {
			t.Errorf("cases[%d] expected res is [%v], but got [%v]", k, outs[k].line, res)
		}
	}
}
