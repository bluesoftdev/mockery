package httpMock

func Method(method string, configFunc func()) {
	Case(RequestKeyStringEquals(method), configFunc)
}
