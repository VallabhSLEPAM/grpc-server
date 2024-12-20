package port

type HelloServicePort interface {
	GenerateHello(string) string
}

type BankServicePort interface {
	FindCurrentBalance(acct string) float64
}
