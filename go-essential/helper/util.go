package utils

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/axgle/mahonia"
	log "myth/go-essential/log/logc"
	"net"
	"os"
	"reflect"
	"strings"
	"unicode"
)

func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

func GetElem(elem reflect.Value) reflect.Value {
	for elem.Kind() == reflect.Ptr || elem.Kind() == reflect.Interface {
		elem = elem.Elem()
	}
	return elem
}

// CamelCaseToUnderscore converts from camel case form to underscore separated form.
// Ex.: MyFunc => my_func
func CamelCaseToUnderscore(str string) string {
	var output []rune
	var segment []rune
	for _, r := range str {
		if !unicode.IsLower(r) && string(r) != "_" {
			output = addSegment(output, segment)
			segment = nil
		}
		segment = append(segment, unicode.ToLower(r))
	}
	output = addSegment(output, segment)
	return string(output)
}

func addSegment(inrune, segment []rune) []rune {
	if len(segment) == 0 {
		return inrune
	}
	if len(inrune) != 0 {
		inrune = append(inrune, '_')
	}
	inrune = append(inrune, segment...)
	return inrune
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func ConvertTcpAddress(addr net.Addr) net.Addr {
	tcpAddr, ok := addr.(*net.TCPAddr)
	if !ok {
		return addr
	}

	if !tcpAddr.IP.IsUnspecified() {
		return addr
	}

	host, err := os.Hostname()
	if err != nil {
		log.Errorf("ServiceRegistry convertTcpAddr Hostname fail, err:%s", err)
		return addr
	}

	ips := make([]net.IP, 0)

	interfaceAddrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Errorf("ServiceRegistry convertTcpAddr InterfaceAddrs fail, err:%s", err)
		return addr
	}

	log.Debugf("ServiceRegistry convertTcpAddr InterfaceAddrs host:[%s], addr:[%s]", host, interfaceAddrs)

	ips = make([]net.IP, 0)
	for _, addr := range interfaceAddrs {
		if ipAddr, ok := addr.(*net.IPNet); ok {
			ips = append(ips, ipAddr.IP)
		} else {
			log.Debugf("bot ip addr addr:[%s], type:[%s]", addr, reflect.TypeOf(addr))
		}
	}

	log.Debugf("ServiceRegistry convertTcpAddr LookupIP host:[%s], ip:[%s]", host, ips)

	found := false
	for _, ip := range ips {
		if ip.IsUnspecified() || ip.IsLoopback() {
			continue
		}

		if !strings.HasPrefix(ip.String(), "10.") &&
			!strings.HasPrefix(ip.String(), "192.168.") &&
			!strings.HasPrefix(ip.String(), "172.56.") {
			continue
		}

		tcpAddr.IP = ip
		found = true
		break
	}

	if !found {
		log.Errorf("ServiceRegistry convertTcpAddr LookupIP no suitable ip found, host:%s", ips)
		return addr
	}

	return tcpAddr
}
