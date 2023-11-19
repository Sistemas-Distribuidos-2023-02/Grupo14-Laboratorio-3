.PHONY: docker-broker docker-vanguardia docker-f1 docker-f2 docker-f3 docker-i1 docker-i2 help

docker-broker: ## Initiates the Docker code for the broker server
	@docker build -f dockerfile.broker -t broker .

docker-vanguardia: ## Initiates the Docker code for the vanguard server
	@docker build -f dockerfile.vanguardia -t vanguardia .

docker-f1: ## Initiates the Docker code for fulcrum 1
	@docker build -f dockerfile.fulcrum --build-arg NUM=1 -t fulcrum1 .

docker-f2: ## Initiates the Docker code for fulcrum 2
	@docker build -f dockerfile.fulcrum --build-arg NUM=2 -t fulcrum2 .

docker-f3: ## Initiates the Docker code for fulcrum 3
	@docker build -f dockerfile.fulcrum --build-arg NUM=3 -t fulcrum3 .

docker-i1: ## Initiates the Docker code for Caiatl
	@docker build -f dockerfile.caiatl -t caiatl .

docker-i2: ## Initiates the Docker code for Osiris
	@docker build -f dockerfile.osiris -t osiris .

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