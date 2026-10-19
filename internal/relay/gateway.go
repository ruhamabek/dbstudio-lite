package relay

import (
	"app-dev/internal/domain"
	"app-dev/internal/service"
	"context"
	"encoding/json"
	"fmt"
)


type HostGateway struct {
	inspector domain.SchemaInspector
	querier   domain.Querier
	guard     *QueryGuard
	runner    *service.Runner
}

func NewHostGateway(inspector domain.SchemaInspector, querier domain.Querier) *HostGateway{
	return &HostGateway{
		  inspector: inspector,
		  querier:   querier,
		  guard:     NewQueryGuard(),
		  runner:    service.NewRunner(querier),
	}
}

func (g *HostGateway) Handle(ctx context.Context, perm Permission, req Request) Response {
	switch req.Method {
	  case MethodListTables:
			tables, err := g.inspector.ListTables(ctx)
			if err != nil {
				return Response{ID: req.ID, Error: err.Error()}
			}
		return makeResponse(req.ID, tables)

	 case MethodGetColumns:
			cols, err := g.inspector.GetColumns(ctx, req.Schema, req.Table)
				if err != nil {
					return Response{ID: req.ID, Error: err.Error()}
				}
	    return makeResponse(req.ID, cols)

     case MethodGetPrimaryKeys:
			pks, err := g.inspector.GetPrimaryKeys(ctx, req.Schema, req.Table)
			if err != nil {
				return Response{ID: req.ID, Error: err.Error()}
			}
	    return makeResponse(req.ID, pks)

	case MethodExecuteQuery:
			if err := g.guard.Validate(perm, req.SQL); err != nil {
				return Response{ID: req.ID, Error: err.Error()}
			}
			result, err := g.runner.Execute(ctx, req.SQL)
			if err != nil {
				return Response{ID: req.ID, Error: err.Error()}
			}
	    return makeResponse(req.ID, result)
		
	default:
		return Response{ID: req.ID, Error: fmt.Sprintf("unknown method: %s", req.Method)}
	}     
}

func makeResponse(id string, result any) Response {
	data, err := json.Marshal(result)
	if err != nil {
		return Response{ID: id, Error: err.Error()}
	}
	return Response{ID: id, Result: data}
}