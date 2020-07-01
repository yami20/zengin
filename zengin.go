package zengin

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"golang.org/x/text/width"
)

// accountType 口座種別
type accountType string

const (
	// AccountTypeOrdinary 普通預金
	AccountTypeOrdinary accountType = "1"
	// AccountTypeCurrent 当座預金
	AccountTypeCurrent accountType = "2"
)

// Sender 送金者
type Sender struct {
	clientCode    string
	clientName    string
	bankCode      string
	bankName      string
	branchCode    string
	branchName    string
	accountType   accountType
	accountNumber string
}

// NewSender 送金者情報を生成する
func NewSender(
	clientCode string,
	clientName string,
	bankCode string,
	bankName string,
	branchCode string,
	branchName string,
	accountType accountType,
	accountNumber string,
) *Sender {
	return &Sender{
		clientCode:    clientCode,
		clientName:    clientName,
		bankCode:      bankCode,
		bankName:      bankName,
		branchCode:    branchCode,
		branchName:    branchName,
		accountType:   accountType,
		accountNumber: accountNumber,
	}
}

// Transfer 送金内容
type Transfer struct {
	bankCode         string
	bankName         string
	branchCode       string
	branchName       string
	accountType      accountType
	accountNumber    string
	accountOwnerName string
	amount           int
}

// NewTransfer 振込１件ずつのデータを作成する
func NewTransfer(
	bankCode string,
	bankName string,
	branchCode string,
	branchName string,
	accountType accountType,
	accountNumber string,
	accountOwnerName string,
	amount int,
) *Transfer {
	return &Transfer{
		bankCode:         bankCode,
		bankName:         bankName,
		branchCode:       branchCode,
		branchName:       branchName,
		accountType:      accountType,
		accountNumber:    accountNumber,
		accountOwnerName: accountOwnerName,
		amount:           amount,
	}
}

// BuildZenginData 全銀データ生成
func (s *Sender) BuildZenginData(date string, transfers []*Transfer) ([]byte, error) {
	contents, err := s.buildUTFZenginData(date, transfers)
	if err != nil {
		return nil, err
	}

	joined := strings.Join(contents, "\r\n")
	sjisRes, _, err := transform.String(japanese.ShiftJIS.NewEncoder(), joined)
	if err != nil {
		return nil, err
	}
	return []byte(sjisRes), nil

}

func (s *Sender) buildUTFZenginData(
	date string,
	transfers []*Transfer,
) ([]string, error) {
	contents := []string{}
	header, err := s.header(date)
	if err != nil {
		return nil, err
	}
	contents = append(contents, header)

	count := 0
	amount := 0
	for _, transfer := range transfers {
		line, err := transfer.toLine()
		if err != nil {
			return nil, err
		}
		contents = append(contents, line)

		count++
		amount += transfer.amount
	}

	contents = append(contents, buildTrailer(count, amount))
	contents = append(contents, ender)
	return contents, nil
}

