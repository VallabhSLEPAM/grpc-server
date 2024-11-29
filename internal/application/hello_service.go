package application

type HelloService struct {
}

func (h *HelloService) GenerateHello(name string) string {
	return "Hello " + name
}
