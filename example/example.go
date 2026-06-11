package main

import (
	"github.com/Meduzz/commando"
	"github.com/Meduzz/dsl"
	"github.com/Meduzz/dsl/app"
	"github.com/Meduzz/dsl/endpoint"
	"github.com/Meduzz/dsl/service"
	"github.com/Meduzz/generator/pkg/inputs"
	"github.com/Meduzz/quickapi/model"
)

type (
	Data struct{}

	Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
)

var (
	_       model.Entity = Data{}
	testApp              = dsl.NewApp("app", func(ab app.AppBuilder) {
		ab.AddService("service", func(sb service.ServiceBuilder) {
			sb.SetDescription("Not used?")
			sb.AddEndpoint("fetch", func(eb endpoint.EndpointBuilder) {
				eb.GET("/api/service/:id")
				eb.Path("id")
				eb.Query("query")
				eb.RequestHeader("Authentication")
			})
			sb.AddEndpoint("store", func(eb endpoint.EndpointBuilder) {
				eb.POST("/api/service")
				eb.RequestHeader("Authentication")
				eb.RequestBody("application/json", &Person{})
			})
		})
	})
)

func main() {
	// name := entity()
	// entities(testApp.Name)
	name := dslService()
	// dslApp()

	err := commando.Execute(name)
	// err := commando.Execute(testApp.Name)

	if err != nil {
		panic(err)
	}
}

func entity() string {
	i := inputs.FromQuickapi(Data{})
	return i.ToCLI("http://localhost:8000")
}

func entities(name string) {
	i := inputs.FromQuickapiEntities(name, Data{})
	i.ToCLI("http://localhost:8000")
}

func dslService() string {
	i := inputs.FromDSLService(testApp.Services[0])
	return i.ToCLI("http://localhost:8000")
}

func dslApp() {
	i := inputs.FromDSLApp(testApp)
	i.ToCLI("http://localhost:8000")
}

func (Data) Name() string {
	return "data"
}

func (Data) Create() any {
	return &Data{}
}

func (Data) CreateArray() any {
	return make([]*Data, 0)
}
