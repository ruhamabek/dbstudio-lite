package relay

import (
	"errors"
	"strings"
)


type Permission string


const (
	PermissionView   Permission = "view"
	PermissionRead   Permission = "read"
	PermissionUpdate Permission = "update"
)


var (
	ErrViewOnlyForbidden = errors.New("view permission does not allow executing SQL queries")
	ErrMutationForbidden = errors.New("mutating queries are forbidden under read permission")
	ErrDDLForbidden      = errors.New("destructive DDL operations are forbidden under update permission")
)


type QueryGuard struct {}

func NewQueryGuard() *QueryGuard {
	 return &QueryGuard{}
}

func (g *QueryGuard) Validate(perm Permission, sql string) error {

	if perm == PermissionView {
		return ErrViewOnlyForbidden
	}

	trimmed := strings.TrimSpace(sql)
	if trimmed == "" {
		return nil
	}

	fields := strings.Fields(trimmed)
	firstWord := strings.ToUpper(fields[0])

	switch perm {
	   case PermissionRead:
		   switch firstWord {
			    case "SELECT", "WITH", "EXPLAIN":
			        return nil
		    default:
			       return ErrMutationForbidden
		   }
		case PermissionUpdate:
			 switch firstWord {
				case "DROP", "TRUNCATE", "ALTER":
					return ErrDDLForbidden
				default:
					return nil
			 }
		default:
			return errors.New("invalid permission")
	}
}
