init:
	./scripts/init.sh

clean:
	rm -rf data

parse:
	python3 parse.py

games:
	mkdir -p data/games
	go run games.go

.PHONY: parser fast elo benchmark

parser:
	cd parser && go run ./...

elo:
	rm -rf data/elo
	mkdir -p data/elo
	python3 ./elo/scripts/team.py

benchmark:
	cd parser && go test -bench ./...

fast:
	cd fast && go run ./...

zip:
	gzip --keep data/analysis/games_kda.json
	gzip --keep data/analysis/games.json

unzip:
	gzip --keep -d data/analysis/games_kda.json.gz
	gzip --keep -d data/analysis/games.json.gz