# Grupo14-Laboratorio-3

## Equivalencia tecnica

Broker = Servidor Central

Fulcrum = Data Node

Informante = Mismo rol que los servidores regionales antes

Vanguardia = Mismo rol que el server de la ONU en el lab 2

## Comandos

* Para compilar protobuf (ejecutar en carpeta proto): protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative <nombre_archivo.proto>