package mongodb

type Client interface {
	Get(num int) int
}

func NewClient() Client {
	return client{}
}

func NewMockClient() Client {
	return clientMock{}
}
