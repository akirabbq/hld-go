package hld

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"unicode/utf16"
	"unicode/utf8"
)

//HSTestString unicode test string including surrogated pairs
const HSTestString = "English للغة العربية ภาษาไทย 中文𪛖𨳒中文简体 ウェブ全体から検索Русский язык"

const BOM_UTF8 = 1
const BOM_UTF16LE = 2
const BOM_UTF16BE = 3

//GetUTFBomType check if bytes have the UTF8 BOM mark, and return BOM len
func GetUTFBomType(bytes []byte) (int, int) {
	// if len(*bytes) < 3 {
	// 	return false
	// }
	if (len(bytes) >= 3) && (bytes[0] == 239) && (bytes[1] == 187) && (bytes[2] == 191) {
		return BOM_UTF8, 3
	} else if (len(bytes) >= 2) && (bytes[0] == 255) && (bytes[1] == 254) {
		return BOM_UTF16LE, 2
	} else if (len(bytes) >= 2) && (bytes[0] == 254) && (bytes[1] == 255) {
		return BOM_UTF16BE, 2
	}
	return 0, 0

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

//DecodeUTF16 decode UTF16 string UTF8, assume little indian and no BOM mark
func DecodeUTF16(b []byte, order binary.ByteOrder) (string, error) {

	if len(b)%2 != 0 {
		return "", fmt.Errorf("Must have even length byte slice")
	}

	u16s := make([]uint16, 2)

	ret := &bytes.Buffer{}

	b8buf := make([]byte, 4)

	lb := len(b)
	for i := 0; i < lb; i += 2 {

		u16s[0] = order.Uint16(b[i:])

		if u16s[0] >= 0xD800 && u16s[0] <= 0xE000 {
			i += 2
			u16s[1] = order.Uint16(b[i:])
		}

		r := utf16.Decode(u16s)
		n := utf8.EncodeRune(b8buf, r[0])
		ret.Write(b8buf[:n])
	}

	return ret.String(), nil
}

//HSFileToString read text file into strings
func HSFileToString(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	bomtype, bomlen := GetUTFBomType(data)
	switch bomtype {
	case BOM_UTF8:
		return string(data[bomlen:]), nil
	case BOM_UTF16LE:
		return DecodeUTF16(data[bomlen:], binary.LittleEndian)
	case BOM_UTF16BE:
		return DecodeUTF16(data[bomlen:], binary.BigEndian)
	default:
		return string(data), nil
		//log.Println("default")
	}
}
