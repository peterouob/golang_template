package tools

import (
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"log"
	"net"
	"net/http"
	"reflect"
	"strings"
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

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		log.Printf("Gateway received: method=%s, url=%s, content-length=%d", r.Method, r.URL, r.ContentLength)
		next.ServeHTTP(w, r)
	})
}

func Matcher(key string) (string, bool) {
	switch strings.ToLower(key) {
	case "authorization":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}
