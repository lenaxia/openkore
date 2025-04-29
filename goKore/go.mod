module github.com/lenaxia/goKore

go 1.21

// Replace directives to map import paths to actual directory structure
replace github.com/lenaxia/goKore/network/send => ./network/send
replace github.com/lenaxia/goKore/network/hooks => ./network/hooks
replace github.com/lenaxia/goKore/network/common => ./network/common
