init:
	./scripts/init.sh

clean:
	rm -rf data

parse:
	python3 parse.py

games:
	mkdir -p data/games
	go run games.go