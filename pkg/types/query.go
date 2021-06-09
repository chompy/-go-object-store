package types

import (
	"strings"
	"time"
)

// Query is a search query.
type Query map[string]interface{}

func (q Query) checkKey(key string, qValue interface{}, oValue interface{}) bool {
	op := getOperator(key)
	switch qValue := qValue.(type) {
	case time.Time:
		{
			return timeMatch(op, qValue, oValue.(time.Time))
		}
	case int:
		{
			return numericMatch(op, float64(qValue), float64(oValue.(int)))
		}
	case int64:
		{
			return numericMatch(op, float64(qValue), float64(oValue.(int64)))
		}
	case int32:
		{
			return numericMatch(op, float64(qValue), float64(oValue.(int32)))
		}
	case float32:
		{
			return numericMatch(op, float64(qValue), float64(oValue.(float32)))
		}
	case float64:
		{
			return numericMatch(op, qValue, oValue.(float64))
		}
	case string:
		{
			return stringMatch(op, qValue, oValue.(string))
		}
	default:
		{
			switch op {
			case opNE:
				{
					return oValue != qValue
				}
			}
			return oValue == qValue
		}
	}
}

// Check checks if given object matches query.
func (q Query) Check(o *Object) bool {
	for k, v := range q {
		kk := strings.Split(k, "__")
		switch strings.ToLower(kk[0]) {
		case "uid":
			{
				if !q.checkKey(k, v, o.UID) {
					return false
				}
				break
			}
		case "created":
			{
				if !q.checkKey(k, v, o.Created) {
					return false
				}
				break
			}
		case "modified":
			{
				if !q.checkKey(k, v, o.Modified) {
					return false
				}
				break
			}
		case "type":
			{
				if !q.checkKey(k, v, o.Type) {
					return false
				}
				break
			}
		default:
			{
				if o.Data[kk[0]] == nil {
					break
				}
				if !q.checkKey(k, v, o.Data[kk[0]]) {
					return false
				}
				break
			}
		}
	}
	return true
}
