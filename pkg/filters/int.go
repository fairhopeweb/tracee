package filters

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/aquasecurity/tracee/pkg/logger"
	"golang.org/x/exp/constraints"
)

const (
	maxNotSetInt int64 = math.MinInt64
	minNotSetInt int64 = math.MaxInt64
)

type IntFilter[T constraints.Signed] struct {
	equal    map[int64]bool
	notEqual map[int64]bool
	min      int64
	max      int64
	is32Bit  bool
	enabled  bool
}

// TODO: Add int16 and int8 filters?

func NewIntFilter() *IntFilter[int64] {
	return newIntFilter[int64](false)
}

func NewInt32Filter() *IntFilter[int32] {
	return newIntFilter[int32](true)
}

func newIntFilter[T constraints.Signed](is32Bit bool) *IntFilter[T] {
	filter := &IntFilter[T]{
		equal:    map[int64]bool{},
		notEqual: map[int64]bool{},
		min:      minNotSetInt,
		max:      maxNotSetInt,
		is32Bit:  is32Bit,
	}

	return filter
}

func (f *IntFilter[T]) Enable() {
	f.enabled = true
}

func (f *IntFilter[T]) Disable() {
	f.enabled = false
}

func (f *IntFilter[T]) Enabled() bool {
	return f.enabled
}

func (f *IntFilter[T]) Minimum() int64 {
	return f.min
}

func (f *IntFilter[T]) Maximum() int64 {
	return f.max
}

func (f *IntFilter[T]) Filter(val interface{}) bool {
	filterable, ok := val.(T)
	if !ok {
		return false
	}
	return f.filter(filterable)
}

// priority goes by (from most significant):
// 1. equality
// 2. greater
// 3. lesser
// 4. non equality
func (f *IntFilter[T]) filter(val T) bool {
	compVal := int64(val)
	result := !f.enabled || f.equal[compVal] || compVal > f.min || compVal < f.max
	if !result && f.notEqual[compVal] {
		return false
	}
	return result
}

func (f *IntFilter[T]) validate(val int64) bool {
	if f.is32Bit {
		return val <= math.MaxInt32 && val >= math.MinInt32
	}
	return true
}

func (f *IntFilter[T]) addEqual(val int64) {
	f.equal[val] = true
}

func (f *IntFilter[T]) addNotEqual(val int64) {
	f.notEqual[val] = true
}

func (f *IntFilter[T]) addLesserThan(val int64) {
	// we want to have the highest max input
	if val > f.max {
		f.max = val
	}
}

func (f *IntFilter[T]) addGreaterThan(val int64) {
	// we want to have the lowest min input
	if val < f.min {
		f.min = val
	}
}

func (f *IntFilter[T]) add(val int64, operator Operator) error {
	if !f.validate(val) {
		return InvalidValue(fmt.Sprint(val))
	}
	switch operator {
	case Equal:
		f.addEqual(val)
	case NotEqual:
		f.addNotEqual(val)
	case Lower:
		f.addLesserThan(val)
	case Greater:
		f.addGreaterThan(val)
	case LowerEqual:
		f.addEqual(val)
		f.addLesserThan(val)
	case GreaterEqual:
		f.addEqual(val)
		f.addGreaterThan(val)
	}
	return nil
}

func (filter *IntFilter[T]) Parse(operatorAndValues string) error {
	if len(operatorAndValues) < 2 {
		return InvalidExpression(operatorAndValues)
	}
	valuesString := string(operatorAndValues[1:])
	operatorString := string(operatorAndValues[0])

	// check for !=
	if operatorString == "!" {
		if len(operatorAndValues) < 3 {
			return InvalidExpression(operatorAndValues)
		}
		operatorString = operatorAndValues[0:2]
		valuesString = operatorAndValues[2:]
	}

	// check for >= and <=
	if (operatorString == ">" || operatorString == "<") && operatorAndValues[1] == '=' {
		if len(operatorAndValues) < 3 {
			return InvalidExpression(operatorAndValues)
		}
		operatorString = operatorAndValues[0:2]
		valuesString = operatorAndValues[2:]
	}

	values := strings.Split(valuesString, ",")
	operator := stringToOperator(operatorString)

	for _, val := range values {
		valInt, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return InvalidValue(val)
		}
		err = filter.add(valInt, operator)
		if err != nil {
			return logger.ErrorFunc(err)
		}
	}

	filter.Enable()

	return nil
}
