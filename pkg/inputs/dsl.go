package inputs

import (
	"fmt"
	"strings"

	"github.com/Meduzz/commando"
	"github.com/Meduzz/commando/builder"
	"github.com/Meduzz/commando/registry"
	"github.com/Meduzz/dsl/app"
	"github.com/Meduzz/dsl/endpoint"
	"github.com/Meduzz/dsl/service"
	serviceref "github.com/Meduzz/dsl/serviceRef"
	"github.com/Meduzz/generator/pkg/actions"
	"github.com/Meduzz/helper/fp/slice"
)

func FromDSLApp(app *app.App) *DslApp {
	return &DslApp{
		app: app,
	}
}

func FromDSLService(service *service.Service) *DslService {
	return &DslService{
		service: service,
	}
}

func (d *DslService) ToCLI(baseUrl string) string {
	httpOnly := slice.Filter(d.service.Endpoints, func(e *endpoint.Endpoint) bool {
		return e.Route.Kind == endpoint.HttpKind
	})

	mapped := slice.Map(httpOnly, endpoint2Api)

	slice.ForEach(mapped, func(a *actions.API) {
		cmd := actions.CreateCommand(a, a.Name)
		handler := actions.CreateHandler(baseUrl, a)

		registry.RegisterHandler(a.Name, handler)
		registry.RegisterCommand(cmd)
	})

	return d.service.Name
}

func (a *DslApp) ToCLI(baseUrl string, exclude ...serviceref.ServiceRef) {
	validExclusions := slice.Filter(exclude, func(s serviceref.ServiceRef) bool {
		it, ok := s.App()

		return ok && it != a.app.Name
	})

	validServices := slice.Filter(a.app.Services, func(s *service.Service) bool {
		return !slice.Fold(validExclusions, false, func(ref serviceref.ServiceRef, agg bool) bool {
			if agg {
				return agg
			}

			it, ok := ref.Service()
			return ok && it != s.Name
		})
	})

	slice.ForEach(validServices, func(s *service.Service) {
		svcCmd := commando.CommandBuilder(s.Name, func(cb builder.CommandBuilder) {
			cb.Description(s.Description)
		})

		httpOnly := slice.Filter(s.Endpoints, func(e *endpoint.Endpoint) bool {
			return e.Route.Kind == endpoint.HttpKind
		})

		mapped := slice.Map(httpOnly, endpoint2Api)

		slice.ForEach(mapped, func(a *actions.API) {
			subName := fmt.Sprintf("%s-%s", s.Name, a.Name)
			cmd := actions.CreateCommand(a, subName)
			handler := actions.CreateHandler(baseUrl, a)

			svcCmd.Children = append(svcCmd.Children, cmd)

			registry.RegisterHandler(subName, handler)
		})

		registry.RegisterCommand(svcCmd)
	})
}

func endpoint2Api(e *endpoint.Endpoint) *actions.API {
	out := &actions.API{}

	out.Name = e.Name
	out.Method = e.Route.Method

	pathParams := slice.Filter(e.Route.In, func(p *endpoint.Param) bool {
		return p.Kind == endpoint.PathKind
	})
	out.Path = slice.Fold(pathParams, e.Route.Path, func(p *endpoint.Param, agg string) string {
		return strings.ReplaceAll(agg, fmt.Sprintf(":%s", p.Name), "%s")
	})

	notBody := slice.Filter(e.Route.In, func(p *endpoint.Param) bool {
		return p.Kind != endpoint.BodyKind
	})

	out.Params = slice.Map(notBody, func(i *endpoint.Param) *actions.Param {
		p := &actions.Param{}

		switch i.Kind {
		case endpoint.PathKind:
			p.Kind = actions.Path
			p.Source = actions.Flag
		case endpoint.HeaderKind:
			p.Kind = actions.Header
			p.Source = actions.Env
		case endpoint.QueryKind:
			p.Kind = actions.Query
			p.Source = actions.Flag
		}

		p.Name = i.Name
		p.Subkind = actions.String

		return p
	})

	bodyParam := slice.Head(slice.Filter(e.Route.In, func(p *endpoint.Param) bool {
		return p.Kind == endpoint.BodyKind
	}))

	if bodyParam != nil && bodyParam.ContentType != "" {
		out.ContentType = bodyParam.ContentType
	}

	return out
}
