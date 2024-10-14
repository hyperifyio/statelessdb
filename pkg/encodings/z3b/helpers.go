package z3b

import (
	"github.com/hyperifyio/statelessdb/pkg/helpers"
)

func compareSet(a, b []byte) bool {
	return helpers.CompareSlices(a, b)
}
