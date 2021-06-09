package types

import "strings"

const (
	opEQ = "__eq"
	opNE = "__ne"
	opGT = "__gt"
	opLT = "__lt"
)

func getOperator(key string) string {
	if strings.HasSuffix(key, opGT) {
		return opGT
	} else if strings.HasSuffix(key, opLT) {
		return opLT
	} else if strings.HasSuffix(key, opNE) {
		return opNE
	}
	return opEQ
}
