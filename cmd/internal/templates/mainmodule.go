package templates

var MainModuleTmpl = EditableHeaderTmpl + `package main

import (
	"github.com/unionj-cloud/go-doudou/v2/framework/grpcx"
	"github.com/unionj-cloud/go-doudou/v2/framework/plugin"
	"github.com/unionj-cloud/go-doudou/v2/framework/rest"
	"github.com/unionj-cloud/toolkit/zlogger"
	_ "{{.PluginPackage}}"
)

func main() {
	srv := rest.NewRestServer()
	grpcServer := grpcx.NewEmptyGrpcServer()
	plugins := plugin.GetServicePlugins()
	for _, key := range plugins.Keys() {
		value, _ := plugins.Get(key)
		value.Initialize(srv, grpcServer, nil)
	}
	defer func() {
		if r := recover(); r != nil {
			zlogger.Info().Msgf("Recovered. Error: %v\n", r)
		}
		for _, key := range plugins.Keys() {
			value, _ := plugins.Get(key)
			value.Close()
		}
	}()
	go func() {
		grpcServer.Run()
	}()
	srv.Run()
}
`
