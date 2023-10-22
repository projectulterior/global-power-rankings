import os 
import json

# script to get the region of each team

tournamentPath = './data/tournaments.json'
leaguePath = './data/leagues.json'

fd = os.open(tournamentPath, os.O_RDONLY)
tournaments = os.read(fd, os.path.getsize(tournamentPath))
tournaments = json.loads(tournaments)
os.close(fd)

fd = os.open(leaguePath, os.O_RDONLY)
leagues = os.read(fd, os.path.getsize(leaguePath))
leagues = json.loads(leagues)
os.close(fd)

# map tournament id to region
tournamentRegion = {}
for league in leagues:
    for tournament in league['tournaments']:
        tournamentRegion[tournament['id']] = league['region']

# count number of games a team plays in a region
teamRegionCount = {}
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
                    blueTeam = game['teams'][0]['id']
                    redTeam = game['teams'][1]['id']

                    if blueTeam not in teamRegionCount:
                        teamRegionCount[blueTeam] = {}
                    if redTeam not in teamRegionCount:
                        teamRegionCount[redTeam] = {}

                    blueTeamRegion = 'Unknown' if t['id'] not in tournamentRegion else tournamentRegion[t['id']]
                    redTeamRegion = 'Unknown' if t['id'] not in tournamentRegion else tournamentRegion[t['id']]

                    if blueTeamRegion not in teamRegionCount[blueTeam]:
                        teamRegionCount[blueTeam][blueTeamRegion] = 0
                    if redTeamRegion not in teamRegionCount[redTeam]:
                        teamRegionCount[redTeam][redTeamRegion] = 0

                    teamRegionCount[blueTeam][blueTeamRegion] += 1
                    teamRegionCount[redTeam][redTeamRegion] += 1

teamRegion = {}
for team, regions in teamRegionCount.items():
    teamRegion[team] = max(regions, key=regions.get)

data_bytes = json.dumps(teamRegion).encode('utf-8')
fd = os.open('team_region.json', os.O_WRONLY | os.O_CREAT)
os.write(fd, data_bytes)
os.close(fd)

print(teamRegion)