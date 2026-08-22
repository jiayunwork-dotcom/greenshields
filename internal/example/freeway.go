// Package example loads and decodes Greenshields example files. It also embeds
// a default highway example so that the web UI and the /api/example endpoint
// work regardless of the process working directory.
package example

import _ "embed"

// FreewayJSON is the embedded highway example used by the web UI and the
// /api/example endpoint. It mirrors example/freeway.json at the repository
// root so that the example is always available, even inside a container where
// the working directory may differ.
//
//go:embed freeway.json
var FreewayJSON []byte
