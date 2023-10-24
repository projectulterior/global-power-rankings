# [Global Power Rankings Hackathon](https://lolglobalpowerrankings.devpost.com/)

Develop a method that ranks the top global LoL Esports teams using official Riot Games data and AWS services.

The working method must provide the following outputs:

-   Tournament rankings: provide team rankings for a given tournament
-   Global rankings: provide current rankings of all teams globally
-   Team rankings: provide rankings for a given list of teams

[API Requirements](https://docs.google.com/document/d/1Klodp4YqE6bIOES026ecmNb_jS5IOntRqLv5EmDAXyc/edit)

## Data

S3: https://power-rankings-dataset-gprhack.s3.us-west-2.amazonaws.com/

-   Game Data: Events sent by the game server

```sh
games/ESPORTSTMNT01:1110148.json.gz
```

-   Esports fixture data: Information about the leagues, tournaments, players, team and schedule

```sh
esports-data/leagues.json.gz
esports-data/tournaments.json.gz
esports-data/players.json.gz
esports-data/teams.json.gz
```

-   Esports mapping data: Joins the livestats data with the esports fixture data

```sh
esports-data/mapping_data.json.gz
```

[Technical Doc](https://docs.google.com/document/d/1wFRehKMJkkRR5zyjEZyaVL9H3ZbhP7_wP0FBE5ID40c/edit#heading=h.4osafmixo0au)

[Athena Setup](https://docs.google.com/document/d/14uhbMUYb7cR_Hg6UWjlAgnN-hSy0ymhz19-_A6eidxI/edit#heading=h.mn6lxq2agqoh)

[S3 Bucket](https://s3.console.aws.amazon.com/s3/buckets/power-rankings-dataset-gprhack)

## Game Data

Game data has been downloaded to an EC2 which can be accessed thru SSH.

1. Download [`production.pem`](https://github.com/projectulterior/devops/blob/master/.keys/production.pem) and move it to your `~/.ssh` directory

2. Copy the following in to your `~/.ssh/config` file

```sh
Host lol-games
    HostName 52.13.219.113
    User ec2-user
    IdentityFile ~/.ssh/production.pem
```
