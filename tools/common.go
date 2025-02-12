package tools

import (
	"errors"
	"fmt"
	"net"
	"reflect"
)

func GetLocalIP() (ipv4 string) {
	var (
		addrs   []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet
		isNetIp bool
		err     error
	)

	addrs, err = net.InterfaceAddrs()
	HandelError("GetLocalIP error ", err)
	for _, addr = range addrs {
		if ipNet, isNetIp = addr.(*net.IPNet); isNetIp && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}

	err = errors.New("not found ip from get local ip")
	HandelError("GetLocalIP error ", err)
	return
}

func CheckStructNil(value interface{}) bool {
	v := reflect.ValueOf(value)

	if v.Kind() != reflect.Struct {
		return false
	}
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.Kind() != reflect.Ptr {
			continue
		}
		if f.IsNil() {
			return false
		}
	}
	return true
}

func FormatAddr(port int) string {
	localP := GetLocalIP()
	return fmt.Sprintf("%s:%d", localP, port)
}
