//+build tools
//go:generate go mod download
//go:generate go install "github.com/psampaz/go-mod-outdated"
//go:generate go mod tidy

package mixpanel

import (
	_ "github.com/psampaz/go-mod-outdated" // required vendor
)