func (s *Sender) header(date string) (string, error) {
	res := ""
	headerFmt := "1210%010s%40s%04s%04s%15s%03s%15s%01s%07s%17s"

	normalizedClientCode, err := normalizeNumber(s.clientCode)
	if err != nil {
		return res, err
	}
	if len([]rune(normalizedClientCode)) > 10 {
		return res, fmt.Errorf("clientCode[%s] is too long. max length is 10", normalizedClientCode)
	}

	normalizedClientName, err := normalizeKana(s.clientName)
	if err != nil {
		return res, err
	}
	if len([]rune(normalizedClientName)) > 40 {
		return res, fmt.Errorf("clientName[%s] is too long. max length is 40", normalizedClientName)
	}

	normalizedDate, err := normalizeNumber(date)
	if err != nil {
		return res, err
	}
	if len([]rune(normalizedDate)) != 4 {
		return res, fmt.Errorf("date[%s] length must be 4", normalizedDate)
	}

	normalizedBankCode, err := normalizeNumber(s.bankCode)
	if err != nil {
		return res, err
	}
	if len([]rune(normalizedBankCode)) > 4 {
		return res, fmt.Errorf("bankCode[%s] is too long. max length is 4", normalizedBankCode)
	}

	normalizedBankName, err := normalizeKana(s.bankName)
	if err != nil {
		return res, err
	}
	if len([]rune(normalizedBankName)) > 15 {
		return res, fmt.Errorf("bankName[%s] is too long. max length is 15", normalizedBankName)
	}

	normalizedBranchCode, err := normalizeNumber(s.branchCode)
	if err != nil {
		return res, err
	}
	if len([]rune(normalizedBranchCode)) > 3 {
		return res, fmt.Errorf("branchCode[%s] is too long. max length is 3", normalizedBranchCode)
	}

	normalizedBranchName, err := normalizeKana(s.branchName)
	if err != nil {
		return res, err
	}
	if len([]rune(normalizedBranchName)) > 15 {
		return res, fmt.Errorf("branchName[%s] is too long. max length is 15", normalizedBranchName)
	}

	normalizedAccountType, err := normalizeNumber(string(s.accountType))
	if err != nil {
		return res, err
	}
	if len([]rune(normalizedAccountType)) > 1 {
		return res, fmt.Errorf("accountType[%s] is too long. max length is 1", normalizedAccountType)
	}

	normalizedAccountNumber, err := normalizeNumber(s.accountNumber)
	if err != nil {
		return res, err
	}
	if len([]rune(normalizedAccountNumber)) > 7 {
		return res, fmt.Errorf("accountNumber[%s] is too long. max length is 7", normalizedAccountNumber)
	}

	res = fmt.Sprintf(
		headerFmt,
		normalizedClientCode,
		normalizedClientName,
		normalizedDate,
		normalizedBankCode,
		normalizedBankName,
		normalizedBranchCode,
		normalizedBranchName,
		normalizedAccountType,
		normalizedAccountNumber,
		"",
	)
	return res, nil

}

func (t *Transfer) toLine() (string, error) {
	res := ""
	lineFmt := "2%04s%15s%03s%15s    %1s%07s%30s%010s1%29s"

	normalizedBankCode, err := normalizeNumber(t.bankCode)
	if err != nil {
		return res, err
	}
	if len([]rune(normalizedBankCode)) > 4 {
		return res, fmt.Errorf("bankCode[%s] is too long. max length is 4", normalizedBankCode)
	}

	normalizedBankName, err := normalizeKana(t.bankName)
	if err != nil {
		return res, err
	}
	if len([]rune(normalizedBankName)) > 15 {
		return res, fmt.Errorf("bankName[%s] is too long. max length is 15", normalizedBankName)
	}

	normalizedBranchCode, err := normalizeNumber(t.branchCode)
	if err != nil {
		return res, err
	}
	if len([]rune(normalizedBranchCode)) > 3 {
		return res, fmt.Errorf("branchCode[%s] is too long. max length is 3", normalizedBranchCode)
	}

	normalizedBranchName, err := normalizeKana(t.branchName)
	if err != nil {
		return res, err
	}
	if len([]rune(normalizedBranchName)) > 15 {
		return res, fmt.Errorf("branchName[%s] is too long. max length is 15", normalizedBranchName)
	}

	normalizedAccountType, err := normalizeNumber(string(t.accountType))
	if err != nil {
		return res, err
	}
	if len([]rune(normalizedAccountType)) > 1 {
		return res, fmt.Errorf("accountType[%s] is too long. max length is 1", normalizedAccountType)
	}

	normalizedAccountNumber, err := normalizeNumber(t.accountNumber)
	if err != nil {
		return res, err
	}
	if len([]rune(normalizedAccountNumber)) > 7 {
		return res, fmt.Errorf("accountNumber[%s] is too long. max length is 7", normalizedAccountNumber)
	}

	normalizedAccountOwnerName, err := normalizeKana(t.accountOwnerName)
	if err != nil {
		return res, err
	}
	if len([]rune(normalizedAccountOwnerName)) > 30 {
		return res, fmt.Errorf("accountOwnerName[%s] is too long. max length is 30", normalizedAccountOwnerName)
	}

	normalizedAmount, err := normalizeNumber(fmt.Sprintf("%d", t.amount))
	if err != nil {
		return res, err
	}
	if len([]rune(normalizedAmount)) > 10 {
		return res, fmt.Errorf("amount[%s] is too long. max length is 10", normalizedAmount)
	}

	res = fmt.Sprintf(
		lineFmt,
		normalizedBankCode,
		normalizedBankName,
		normalizedBranchCode,
		normalizedBranchName,
		normalizedAccountType,
		normalizedAccountNumber,
		normalizedAccountOwnerName,
		normalizedAmount,
		"",
	)
	return res, nil
}

