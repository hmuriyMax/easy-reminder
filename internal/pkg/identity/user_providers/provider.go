package user_providers

type Provider interface {
	GetTokenByUser(user string) (string, error)
}
