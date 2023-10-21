import os
import json
import math


tournamentPath = './flat_games.json'

fd = os.open(tournamentPath, os.O_RDONLY)
tournaments = os.read(fd, os.path.getsize(tournamentPath))
tournaments = json.loads(tournaments)
os.close(fd)

INIT_K = 4
MIN_K = 4
DEFAULT_ELO = 1500

eloDB = {}
gameCount = {}
for gameID, gameResult in tournaments.items():
    redTeam = gameResult['red']
    blueTeam = gameResult['blue']
    winner = gameResult['winner']

    kFactor_RED = (INIT_K - MIN_K) * math.e ** (0 if redTeam not in gameCount else -gameCount[redTeam] / 100) + MIN_K
    kFactor_BLUE = (INIT_K - MIN_K) * math.e ** (0 if blueTeam not in gameCount else -gameCount[blueTeam] / 100) + MIN_K

    currentRedElo = eloDB[redTeam] if redTeam in eloDB else DEFAULT_ELO
    currentBlueElo = eloDB[blueTeam] if blueTeam in eloDB else DEFAULT_ELO

    expectedRed = 1 / (1 + 10 ** ((currentRedElo - currentBlueElo) / 400))
    expectedBlue = 1 / (1 + 10 ** ((currentBlueElo - currentRedElo) / 400))

    if winner == 'red':
        newRedElo = currentRedElo + kFactor_RED * (1 - expectedRed)
        newBlueElo = currentBlueElo + kFactor_BLUE * (0 - expectedBlue)
    elif winner == 'blue':
        newRedElo = currentRedElo + kFactor_RED * (0 - expectedRed)
        newBlueElo = currentBlueElo + kFactor_BLUE * (1 - expectedBlue)
    
    eloDB[redTeam] = newRedElo
    eloDB[blueTeam] = newBlueElo

    gameCount[redTeam] = gameCount[redTeam] + 1 if redTeam in gameCount else 1
    gameCount[blueTeam] = gameCount[blueTeam] + 1 if blueTeam in gameCount else 1

teamPath = './data/teams.json'
fd = os.open(teamPath, os.O_RDONLY)
teams = os.read(fd, os.path.getsize(tournamentPath))
teams = json.loads(teams)
os.close(fd)

teamMap = {}
for team in teams:
    teamMap[team['team_id']] = team['name']

eloArr = []
for team, elo in eloDB.items():
    eloArr.append({'id': team, 'elo': elo, 'name': None if team not in teamMap else teamMap[team]})
    eloArr.sort(key=lambda x: x['elo'], reverse=True)

data_bytes = json.dumps(eloDB).encode('utf-8')
fd = os.open('elo.json', os.O_WRONLY | os.O_CREAT)
os.write(fd, data_bytes)
os.close(fd)

data_bytes = json.dumps(eloArr).encode('utf-8')
fd = os.open('eloList.json', os.O_WRONLY | os.O_CREAT)
os.write(fd, data_bytes)
os.close(fd)




