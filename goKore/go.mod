module github.com/lenaxia/goKore

go 1.21

// Replace directives to map import paths to actual directory structure
replace github.com/lenaxia/goKore/network/send => ./network/send

replace github.com/lenaxia/goKore/network/hooks => ./network/hooks

replace github.com/lenaxia/goKore/network/common => ./network/common

replace github.com/lenaxia/goKore/network/receive => ./network/receive

replace github.com/lenaxia/goKore/network/receive/factory => ./network/receive/factory

replace github.com/lenaxia/goKore/network/receive/handlers/login => ./network/receive/handlers/login

replace github.com/lenaxia/goKore/network/receive/base => ./network/receive/base

replace github.com/lenaxia/goKore/network/receive/types => ./network/receive/types

replace github.com/lenaxia/goKore/network/send/factory => ./network/send/factory

replace github.com/lenaxia/goKore/network/send/handlers/game => ./network/send/handlers/game

replace github.com/lenaxia/goKore/network/send/handlers/login => ./network/send/handlers/login

replace github.com/lenaxia/goKore/network/send/handlers/servers => ./network/send/handlers/servers

// Replace directives for mikekao imports (after files were moved)
replace github.com/mikekao/openkore/goKore/network/implementation/network/connection => ./network/connection

replace github.com/mikekao/openkore/goKore/network/implementation/network/hooks => ./network/hooks

replace github.com/mikekao/openkore/goKore/network/implementation/network/protocol => ./network/protocol

replace github.com/mikekao/openkore/goKore/network/implementation/network/receive/core => ./network/receive/core

replace github.com/mikekao/openkore/goKore/network/implementation/network/receive/game/actor => ./network/receive/game/actor

replace github.com/mikekao/openkore/goKore/network/implementation/network/receive/security => ./network/receive/security

replace github.com/mikekao/openkore/goKore/network/implementation/network/send/core => ./network/send/core

replace github.com/mikekao/openkore/goKore/network/implementation/network/send/game/actor => ./network/send/game/actor

replace github.com/mikekao/openkore/goKore/network/implementation/network/config => ./network/config

replace github.com/mikekao/openkore/goKore/network/implementation/network/servers => ./network/servers

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
