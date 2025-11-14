.PHONY: all give_permissions clean run migrate key build

all: clean build run

run:
	./bin/bastion.so

build:
	./scripts/build.sh

migrate:
	./scripts/migrate_down.sh
	./scripts/migrate_up.sh

key:
	./scripts/generate_encryption_key.sh

give_permissions:
	chmod +x ./scripts/build.sh
	chmod +x ./scripts/clean.sh
	chmod +x ./scripts/migrate_up.sh
	chmod +x ./scripts/migrate_down.sh
	chmod +x ./scripts/generate_encryption_key.sh

clean:
	./scripts/clean.sh
