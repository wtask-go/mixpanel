//+build tools
//go:generate go mod download
//go:generate go install "github.com/psampaz/go-mod-outdated"
//go:generate go install "github.com/deepmap/oapi-codegen/cmd/oapi-codegen"
//go:generate go mod tidy

package mixpanel

import (
	_ "github.com/deepmap/oapi-codegen/cmd/oapi-codegen" // required vendor
	_ "github.com/psampaz/go-mod-outdated"               // required vendor
)
