# Grupo14-Laboratorio-3

## Equivalencia tecnica

Broker = Servidor Central

Fulcrum = Data Node

Informante = Mismo rol que los servidores regionales antes

Vanguardia = Mismo rol que el server de la ONU en el lab 2

## Comandos

* Para compilar protobuf (ejecutar en carpeta proto): `protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative <nombre_archivo.proto>`

## Máquinas Virtuales

Máquina - Contraseña

- VM1: dist053 - Svwg5wPVPZT4

- VM2: dist054 - Zxq4deXdBnXy

- VM3: dist055 - DJekyBFztABd

- VM4: dist056 - TdRTwg8Cp775

## Distribución de Servidores en VMs

* VM1: Broker
* VM2: Informante 1, Fulcrum 1
* VM3: Informante 2, Fulcrum 2
* VM4: Vanguardia, Fulcrum 3
