package hld

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode/utf16"
	"unicode/utf8"
)

const (
	//HSTestString unicode test string including surrogated pairs
	HSTestString = "English للغة العربية ภาษาไทย 中文𪛖𨳒中文简体 ウェブ全体から検索Русский язык"

	//BOMUtf8 utf8 bom
	BOMUtf8 = 1
	//BOMUtf16 BOM
	BOMUtf16 = 2
	//BOMUtf16be BOM
	BOMUtf16be = 3
)

//JSTimeToTime JavaScript time to golang time
func JSTimeToTime(t int64) time.Time {
	//1510658333603 - get the most right 3 digits, and convert to nano second
	//javascript time is only down to miliseconds
	return time.Unix(t/1000, (t%1000)*1000*1000)
}

//GetUTFBomType check if bytes have the UTF8 BOM mark, and return BOM len
func GetUTFBomType(bytes []byte) (int, int) {
	// if len(*bytes) < 3 {
	// 	return false
	// }
	if (len(bytes) >= 3) && (bytes[0] == 239) && (bytes[1] == 187) && (bytes[2] == 191) {
		return BOMUtf8, 3
	} else if (len(bytes) >= 2) && (bytes[0] == 255) && (bytes[1] == 254) {
		return BOMUtf16, 2
	} else if (len(bytes) >= 2) && (bytes[0] == 254) && (bytes[1] == 255) {
		return BOMUtf16be, 2
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

//BytesToFile bytes to file
func BytesToFile(bytes []byte, filename string) error {

	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		return err
	}
	_, err = file.Write(bytes)
	return err

}

//FileToBytes read text file byte, will omit the UTF BOM
func FileToBytes(filename string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	_, bomlen := GetUTFBomType(bytes)
	if bomlen > 0 {
		return bytes[bomlen:], nil
	}
	return bytes, nil
}

//HSFileToString read text file into strings
func HSFileToString(filename string) (string, error) {
	if file, err := os.Open(filename); err == nil {
		data, err := ioutil.ReadAll(file)
		if err != nil {
			return "", err
		}

		bomtype, bomlen := GetUTFBomType(data)
		switch bomtype {
		case BOMUtf8:
			return string(data[bomlen:]), nil
		case BOMUtf16:
			return DecodeUTF16(data[bomlen:], binary.LittleEndian)
		case BOMUtf16be:
			return DecodeUTF16(data[bomlen:], binary.BigEndian)
		default:
			return string(data), nil
			//log.Println("default")
		}
	} else {
		return "", err
	}
}

//HSStringList string list
type HSStringList struct {
	text      string
	Lines     []string
	lineBreak string
}

//AssignString assign string and linebreak, if linebreak is empty then linebreak will be "\n"
func (sl *HSStringList) AssignString(text string, lineBreak string) {
	sl.text = text
	sl.lineBreak = lineBreak
	if sl.lineBreak == "" {
		sl.lineBreak = "\n"
	}
	sl.Lines = strings.Split(sl.text, sl.lineBreak)
}

//Count return number of line
func (sl *HSStringList) Count() int {
	return len(sl.Lines)
}

//LoadFromFile load strings from file
func (sl *HSStringList) LoadFromFile(filename string) bool {
	if s, err := HSFileToString(filename); err == nil {
		sl.AssignString(s, "\n")
		return true
	}
	return false
}

//ChangeFileExt change file extension
func ChangeFileExt(filename string, ext string) string {
	oldExt := filepath.Ext(filename)
	if oldExt != "" {
		return filename[0:len(filename)-len(oldExt)] + ext
	}
	return filename + ext
}

//CopyFile copy file
func CopyFile(source string, target string) error {
	file, err := os.Open(source)
	if err != nil {
		return err
	}
	defer file.Close()

	w, err := os.Create(target)
	if err != nil {
		return err
	}
	defer w.Close()

	buffSize := 32 * 1024
	buffer := make([]byte, buffSize)

	for {
		n, err := file.Read(buffer)
		if err != nil {
			return err
		}
		if n > 0 {
			w.Write(buffer[:n])
		}
		if n < buffSize {
			break
		}
	}
	return nil
}

//AppendPathSlash append slash to a directory
func AppendPathSlash(pathname string) string {
	if (len(pathname) > 0) && (pathname[len(pathname)-1] != filepath.Separator) {
		return pathname + "/"
	}
	return pathname

}

//GetDirFiles get all files under a directory
func GetDirFiles(dirname string, r bool) []string {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil
	}
	dirname, _ = filepath.Abs(dirname)
	var results []string
	for _, file := range files {
		if !file.IsDir() {
			results = append(results, filepath.Join(dirname, file.Name()))
		} else if r {
			results = append(results, GetDirFiles(filepath.Join(dirname, file.Name()), r)...)

		}
	}
	return results
}

//FilterFiles filter out file list by extension
func FilterFiles(files []string, ext string) []string {
	var outfiles []string
	re := regexp.MustCompile(fmt.Sprintf("(.)\\.(?i)(%s)$", ext))
	for _, file := range files {
		if re.MatchString(file) {
			outfiles = append(outfiles, file)
		}
	}
	return outfiles
}
