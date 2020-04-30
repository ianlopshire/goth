package goth

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

// Provider needs to be implemented for each 3rd party authentication provider
// e.g. Facebook, Twitter, etc...
type Provider interface {
	Name() string
	SetName(name string)
	BeginAuth(state string) (Session, error)
	UnmarshalSession(string) (Session, error)
	FetchUser(Session) (User, error)
	Debug(bool)
	RefreshToken(refreshToken string) (*oauth2.Token, error) //Get new access token based on the refresh token
	RefreshTokenAvailable() bool                             //Refresh token is provided by auth provider or not
}

const NoAuthUrlErrorMessage = "an AuthURL has not been set"

var providers = Providers{}

// Providers is a set of known/available providers.
type Providers map[string]Provider

// Use adds a list of providers to the set.
//
// Use can be called multiple times. If you pass the same provider more than once, the
// last will be used.
func (p Providers) Use(viders ...Provider) {
	for _, provider := range viders {
		p[provider.Name()] = provider
	}
}

// GetProvider returns a provider by name. If the provider has not been added to the set,
// an error will be returned.
func (p Providers) Get(name string) (Provider, error) {
	provider := providers[name]
	if provider == nil {
		return nil, fmt.Errorf("no provider for %s exists", name)
	}
	return provider, nil
}

func (p *Providers) Clear() {
	*p = Providers{}
}

// UseProviders adds a list of available providers for use with Goth.
// Can be called multiple times. If you pass the same provider more
// than once, the last will be used.
func UseProviders(viders ...Provider) {
	providers.Use(viders...)
}

// GetProviders returns a list of all the providers currently in use.
func GetProviders() Providers {
	return providers
}

// GetProvider returns a previously created provider. If Goth has not
// been told to use the named provider it will return an error.
func GetProvider(name string) (Provider, error) {
	return providers.Get(name)
}

// ClearProviders will remove all providers currently in use.
// This is useful, mostly, for testing purposes.
func ClearProviders() {
	providers.Clear()
}

// ContextForClient provides a context for use with oauth2.
func ContextForClient(h *http.Client) context.Context {
	if h == nil {
		return oauth2.NoContext
	}
	return context.WithValue(oauth2.NoContext, oauth2.HTTPClient, h)
}

// HTTPClientWithFallBack to be used in all fetch operations.
func HTTPClientWithFallBack(h *http.Client) *http.Client {
	if h != nil {
		return h
	}
	return http.DefaultClient
}
