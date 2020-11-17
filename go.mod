module github.com/vdbulcke/waypoint-plugin-hello-world

go 1.14

require (
	github.com/cloudflare/cfssl v1.5.0
	github.com/golang/protobuf v1.4.3
	github.com/hashicorp/waypoint-plugin-sdk v0.0.0-20201021094150-1b1044b1478e
	github.com/mitchellh/go-glint v0.0.0-20201015034436-f80573c636de
	google.golang.org/protobuf v1.25.0
)

// replace github.com/hashicorp/waypoint-plugin-sdk => ../../waypoint-plugin-sdk
