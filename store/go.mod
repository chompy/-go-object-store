module gitlab.com/contextualcode/go-object-store/store

go 1.17

replace gitlab.com/contextualcode/go-object-store/types => ../types

require (
	github.com/caibirdme/yql v0.0.0-20210122071211-a800d6de28a0
	github.com/matoous/go-nanoid/v2 v2.0.0
	github.com/philippgille/gokv v0.6.0
	github.com/philippgille/gokv/file v0.6.0
	github.com/philippgille/gokv/redis v0.6.0
	github.com/philippgille/gokv/syncmap v0.6.0
	github.com/pkg/errors v0.9.1
	gitlab.com/contextualcode/go-object-store/types v0.0.1
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	github.com/antlr/antlr4 v0.0.0-20210121092344-5dce78c87a9e // indirect
	github.com/go-redis/redis v6.15.6+incompatible // indirect
	github.com/philippgille/gokv/encoding v0.0.0-20191011213304-eb77f15b9c61 // indirect
	github.com/philippgille/gokv/util v0.0.0-20191011213304-eb77f15b9c61 // indirect
)
