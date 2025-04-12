# Os obxectivos .PHONY son aqueles que non xeran ficheiros co mesmo nome.
# Sen esta declaración, 'make' pode omitir a execución dalgúns comandos se
# existe un ficheiro co mesmo nome do obxectivo ou se considera que o 
# obxectivo está actualizado ao non ter dependencias que cambiaran.
# Neste caso dado que existe un cartafol chamado test no noso proxecto
# o make confundiase e trataba de actualizar este ficheiro en lugares de 
# executar o comando test. Chegaría con ".PHONY: test" neste caso
# pero engado todos por se acaso.
.PHONY: update build run run-dev run-dev-windows test

DEV-ENVIRONMENT=development
PROD-ENVIRONMENT=production

update:
	go mod tidy

build: update
	go build -o ./deployment/${PROD-ENVIRONMENT}/aggregationFramework cmd/main.go

run:
	export ENVIRONMENT="${PROD-ENVIRONMENT}" && go run cmd/main.go

run-dev:
	export ENVIRONMENT="${DEV-ENVIRONMENT}" && go run ./cmd/main.go

run-dev-windows: 
	set ENVIRONMENT=${DEV-ENVIRONMENT} && go run ./cmd/main.go

test:
	go generate -v ./internal/... && go test ./internal/...