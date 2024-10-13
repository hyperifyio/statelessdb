package z3b

import "statelessdb/internal/helpers"

func compareSet(a, b []byte) bool {
	return helpers.CompareSlices(a, b)
}
