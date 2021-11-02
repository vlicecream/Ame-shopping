package utils

import (
	"go.uber.org/zap"
	"net"
)

func DynamicPort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		zap.L().Error("DynamicPort net.ResolveTCPAddr failed", zap.Error(err))
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		zap.L().Error("DynamicPort net.ListenTCP failed", zap.Error(err))
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
