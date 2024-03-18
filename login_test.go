package main

import (
	"testing"
	"wxcloudrun-golang/service"
)

func TestLogin(t *testing.T) {
	t.Log("TestLogin")
	if ret, err := service.Login(); err != nil {
		t.Error(err)
	} else {
		t.Log(ret)
	}
}
