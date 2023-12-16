package service

import "helptools/internal/service/proxy"

func Bind() []interface{} {
	return []interface{}{
		&proxy.Proxy{},
	}
}
