package main

import (
	"github.com/gogo/protobuf/vanity/command"
	"github.com/influxdata/yarpc/cmd/protoc-gen-yarpc/generator"
	"github.com/influxdata/yarpc/cmd/protoc-gen-yarpc/yarpc"
)

func main() {
	g := generator.New()
	g.Request = command.Read()
	g.Suffix = ".yarpc.go"
	g.CommandLineParameters(g.Request.GetParameter())

	// Create a wrapped version of the Descriptors and EnumDescriptors that
	// point to the file that defines them.
	g.WrapTypes()

	g.SetPackageNames()
	g.BuildTypeNameMap()

	p := &yarpc.Plugin{}
	g.GeneratePlugin(p)

	command.Write(g.Response)
}
