package port

type HelloServicePort interface {
	GenerateHello(string) string
}
