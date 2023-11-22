.PHONY: docker-broker docker-vanguardia docker-f1 docker-f2 docker-f3 docker-i1 docker-i2 help

docker-broker: ## Initiates the Docker code for the broker server
	@docker build -f dockerfile.broker -t broker .
	@docker run -p 50051:50051 -d --name broker broker

docker-vanguardia: ## Initiates the Docker code for the vanguard server
	@docker build -f dockerfile.vanguardia -t vanguardia .
	@docker run -i -p 50050:50050 --name vanguardia vanguardia

docker-f1: ## Initiates the Docker code for fulcrum 1
	@docker build -f dockerfile.fulcrum --build-arg NUM=1 -t fulcrum1 .
	@docker run -p 50056:50056 -d --name fulcrum1 fulcrum1

docker-f2: ## Initiates the Docker code for fulcrum 2
	@docker build -f dockerfile.fulcrum --build-arg NUM=2 -t fulcrum2 .
	@docker run -p 50057:50057 -d --name fulcrum2 fulcrum2

docker-f3: ## Initiates the Docker code for fulcrum 3
	@docker build -f dockerfile.fulcrum --build-arg NUM=3 -t fulcrum3 .
	@docker run -p 50058:50058 -d --name fulcrum3 fulcrum3

docker-i1: ## Initiates the Docker code for Caiatl
	@docker build -f dockerfile.caiatl -t caiatl .
	@docker run -i -p 50050:50050 --name caiatl caiatl

docker-i2: ## Initiates the Docker code for Osiris
	@docker build -f dockerfile.osiris -t osiris .
	@docker run -i -p 50050:50050 --name osiris osiris

clean: ## Remove all Docker containers and images
	@docker rm -f $$(docker ps -a -q) || true
	@docker rmi -f $$(docker images -q) || true

help: ## Display this help message
	@echo "Usage:"
	@echo "  make <target>"
	@echo "Targets:"
	@echo "  docker-broker   Initiates the Docker code for the broker server"
	@echo "  docker-vanguardia Initiates the Docker code for the vanguard server"
	@echo "  docker-f1	   Initiates the Docker code for fulcrum 1"
	@echo "  docker-f2	   Initiates the Docker code for fulcrum 2"
	@echo "  docker-f3	   Initiates the Docker code for fulcrum 3"
	@echo "  docker-i1	   Initiates the Docker code for Caiatl"
	@echo "  docker-i2	   Initiates the Docker code for Osiris"
	@echo "  clean		   Remove all Docker containers and images"
