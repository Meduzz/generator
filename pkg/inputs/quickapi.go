package inputs

import (
	"fmt"

	"github.com/Meduzz/commando"
	"github.com/Meduzz/commando/builder"
	"github.com/Meduzz/commando/registry"
	"github.com/Meduzz/generator/pkg/actions"
	myflags "github.com/Meduzz/generator/pkg/flags"
	"github.com/Meduzz/helper/fp/slice"
	"github.com/Meduzz/quickapi/model"
)

func FromQuickapi(entity model.Entity) *QuickapiEntity {
	return &QuickapiEntity{
		entity: entity,
	}
}

func FromQuickapiEntities(name string, entities ...model.Entity) *QuickapiEntities {
	return &QuickapiEntities{
		name:     name,
		entities: entities,
	}
}

func (q *QuickapiEntity) ToCLI(baseUrl string) string {
	name := q.entity.Name()

	// TODO move somewhere sensible
	registry.RegisterVisitor(myflags.MapFlagVisitor())

	slice.ForEach(templates, func(t *actions.API) {
		t.Params = slice.Map(t.Params, func(p *actions.Param) *actions.Param {
			if p.Name == "entity" {
				p.Value = name
			}

			return p
		})
		cmd := actions.CreateCommand(t, t.Name)
		handler := actions.CreateHandler(baseUrl, t)

		registry.RegisterHandler(t.Name, handler)
		registry.RegisterCommand(cmd)
	})

	return name
}

func (q *QuickapiEntities) ToCLI(baseUrl string) {
	// TODO move somewhere sensible
	registry.RegisterVisitor(myflags.MapFlagVisitor())

	slice.ForEach(q.entities, func(e model.Entity) {
		name := e.Name()
		cmd := commando.CommandBuilder(name, func(cb builder.CommandBuilder) {
			cb.Description(fmt.Sprintf("Actions fo %s entity", name))
		})

		slice.ForEach(templates, func(t *actions.API) {
			subName := fmt.Sprintf("%s-%s", name, t.Name)
			t.Params = slice.Map(t.Params, func(p *actions.Param) *actions.Param {
				if p.Name == "entity" {
					p.Value = name
				}

				return p
			})
			subCmd := actions.CreateCommand(t, subName)
			subCmdHandler := actions.CreateHandler(baseUrl, t)

			cmd.Children = append(cmd.Children, subCmd)

			registry.RegisterHandler(subName, subCmdHandler)
		})

		registry.RegisterCommand(cmd)
	})
}

// TODO descriptions, both api and flags
var templates = []*actions.API{
	{
		Name:        "create",
		Method:      "POST",
		Path:        "/api/%s/",
		ContentType: "application/json",
		Params: []*actions.Param{
			{
				Name:    "entity",
				Kind:    actions.Path,
				Subkind: actions.String,
				Source:  actions.Flag,
			},
		},
	}, {
		Name:        "read",
		Method:      "GET",
		Path:        "/api/%s/%d",
		ContentType: "application/json",
		Params: []*actions.Param{
			{
				Name:    "entity",
				Kind:    actions.Path,
				Subkind: actions.String,
				Source:  actions.Flag,
			}, {
				Name:    "id",
				Kind:    actions.Path,
				Subkind: actions.Int,
				Source:  actions.Flag,
			},
		},
	}, {
		Name:        "update",
		Method:      "PUT",
		Path:        "/api/%s/%d",
		ContentType: "application/json",
		Params: []*actions.Param{
			{
				Name:    "entity",
				Kind:    actions.Path,
				Subkind: actions.String,
				Source:  actions.Flag,
			}, {
				Name:    "id",
				Kind:    actions.Path,
				Subkind: actions.Int,
				Source:  actions.Flag,
			},
		},
	}, {
		Name:        "delete",
		Method:      "DELETE",
		Path:        "/api/%s/%d",
		ContentType: "application/json",
		Params: []*actions.Param{
			{
				Name:    "entity",
				Kind:    actions.Path,
				Subkind: actions.String,
				Source:  actions.Flag,
			}, {
				Name:    "id",
				Kind:    actions.Path,
				Subkind: actions.Int,
				Source:  actions.Flag,
			},
		},
	}, {
		Name:        "list",
		Method:      "GET",
		Path:        "/api/%s/",
		ContentType: "application/json",
		Params: []*actions.Param{
			{
				Name:    "entity",
				Kind:    actions.Path,
				Subkind: actions.String,
				Source:  actions.Flag,
			}, {
				Name:    "skip",
				Kind:    actions.Query,
				Subkind: actions.Int,
				Value:   0,
				Source:  actions.Flag,
			}, {
				Name:    "take",
				Kind:    actions.Query,
				Subkind: actions.Int,
				Value:   25,
				Source:  actions.Flag,
			}, {
				Name:    "where",
				Kind:    actions.Query,
				Subkind: actions.Map,
				Source:  actions.Flag,
			}, {
				Name:    "sort",
				Kind:    actions.Query,
				Subkind: actions.Map,
				Source:  actions.Flag,
			},
		},
	},
}
