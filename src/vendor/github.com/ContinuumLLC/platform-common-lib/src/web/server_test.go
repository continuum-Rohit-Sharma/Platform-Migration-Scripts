package web

import "testing"

func TestSomething(t *testing.T) {
	var f ServerFactory
	f = ServerFactoryImpl{}
	srv := f.GetServer(&ServerConfig{})
	_, ok := srv.(*muxConfig)
	if !ok {
		t.Error("Server not muxServer")
	}
}
