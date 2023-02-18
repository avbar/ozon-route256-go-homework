
build-all:
	cd checkout && make build
	cd loms && make build
	cd notifications && make build

run-all: build-all
	sudo docker compose up --force-recreate --build

precommit:
	cd checkout && make precommit
	cd loms && make precommit
	cd notifications && make precommit
