package mongodb

type clientMock struct {
	// mock logic, nothing real will happen here.
}

func (clientMock) Get(num int) int {
	return num * 3
}
