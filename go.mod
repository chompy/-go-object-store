module gitlab.com/contextualcode/go-object-store

go 1.16

replace (
	gitlab.com/contextualcode/go-object-store/http => ./http
	gitlab.com/contextualcode/go-object-store/store => ./store
	gitlab.com/contextualcode/go-object-store/types => ./types
)

require (
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	gitlab.com/contextualcode/go-object-store/http v0.0.0-00010101000000-000000000000
	gitlab.com/contextualcode/go-object-store/store v0.0.0-00010101000000-000000000000
	gitlab.com/contextualcode/go-object-store/types v0.0.1
)
