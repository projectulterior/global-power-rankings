#!/bin/bash

S3='https://power-rankings-dataset-gprhack.s3.us-west-2.amazonaws.com/'

mkdir -p data
cd data

function get {
    echo downloading $1...
    curl $S3$1 --output $2
    gzip --force --decompress $2
}

get esports-data/leagues.json.gz leagues.json.gz
get esports-data/tournaments.json.gz tournaments.json.gz
get esports-data/players.json.gz player.json.gz
get esports-data/teams.json.gz teams.json.gz