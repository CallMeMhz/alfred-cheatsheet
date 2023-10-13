package provider

import "github.com/callmemhz/godash/model"

type Provider interface {
	Search(namespace, keyword string) ([]model.Entry, error)
	Close()
}

type ProviderFactory interface {
	NewProvider() (Provider, error)
}
