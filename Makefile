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
	gzip data/analysis/games_kda.json

unzip:
	gzip -d data/analysis/games_kda.json.gz

initElo:
	if [ -f data/elo/inital_elo.json ]; then rm data/elo/inital_elo.json; fi
	python3 ./elo/scripts/appearances.py
