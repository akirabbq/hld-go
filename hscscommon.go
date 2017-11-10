package hld

import (
	"bytes"
	"fmt"
	"os"
	"unicode/utf16"
	"unicode/utf8"
)

//HSTestString unicode test string including surrogated pairs
const HSTestString = "English للغة العربية ภาษาไทย 中文𪛖𨳒中文简体 ウェブ全体から検索Русский язык"

const BOM_UTF8 = 1
const BOM_UTF16LE = 2
const BOM_UTF16BE = 3

//GetUTFBomType check if bytes have the UTF8 BOM mark
func GetUTFBomType(bytes []byte) int {
	// if len(*bytes) < 3 {
	// 	return false
	// }
	if (len(bytes) >= 3) && (bytes[0] == 239) && (bytes[1] == 187) && (bytes[2] == 191) {
		return BOM_UTF8
	} else if (len(bytes) >= 2) && (bytes[0] == 255) && (bytes[1] == 254) {
		return BOM_UTF16LE
	} else if (len(bytes) >= 2) && (bytes[0] == 254) && (bytes[1] == 255) {
		return BOM_UTF16BE
	}
	return 0

}

//FileExists _
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

//IsDir _
func IsDir(name string) bool {
	f, err := os.Stat(name)
	return ((err == nil) && f.IsDir())
}

//SameFile _
func SameFile(filename1, filename2 string) bool {
	f1, _ := os.Stat(filename1)
	f2, _ := os.Stat(filename2)

	return os.SameFile(f1, f2)
}

//DecodeUTF16 decode UTF16 string UTF8
func DecodeUTF16(b []byte) (string, error) {

	if len(b)%2 != 0 {
		return "", fmt.Errorf("Must have even length byte slice")
	}

	u16s := make([]uint16, 2)

	ret := &bytes.Buffer{}

	b8buf := make([]byte, 4)

	lb := len(b)
	for i := 0; i < lb; i += 2 {
		u16s[0] = uint16(b[i]) + (uint16(b[i+1]) << 8)

		if u16s[0] >= 0xD800 && u16s[0] <= 0xE000 {
			i += 2
			u16s[1] = uint16(b[i]) + (uint16(b[i+1]) << 8)
		}

		r := utf16.Decode(u16s)

		n := utf8.EncodeRune(b8buf, r[0])
		ret.Write(b8buf[:n])
	}

	return ret.String(), nil
}
