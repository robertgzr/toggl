module github.com/robertgzr/toggl

require (
	github.com/dougEfresh/gtoggl v0.0.0-20181020132158-0c52db5c669b
	github.com/dougEfresh/gtoggl-api v0.0.0-20181029170833-3dde07b8438e
	github.com/hashicorp/golang-lru v0.5.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/throttled/throttled v2.2.4+incompatible // indirect
	github.com/urfave/cli v1.20.0
)

replace github.com/dougEfresh/gtoggl-api => github.com/robertgzr/gtoggl-api latest
