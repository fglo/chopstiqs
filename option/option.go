package option

type OptInt struct {
	val   int
	isSet bool
}

var EmptyInt = OptInt{}

func Int(val int) OptInt {
	return OptInt{
		val:   val,
		isSet: true,
	}
}

func (i OptInt) Val() int {
	return i.val
}

func (i OptInt) IsEmpty() bool {
	return !i.isSet
}

func (i OptInt) IsSet() bool {
	return i.isSet
}

type OptFloat struct {
	val   float64
	isSet bool
}

var EmptyFloat = OptFloat{}

func Float(val float64) OptFloat {
	return OptFloat{
		val:   val,
		isSet: true,
	}
}

func (f OptFloat) Val() float64 {
	return f.val
}

func (f OptFloat) IsEmpty() bool {
	return !f.isSet
}

func (f OptFloat) IsSet() bool {
	return f.isSet
}
