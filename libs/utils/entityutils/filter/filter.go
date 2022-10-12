package filter

type Operation int64

const (
	Eq Operation = iota
	In
	NotEq
)

type Param struct {
	Op    Operation
	Value interface{}
}

type Filter struct {
	Params map[string]Param
}

func NewFilter() Filter {
	return Filter{
		Params: make(map[string]Param),
	}
}

func (f *Filter) Add(
	field string,
	Op Operation,
	value interface{},
) {
	f.Params[field] = Param{
		Op:    Op,
		Value: value,
	}
}

func (f *Filter) Match(
	candidate map[string]interface{},
) bool {
	for field := range f.Params {
		cv, ok := candidate[field]
		if !ok {
			return false
		}
		switch f.Params[field].Op {
		case Eq:
			if f.Params[field].Value != cv {
				return false
			}
		case NotEq:
			if f.Params[field].Value == cv {
				return false
			}
		case In:
			match := false
			sv, ok := f.Params[field].Value.([]interface{})
			if !ok {
				return false
			}
			for _, v := range sv {
				if v == cv {
					match = true
				}
			}
			if !match {
				return false
			}
		}
	}

	return true
}
