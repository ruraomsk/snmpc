module github.com/ruraomsk/snmpc

go 1.23.4

replace github.com/ruraomsk/potop => ../potop

require (
	github.com/gosnmp/gosnmp v1.42.1
	github.com/ruraomsk/potop v0.0.0-00010101000000-000000000000
)

require (
	github.com/goburrow/serial v0.1.0 // indirect
	golang.org/x/text v0.13.0 // indirect
)
