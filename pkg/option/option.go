package option

type OptInt struct {
	val    int
	hasVal bool
}

var EmptyInt = OptInt{}

func Int(val int) OptInt {
	return OptInt{
		val:    val,
		hasVal: true,
	}
}

func (i OptInt) Val() int {
	return i.val
}

func (i OptInt) IsEmpty() bool {
	return !i.hasVal
}

func (i OptInt) HasVal() bool {
	return i.hasVal
}

type OptFloat struct {
	val    float64
	hasVal bool
}

var EmptyFloat = OptFloat{}

func Float(val float64) OptFloat {
	return OptFloat{
		val:    val,
		hasVal: true,
	}
}

func (f OptFloat) Val() float64 {
	return f.val
}

func (f OptFloat) IsEmpty() bool {
	return !f.hasVal
}

func (f OptFloat) HasVal() bool {
	return f.hasVal
}
