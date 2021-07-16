module gitlab.com/contextualcode/go-object-store/http

go 1.17

replace gitlab.com/contextualcode/go-object-store/types => ../types

replace gitlab.com/contextualcode/go-object-store/store => ../store

require (
	github.com/pkg/errors v0.9.1
	gitlab.com/contextualcode/go-object-store/store v0.0.0-00010101000000-000000000000
	gitlab.com/contextualcode/go-object-store/types v0.0.1
)

require (
	github.com/antlr/antlr4 v0.0.0-20210121092344-5dce78c87a9e // indirect
	github.com/caibirdme/yql v0.0.0-20210122071211-a800d6de28a0 // indirect
	github.com/go-redis/redis v6.15.6+incompatible // indirect
	github.com/matoous/go-nanoid/v2 v2.0.0 // indirect
	github.com/philippgille/gokv v0.6.0 // indirect
	github.com/philippgille/gokv/encoding v0.0.0-20191011213304-eb77f15b9c61 // indirect
	github.com/philippgille/gokv/file v0.6.0 // indirect
	github.com/philippgille/gokv/redis v0.6.0 // indirect
	github.com/philippgille/gokv/syncmap v0.6.0 // indirect
	github.com/philippgille/gokv/util v0.0.0-20191011213304-eb77f15b9c61 // indirect
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
