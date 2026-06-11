package actions

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Meduzz/commando"
	"github.com/Meduzz/commando/builder"
	"github.com/Meduzz/commando/flags"
	"github.com/Meduzz/commando/model"
	myflags "github.com/Meduzz/generator/pkg/flags"
	"github.com/Meduzz/generator/pkg/future"
	"github.com/Meduzz/helper/fp/slice"
	"github.com/Meduzz/helper/http/client"
	"github.com/Meduzz/helper/http/herror"
	"github.com/Meduzz/helper/utilz"
	"github.com/spf13/cobra"
)

func CreateCommand(t *API, handlerRef string) *model.Command {
	return commando.CommandBuilder(t.Name, func(cb builder.CommandBuilder) {
		cb.HandlerRef(handlerRef)

		urlFlags := slice.Filter(t.Params, func(p *Param) bool {
			return p.Source == Flag && (p.Kind == Path || p.Kind == Query)
		})

		slice.ForEach(urlFlags, func(p *Param) {
			var f *model.Flag

			switch p.Subkind {
			case String:
				value := ""

				if p.Value != nil {
					value = p.Value.(string)
				}

				f = flags.StringFlag(p.Name, value, "")
			case Int:
				value := 0

				if p.Value != nil {
					value = p.Value.(int)
				}

				f = flags.IntFlag(p.Name, value, "")
			case Map:
				f = myflags.MapFlag(p.Name)
			case Array:
				f = myflags.Array(p.Name)
			}

			if f != nil {
				cb.Flag(f)
			}
		})
	})
}

// TODO a lot of these errors will need some love to make sense to the token munchers.
func CreateHandler(baseUrl string, t *API) model.Handler {
	return func(c *cobra.Command, s []string) error {
		// set flag params
		flagParams := slice.Filter(t.Params, func(p *Param) bool {
			return p.Source == Flag
		})

		errorz := slice.Map(flagParams, func(p *Param) error {
			switch p.Subkind {
			case String:
				value, err := c.Flags().GetString(p.Name)

				if err != nil {
					return err
				}

				p.Value = value
			case Int:
				value, err := c.Flags().GetInt(p.Name)

				if err != nil {
					return err
				}

				p.Value = value
			case Map:
				value, err := c.Flags().GetStringToString(p.Name)

				if err != nil {
					return err
				}

				p.Value = value
			case Array:
				value, err := c.Flags().GetStringSlice(p.Name)

				if err != nil {
					return err
				}

				p.Value = value
			}
			return nil
		})

		err := slice.Fold(errorz, nil, func(in, agg error) error {
			if agg != nil {
				return errors.Join(agg, in)
			}

			return in
		})

		if err != nil {
			return err
		}

		// set env params
		envParams := slice.Filter(t.Params, func(p *Param) bool {
			return p.Source == Env
		})

		slice.ForEach(envParams, func(p *Param) {
			value := utilz.Env(p.Name, "")
			if value != "" {
				p.Value = value
			}
		})

		// prepare the data
		pathParams := slice.Filter(t.Params, func(p *Param) bool {
			return p.Kind == Path
		})
		pathValues := slice.Map(pathParams, func(p *Param) any {
			return p.Value
		})

		queryParams := slice.Filter(t.Params, func(p *Param) bool {
			return p.Kind == Query
		})

		query := strings.Join(slice.Filter(slice.Map(queryParams, func(p *Param) string {
			return p.Query()
		}), func(it string) bool {
			return it != ""
		}), "&")

		// set body if present
		var body []byte
		if t.Method != "GET" {
			cb := future.Future(io.ReadAll)
			body, err = cb(os.Stdin)

			if err != nil {
				if err != future.ErrTimeout {
					return err
				}
			}
		}

		// build request
		format := t.Path
		path := fmt.Sprintf(format, pathValues...)

		if query != "" {
			path = fmt.Sprintf("%s?%s", path, query)
		}

		headers := slice.Filter(t.Params, func(p *Param) bool {
			return p.Kind == Header
		})

		req, err := client.NewRequest(t.Method, fmt.Sprintf("%s%s", baseUrl, path), body, "") // TODO content type

		if err != nil {
			return err
		}

		slice.ForEach(headers, func(p *Param) {
			req.Header(p.Name, p.Header())
		})

		res, err := req.DoDefault()

		if err != nil {
			return err
		}

		err = herror.IsError(res.Code())

		if err != nil {
			return err
		}

		returnBody, err := res.AsText()

		if err != nil {
			return err
		}

		println(returnBody) // TODO lets evaluate this one...

		return nil
	}
}
