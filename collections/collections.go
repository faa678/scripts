package collections

type Collection interface {
	Peek() interface{}
	Offer(ele interface{}) error
	Poll() (interface{}, error)
	SizeOf() int64
	CapOf() int64
	IsEmpty() bool
}
