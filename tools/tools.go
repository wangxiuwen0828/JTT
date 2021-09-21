package tools

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/axgle/mahonia"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

//MD5加密，生成32位MD5字符串
func SignMD5(signBefore string) string {
	h := md5.New()
	h.Write([]byte(signBefore))
	cipherStr := h.Sum(nil)
	result := hex.EncodeToString(cipherStr)
	return result
}

func GetRangeNum(min, max int64) int64 {
	num := rand.Int63n(max-min) + min
	if num%2 == 0 {
		return num
	} else {
		return num - 1
	}
}

//字符串拼接
func StringsJoin(str ...string) string {
	var b bytes.Buffer
	strLen := len(str)
	if strLen == 0 {
		return ""
	}
	for i := 0; i < strLen; i++ {
		b.WriteString(str[i])
	}

	return b.String()
}

//生成UUID
func GetUUID() string {
	return strings.ToUpper(uuid.Must(uuid.NewV4()).String())
}

//返回当前工作目录的绝对路径
func GetAbsPath() string {
	path, err := filepath.Abs(filepath.Dir(os.Args[0])) //作为服务时使用
	//path, err := os.Getwd()
	if err != nil {
		log.Fatalln("GetAbsPath() error: ", err)
	}
	path = filepath.ToSlash(path)
	return StringsJoin(path, "/")
}

//根据文件路径创建对应的文件夹
func MkdirAllFile(path, appPath string) string {
	//filePath文件路径，fileName文件名
	filePath, fileName := filepath.Split(path)
	//日志文件的绝对路径
	filePath = StringsJoin(appPath, filePath)
	//panic日志文件的绝对路径+文件名
	path = StringsJoin(filePath, fileName)
	//根据文件路径创建对应的文件夹
	err := os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		log.Fatalln("panicFile failed to MkdirAll(): ", err)
	}
	return path
}

//转义还原，验证校验码
func Changebodymess(bufData []byte) (message []byte, err string) {
	message = bufData[1 : len(bufData)-1]

	message = bytes.ReplaceAll(message, []byte{0x7d, 0x02}, []byte{0x7e})
	message = bytes.ReplaceAll(message, []byte{0x7d, 0x01}, []byte{0x7d})

	var num = message[0]
	for i := 1; i < len(message)-1; i++ {
		num = num ^ message[i]
	}

	if num != message[len(message)-1] {
		fmt.Println("校验码错误")
		err = "校验码错误"
	}
	return message[:len(message)-1], err
}

//计算校验码转义封装
func ParaphraseMess(bufData []byte) (sendMes []byte) {
	sendMes = bufData
	num := sendMes[0]
	for i := 1; i < len(sendMes); i++ {
		num = num ^ sendMes[i]
	}
	sendMes = append(sendMes, num)

	sendMes = bytes.ReplaceAll(sendMes, []byte{0x7d}, []byte{0x7d, 0x01})
	sendMes = bytes.ReplaceAll(sendMes, []byte{0x7e}, []byte{0x7d, 0x02})

	sendMes = append([]byte{126}, sendMes...)
	sendMes = append(sendMes, 126)
	return
	//fmt.Println(sendMes)
}

//uint16转byte
func Uint16ToByte(i uint16) []byte {
	var buf = make([]byte, 2)
	binary.BigEndian.PutUint16(buf, i)
	return buf
}

//uint32转byte
func Uint32ToByte(i uint32) []byte {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, i)
	return buf
}

//int转8位byte
func Int64ToByte(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

//uint64转byte
func Uint64ToByte(i uint64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, i)
	return buf
}

//2进制到10进制
func TwoToInt64(two string) int64 {
	var middle []int64
	for _, v := range two {
		if string(v) == "1" {
			middle = append(middle, 1)
		} else if string(v) == "0" {
			middle = append(middle, 0)
		} else {
			return -1
		}
	}
	length := len(middle)
	if middle != nil {
		result := middle[length-1]
		for i := 1; i < length; i++ {
			result1 := middle[length-i-1]
			for j := 0; j < i; j++ {
				result1 *= 2
			}
			result += result1
		}
		return result
	}
	return -1
}

//uint32/uint16等转2进制再转为切片
func UintTo2byte(num uint64, length int) []byte {
	str2 := strconv.FormatUint(num, 2)
	for len(str2) < length {
		str2 = "0" + str2
	}
	return []byte(str2)
}

//byte倒序
func ReverseOrder(originalByte []byte) []byte {
	length := len(originalByte)
	var newByte []byte
	for i := 0; i < length; i++ {
		newByte = append(newByte, originalByte[length-1-i])
	}
	return newByte
}

//生成16进制字符串
func GetRandomString(lenth int) string {
	str := "0123456789abcdef"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < lenth; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

//float64转为度分秒
func Float64toDegree(i float64) (degree string) {
	degreeInt, minuteFloat := math.Modf(i)
	minute := minuteFloat * 60
	minuteInt, secondFloat := math.Modf(minute)
	second := secondFloat * 60
	degree = strconv.FormatFloat(degreeInt, 'f', 0, 64) + "°" + strconv.FormatFloat(minuteInt, 'f', 0, 64) + "‘" +
		strconv.FormatFloat(second, 'f', 1, 64) + "“"
	return
}

func GBKToUTF8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func UTF8ToGBK(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}

	return d, nil
}

//已验证均可以实现
func ConvertGBKToUTF8(str string) string {
	dec := mahonia.NewDecoder("GBK")
	return dec.ConvertString(str)
}

func ConvertUTF8ToGBK(str string) string {
	enc := mahonia.NewEncoder("GBK")
	return enc.ConvertString(str)
}