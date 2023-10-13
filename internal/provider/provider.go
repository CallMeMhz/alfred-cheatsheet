package provider

import "github.com/callmemhz/alfred-cheatsheet/model"

type Provider interface {
	Search(namespace, keyword string) ([]model.Entry, error)
	Close()
}

type ProviderFactory interface {
	NewProvider() (Provider, error)
}
