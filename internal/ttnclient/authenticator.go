package ttnclient

import (
	"fmt"
	"net/http"
)

type Authenticator interface {
	Authenticate(request *http.Request) error
}

type ApiKeyAuthenticator struct {
	ApiKey string
}

func (a ApiKeyAuthenticator) Authenticate(request *http.Request) error {
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.ApiKey))
	return nil
}
