module github.com/robertgzr/toggl

go 1.13

require (
	github.com/dougEfresh/gtoggl-api v0.0.0-20190813143908-d4a77f6fd0a6
	github.com/pkg/errors v0.8.1
	github.com/urfave/cli v1.20.0
)

replace github.com/dougEfresh/gtoggl-api => github.com/robertgzr/gtoggl-api v0.1.1-0.20191206154203-6c421c2736b5 // token_always branch