func buildTrailer(
	num int,
	amount int,
) string {
	return fmt.Sprintf("8%06d%012d%101s", num, amount, "")
}

func normalizeKana(s string) (string, error) {
	narrow := strings.ToUpper(width.Narrow.String(s))
	res := []rune{}
	for _, r := range narrow {
		if mapped, ok := kanaMap[r]; ok {
			res = append(res, []rune(mapped)...)
		} else {
			res = append(res, r)
		}
	}

	if !kanaSymbol.MatchString(string(res)) {
		return "", fmt.Errorf("%s includes invalid character. available characters are %s", string(res), kanaSymbol)
	}
	return string(res), nil
}

func normalizeNumber(s string) (string, error) {
	narrow := width.Narrow.String(s)
	if !num.MatchString(narrow) {
		return "", fmt.Errorf("%s includes invalid character. available charaters are %s", narrow, num)
	}
	return narrow, nil
}

var ender = fmt.Sprintf("9%119s", "")
var kanaSymbol = regexp.MustCompile(`^[ｦ-ﾟ()\-\.¥\s\/]*$`)
var num = regexp.MustCompile(`^[0-9]*$`)
var kanaMap = map[rune]string{ // 小文字・濁点・半濁点の半角変換解決するためのmap
	rune('ァ'): "ｱ",
	rune('ィ'): "ｲ",
	rune('ゥ'): "ｳ",
	rune('ェ'): "ｴ",
	rune('ォ'): "ｵ",
	rune('ッ'): "ﾂ",
	rune('ャ'): "ﾔ",
	rune('ュ'): "ﾕ",
	rune('ョ'): "ﾖ",
	rune('ｧ'): "ｱ",
	rune('ｨ'): "ｲ",
	rune('ｩ'): "ｳ",
	rune('ｪ'): "ｴ",
	rune('ｫ'): "ｵ",
	rune('ｯ'): "ﾂ",
	rune('ｬ'): "ﾔ",
	rune('ｭ'): "ﾕ",
	rune('ｮ'): "ﾖ",
	rune('ヴ'): "ｳﾞ",
	rune('ガ'): "ｶﾞ",
	rune('ギ'): "ｷﾞ",
	rune('グ'): "ｸﾞ",
	rune('ゲ'): "ｹﾞ",
	rune('ゴ'): "ｺﾞ",
	rune('ザ'): "ｻﾞ",
	rune('ジ'): "ｼﾞ",
	rune('ズ'): "ｽﾞ",
	rune('ゼ'): "ｾﾞ",
	rune('ゾ'): "ｿﾞ",
	rune('ダ'): "ﾀﾞ",
	rune('ヂ'): "ﾁﾞ",
	rune('ヅ'): "ﾂﾞ",
	rune('デ'): "ﾃﾞ",
	rune('ド'): "ﾄﾞ",
	rune('バ'): "ﾊﾞ",
	rune('ビ'): "ﾋﾞ",
	rune('ブ'): "ﾌﾞ",
	rune('ベ'): "ﾍﾞ",
	rune('ボ'): "ﾎﾞ",
	rune('パ'): "ﾊﾟ",
	rune('ピ'): "ﾋﾟ",
	rune('プ'): "ﾌﾟ",
	rune('ペ'): "ﾍﾟ",
	rune('ポ'): "ﾎﾟ",
}
