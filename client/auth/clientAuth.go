package auth

import "context"

type Authentication struct {
	User   string
	Passwd string
}

func (a *Authentication) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"user": a.User, "passwd": a.Passwd}, nil
}

func (a *Authentication) RequireTransportSecurity() bool {

	return false
}
