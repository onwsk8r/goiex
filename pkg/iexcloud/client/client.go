package client

import (
	"github.com/onwsk8r/goiex/internal/api"
	"github.com/onwsk8r/goiex/pkg/iexcloud"
)

func Internal(token, version string) (iexcloud.Client, error) {
	return api.NewClient(token, version)
}
