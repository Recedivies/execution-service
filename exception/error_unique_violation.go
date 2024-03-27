package exception

import (
	"github.com/lib/pq"
)

func IsUniqueViolation(err error) (string, bool) {
	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Code.Name() {
		case "unique_violation":
			return pqErr.Constraint, true
		}
	}
	return "", false
}
