module github.com/kiwicom/easycql

go 1.13

require gopkg.in/inf.v0 v0.9.1

require (
	github.com/gocql/gocql v0.0.0-20190922122429-7b17705d7514
	github.com/mailru/easyjson v0.7.0
	github.com/stretchr/testify v1.3.0
)

//replace github.com/gocql/gocql => github.com/kiwicom/gocql v0.0.0-20190920161558-af217b7dd1b8
replace github.com/gocql/gocql => /home/martin/Projects/gocql
