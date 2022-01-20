module github.com/cserrant/image-transform-cli

go 1.15

replace github.com/Chadius/creating-symmetry v0.0.0-20220116205726-9731f64f2462 => ../creating-symmetry

require (
	github.com/Chadius/creating-symmetry v0.0.0-20220116205726-9731f64f2462
	github.com/chadius/image-transform-server v0.0.0-20220119003333-85a12379ad24
	github.com/maxbrunsfeld/counterfeiter/v6 v6.4.1
	github.com/stretchr/testify v1.7.0
	github.com/twitchtv/twirp v8.1.1+incompatible
	google.golang.org/protobuf v1.23.0
)
