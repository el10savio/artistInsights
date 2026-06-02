.PHONY: docker-up docker-down docker-restart docker-wait \
        migrate-up migrate-down migrate-status migrate-create \
        server-build server-up server-down \
        provision \
        load-test load-full load-artists load-users truncate-listens \
        test test-integration \
        frontend-install frontend-dev frontend-build

COMPOSE_PROJECT := artistinsights
NETWORK         := $(COMPOSE_PROJECT)_default
CH_DSN          := "clickhouse://clickhouse:9000?database=artist_insights&username=admin&password=pass"
MIGRATE         := docker run --rm \
                     -v $(PWD)/migrations:/migrations \
                     --network $(NETWORK) \
                     migrate/migrate

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-clean:
	docker-compose down -v

docker-restart:
	docker-compose restart

docker-wait:
	@echo "Waiting for ClickHouse..."
	@until docker exec artistinsights_clickhouse wget -q --spider http://localhost:8123/ping 2>/dev/null; \
	  do sleep 1; done
	@docker exec artistinsights_clickhouse clickhouse-client \
	  --query "CREATE USER IF NOT EXISTS admin IDENTIFIED WITH plaintext_password BY 'pass'; \
	           GRANT ALL ON artist_insights.* TO admin;"
	@echo "ClickHouse is ready."

migrate-up: docker-up docker-wait
	$(MIGRATE) -path /migrations -database $(CH_DSN) up

migrate-down: docker-wait
	$(MIGRATE) -path /migrations -database $(CH_DSN) down 1

migrate-status: docker-wait
	$(MIGRATE) -path /migrations -database $(CH_DSN) version

migrate-create:
	@test -n "$(name)" || (echo "Usage: make migrate-create name=<migration_name>" && exit 1)
	docker run --rm -v $(PWD)/migrations:/migrations migrate/migrate \
	  create -ext sql -dir /migrations -seq $(name)

server-build:
	docker build -t artistinsights_server .

server-up: docker-up docker-wait
	docker-compose up -d server

server-down:
	docker-compose stop server

provision: migrate-up server-up
	@echo "Provisioning complete."

TSV         := dataset/lastfm-dataset-1K/userid-timestamp-artid-artname-traid-traname.tsv
PROFILE_TSV := dataset/lastfm-dataset-1K/userid-profile.tsv
AWK_COLS    := awk -F'\t' 'BEGIN{OFS="\t"} $$3!="" && $$5!="" && $$6!="" {gsub(/\\/, "\\\\", $$6); gsub(/\\/, "\\\\", $$4); print $$3, $$5, $$6, $$1, $$2, $$4}'
AWK_ARTISTS := awk -F'\t' 'BEGIN{OFS="\t"} NR>1 && $$3!="" && $$4!="" {gsub(/\\/, "\\\\", $$4); print $$3, $$4}'
AWK_USERS   := awk -F'\t' 'BEGIN{OFS="\t"} NR>1 && $$1!="" && $$4!="" {print $$1, $$4}'
CH_INSERT := docker exec -i artistinsights_clickhouse clickhouse-client \
               --user admin --password pass \
               --database artist_insights \
               --date_time_input_format=best_effort \
               --query "INSERT INTO listens(artist_id, track_id, track_name, user_id, ts, artist_name) FORMAT TabSeparated"
CH_INSERT_ARTISTS := docker exec -i artistinsights_clickhouse clickhouse-client \
               --user admin --password pass \
               --database artist_insights \
               --query "INSERT INTO artists(artist_id, artist_name) FORMAT TabSeparated"
CH_INSERT_USERS := docker exec -i artistinsights_clickhouse clickhouse-client \
               --user admin --password pass \
               --database artist_insights \
               --query "INSERT INTO users(user_id, country) FORMAT TabSeparated"

load-artists: docker-wait
	$(AWK_ARTISTS) $(TSV) | sort -u | $(CH_INSERT_ARTISTS)
	@echo "Loaded artists."

load-users: docker-wait
	$(AWK_USERS) $(PROFILE_TSV) | sort -u | $(CH_INSERT_USERS)
	@echo "Loaded users."

load-test: docker-wait
	head -5000 $(TSV) | $(AWK_COLS) | $(CH_INSERT)
	$(AWK_USERS) $(PROFILE_TSV) | sort -u | $(CH_INSERT_USERS)
	@echo "Loaded 5k test rows."

load-full: docker-wait
	$(AWK_COLS) $(TSV) | $(CH_INSERT)
	$(AWK_USERS) $(PROFILE_TSV) | sort -u | $(CH_INSERT_USERS)
	@echo "Loaded full dataset."

test:
	go test ./src/...

test-integration: docker-wait
	go test -tags integration ./src/infra/...

truncate-listens: docker-wait
	docker exec artistinsights_clickhouse clickhouse-client \
	  --user admin --password pass \
	  --query "TRUNCATE TABLE artist_insights.listens"
	@echo "listens table cleared."

frontend-install:
	cd frontend && npm install

frontend-dev: frontend-install
	cd frontend && npm run dev

frontend-build: frontend-install
	cd frontend && npm run build
