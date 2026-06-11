package inputs

import (
	"github.com/Meduzz/dsl/app"
	"github.com/Meduzz/dsl/service"
	"github.com/Meduzz/quickapi/model"
)

type (
	DslApp struct {
		app *app.App
	}

	DslService struct {
		service *service.Service
	}

	QuickapiEntity struct {
		entity model.Entity
	}

	QuickapiEntities struct {
		name     string
		entities []model.Entity
	}
)
