module "main"

go 1.16

require (
	app v1.0.0
)

replace (
	app => ./app
)