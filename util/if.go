package util

func If(condition bool, trueValue, falseValue interface{}) interface{} {
	if condition {
		return trueValue
	}
	return falseValue
}

func Iff(condition bool, trueF func() interface{}, falseF func() interface{}) interface{} {
	if condition {
		return trueF()
	}
	return falseF()
}