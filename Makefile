init:
	./scripts/init.sh

clean:
	rm -rf data

parse:
	python3 parse.py

games:
	mkdir -p data/games
	go run games.go

.PHONY: parser elo

parser:
	cd parser && go run ./...

elo:
	rm -rf data/elo
	mkdir -p data/elo
	python3 ./elo/scripts/team.py