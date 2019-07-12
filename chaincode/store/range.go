package store

type Range struct {
	First interface{}
	Last  interface{}
}

func R(first, last interface{}) *Range {
	return &Range{First: first, Last: last}
}
