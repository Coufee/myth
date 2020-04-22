package utils

import (
	"path"
	"reflect"
	"regexp"
	"strings"
	"unicode/utf8"
)

//检测接口控制
func VerifyNil(i interface{}) bool {
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}

	return false
}

//检验邮箱
func VerifyEmail(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

//校验手机号
func VerifyTel(tel string) bool {
	pattern := "^1([123456789][0-9])\\d{8}$"
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(tel)
}

//校验qq
func VerifyQq(qq string) bool {
	pattern := "^[1-9][0-9]{4,14}$"
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(qq)
}

//检验后缀名
func VerifySuffix(filePath, suffix string) bool {
	fileNameSuffix := path.Base(filePath)
	fileSuffix := strings.ToLower(fileNameSuffix)
	fileSuffix = strings.Trim(fileSuffix, ".")
	if fileSuffix == suffix {
		return true
	}

	return false
}

//检验后缀名
func VerifyMapSuffix(filePath string, suffixMap map[string]interface{}) bool {
	fileNameSuffix := path.Base(filePath)
	fileSuffix := strings.ToLower(fileNameSuffix)
	fileSuffix = strings.Trim(fileSuffix, ".")
	if _, ok := suffixMap[fileSuffix]; ok {
		return true
	}

	return false
}

//正则校验是否正确url
func VerifyDomain(str string) bool {
	reg := `^((ht|f)tps?):\/\/[\w\-]+(\.[\w\-]+)+([\w\-\.\{\},@?^=%&:\/~\+#]*[\w\-\@?^=%&\/~\+#])?$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(str)
}

//校验域名传递值数组
func VerifyDomainList(arr []string) bool {
	for _, domain := range arr {
		if !VerifyDomain(domain) {
			return false
		}
	}

	return true
}

//校验字符串长度
func VerifyStringLength(str string, minLen, maxLen int) bool {
	strLen := utf8.RuneCountInString(str)
	if strLen > maxLen || strLen < minLen {
		return false
	}

	return true
}

//校验Ip正确性
func VerifyIp(ip string) bool {
	reg := `^(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|[1-9])\.(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)\.(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)\.(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(ip)

}

//校验设备号
func VerifyDevices(device string) bool {
	reg := `^[0-9a-zA-Z_-]{1,}$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(device)
}

//校验iMei 15位 数字字母组合 如: 123456789abcdef
func VerifyIMei(iMei string) bool {
	reg := regexp.MustCompile(`^[0-9a-zA-z]{15}$`)
	return reg.MatchString(iMei)
}

//校验idfa 8位–4位–4位–4位–12位, 数字字母组合
func VerifyIdfa(idfa string) bool {
	reg := regexp.MustCompile(`^[0-9a-zA-z]{8}[-][0-9a-zA-z]{4}[-][0-9a-zA-z]{4}[-][0-9a-zA-z]{4}[-][0-9a-zA-z]{12}$`)
	return reg.MatchString(idfa)
}

//校验md5 16或者32
func VerifyMd516Or32(md string) bool {
	bl16, _ := regexp.MatchString(`^[0-9a-zA-Z]{16}$`, md)
	bl32, _ := regexp.MatchString(`^[0-9a-zA-Z]{32}$`, md)

	if bl16 || bl32 {
		return true
	}

	return false
}
