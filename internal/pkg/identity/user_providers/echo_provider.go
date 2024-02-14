package user_providers

type EchoProvider struct {
}

func NewEchoProvider() *EchoProvider {
	return &EchoProvider{}
}

func (p *EchoProvider) GetTokenByUser(user string) (string, error) {
	return user, nil
}
