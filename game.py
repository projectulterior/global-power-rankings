import os
import json

tournamentPath = './data/tournaments.json'

fd = os.open(tournamentPath, os.O_RDONLY)
tournaments = os.read(fd, os.path.getsize(tournamentPath))
tournaments = json.loads(tournaments)
os.close(fd)

gameData = {}
for t in tournaments:
    stages = t['stages']
    for stage in stages:
        sections = stage['sections']
        for section in sections:
            matches = section['matches']
            for match in matches:
                games = match['games']
                for game in games:
                    if game['state'] != 'completed':
                        continue

                    gameID = game['id']
                    gameData[gameID] = {}

                    gameData[gameID]['blue'] = game['teams'][0]['id']
                    gameData[gameID]['red'] = game['teams'][1]['id']

                    gameData[gameID]['winner'] = 'blue' if game['teams'][0]['result']['outcome'] == 'win' else 'red'

print(len(gameData))

data_bytes = json.dumps(gameData).encode('utf-8')
fd = os.open('flat_games.json', os.O_WRONLY | os.O_CREAT)
os.write(fd, data_bytes)
os.close(fd)
