from elods import EloDBQuery, EloDB, SortBy, EloGenerator
from util import getFile, writeFile, regularize, sigmaRegularization, winProb
from trueskill import Rating, quality_1vs1, rate_1vs1, TrueSkill
import math

eloGen = EloGenerator()

teamDB = EloDB()
regionDB = EloDB()

INIT_K = 12
MIN_K = 4
DECAY_RATE = math.e
DEFAULT_ELO = 1500

MAIN_PATH = './data/'
TOURNAMENT_PATH = MAIN_PATH + 'tournaments.json'
TEAM_PATH = MAIN_PATH + 'teams.json'
GAME_PATH = MAIN_PATH + 'games_ts.json'
REGION_PATH = MAIN_PATH + 'team_region.json'

TEAM_ELO_PATH = MAIN_PATH + 'elo/elo.json'
REGION_ELO_PATH = MAIN_PATH + 'elo/region_elo.json'
SORTED_CURRENT_ELO_PATH = MAIN_PATH + 'elo/sorted_elo.json'
SORTED_REGION_ELO_PATH = MAIN_PATH + 'elo/sorted_region_elo.json'
PREDICTIONS_PATH = MAIN_PATH + 'elo/predictions.json'

REGION_K = 4

BUFF_MULTIPLER = 5

# TRUESKILL CONSTANTS'
INIT_MU = 25
INIT_SIGMA = INIT_MU / 3

INIT_REGION_MU = 5
INIT_REGION_SIGMA = INIT_REGION_MU / 2


games = getFile(GAME_PATH)
games = list(filter(lambda x: x.get('end_time') is not None, games))
games.sort(key=lambda x: x['end_time'])
teams = getFile(TEAM_PATH)
regions = getFile(REGION_PATH)

def getRegularizedELOs():
    sum = 0
    for teamID in teamDB.eloDB.keys():
        sum += teamDB.getCurrent(teamID)['elo']
    average = sum / (len(teamDB.eloDB.keys()) if len(teamDB.eloDB.keys()) > 0 else 1)

    teamMap = {}
    for team in teams:
        teamMap[team['team_id']] = team['name']

    eloMap = {}
    for teamID in teamDB.eloDB.keys():
        # numMatches = gameCount[teamID]

        currentELO = teamDB.getCurrent(teamID)
        currentELO['id'] = teamID
        if teamID not in teamMap:
            continue

        currentELO['name'] = teamMap[teamID]
        # currentELO['elo'] = regularize(currentELO['elo'], gameCount[teamID], average, 96, 1.05)
        currentELO['elo'] = sigmaRegularization(currentELO['elo'], currentELO['metadata']['sigma'], average, 0.245)
        eloMap[teamID] = currentELO

    return eloMap

gameCount = {}
bucket = [None] * 100
for game in games:
    if game['end_time'] < '2022-01-01':
        continue

    redTeam = game['red_id']
    blueTeam = game['blue_id']
    winner = game['winner']

    redRegion = regions[redTeam] if redTeam in regions else 'Unknown'
    blueRegion = regions[blueTeam] if blueTeam in regions else 'Unknown'
    redRegionMU = regionDB.getCurrent(redRegion)['elo'] if redRegion in regionDB.eloDB else INIT_REGION_MU
    redRegionSIGMA = regionDB.getCurrent(redRegion)['metadata']['sigma'] if redRegion in regionDB.eloDB else INIT_REGION_SIGMA
    blueRegionMU = regionDB.getCurrent(blueRegion)['elo'] if blueRegion in regionDB.eloDB else INIT_REGION_MU
    blueRegionSIGMA = regionDB.getCurrent(blueRegion)['metadata']['sigma'] if blueRegion in regionDB.eloDB else INIT_REGION_SIGMA

    #find minimum elo among all regions
    minELO = INIT_MU
    for region, events in regionDB.getAll().items():
        currentELO = regionDB.getCurrent(region)
        if currentELO['elo'] < minELO:
            minELO = currentELO['elo']
    blueRegionBuff = (blueRegionMU - minELO - blueRegionSIGMA)
    redRegionBuff = (redRegionMU - minELO - redRegionSIGMA)

    # only buff after some data is revealed
    if game['end_time'] >= '2022-01-01':
        blueRegionBuff *= BUFF_MULTIPLER
        redRegionBuff *= BUFF_MULTIPLER

    currentRedMU = (teamDB.getCurrent(redTeam)['elo'] if redTeam in teamDB.eloDB else INIT_MU)
    currentBlueMU = (teamDB.getCurrent(blueTeam)['elo'] if blueTeam in teamDB.eloDB else INIT_MU)
    currentRedSIGMA = (teamDB.getCurrent(redTeam)['metadata']['sigma'] if redTeam in teamDB.eloDB else INIT_SIGMA)
    currentBlueSIGMA = (teamDB.getCurrent(blueTeam)['metadata']['sigma'] if blueTeam in teamDB.eloDB else INIT_SIGMA)
    buffedRedELO = Rating(currentRedMU + redRegionBuff, currentRedSIGMA)
    buffedBlueELO = Rating(currentBlueMU + blueRegionBuff, currentBlueSIGMA)

    # for prediction purposes
    regularELOS = getRegularizedELOs()

    if redTeam in regularELOS and blueTeam in regularELOS and game['end_time'] >= '2022-01-01':
        redWinPercent = math.floor(winProb(
                Rating(regularELOS[redTeam]['elo'], regularELOS[redTeam]['metadata']['sigma']), 
                Rating(regularELOS[blueTeam]['elo'], regularELOS[blueTeam]['metadata']['sigma'])
        ) * 100)
        
        if bucket[redWinPercent] is None:
            bucket[redWinPercent] = {
                'correct': 0,
                'total': 0
            }

        bucket[redWinPercent]['total'] += 1
        if winner == 'red':
            bucket[redWinPercent]['correct'] = bucket[redWinPercent]['correct'] + 1

    # new ELO based on results
    newRedELO, newBlueELO = None, None
    if winner == 'red':
        newRedELO, _ = rate_1vs1(Rating(currentRedMU, currentRedSIGMA), buffedBlueELO)
        _, newBlueELO = rate_1vs1(buffedRedELO, Rating(currentBlueMU, currentBlueSIGMA))
    else:
        newBlueELO, _ = rate_1vs1(Rating(currentBlueMU, currentBlueSIGMA), buffedRedELO)
        _, newRedELO = rate_1vs1(buffedBlueELO, Rating(currentRedMU, currentRedSIGMA))
    
    teamDB.insert(redTeam, game['end_time'], newRedELO.mu, metadata={'mu': newRedELO.mu, 'sigma': newRedELO.sigma})
    teamDB.insert(blueTeam, game['end_time'], newBlueELO.mu, metadata={'mu': newBlueELO.mu, 'sigma': newBlueELO.sigma})


    if redRegion != blueRegion:
        newBlueRegionELO, newRedRegionELO = None, None
        if winner == 'blue':
            newBlueRegionELO, newRedRegionELO = rate_1vs1(Rating(blueRegionMU, blueRegionSIGMA), Rating(redRegionMU, redRegionSIGMA))
        else:
            newRedRegionELO, newBlueRegionELO  = rate_1vs1(Rating(redRegionMU, redRegionSIGMA), Rating(blueRegionMU, blueRegionSIGMA))

        regionDB.insert(redRegion, game['end_time'], newRedRegionELO.mu, metadata={'mu': newRedRegionELO.mu, 'sigma': newRedRegionELO.sigma})
        regionDB.insert(blueRegion, game['end_time'], newBlueRegionELO.mu, metadata={'mu': newBlueRegionELO.mu, 'sigma': newBlueRegionELO.sigma})
    
    gameCount[redTeam] = gameCount[redTeam] + 1 if redTeam in gameCount else 1
    gameCount[blueTeam] = gameCount[blueTeam] + 1 if blueTeam in gameCount else 1 

# write predictions to disk
writeFile(PREDICTIONS_PATH, bucket)

# write sorted elos for each team
sortedELOS = []
for teamID in teamDB.eloDB.keys():
    currentELO = teamDB.getCurrent(teamID)
    currentELO['id'] = teamID
    sortedELOS.append(currentELO)
sortedELOS.sort(key=lambda x: x['elo'], reverse=True)
writeFile(SORTED_CURRENT_ELO_PATH, sortedELOS)

# write sorted elos for each region
sortedELOS = []
for regionName in regionDB.eloDB.keys():
    currentELO = regionDB.getCurrent(regionName)
    currentELO['id'] = regionName
    sortedELOS.append(currentELO)
sortedELOS.sort(key=lambda x: x['elo'], reverse=True)
writeFile(SORTED_REGION_ELO_PATH, sortedELOS)

# flush databases to disk
TEAM_ELO_PATH = MAIN_PATH + 'elo/elo.json'
REGION_ELO_PATH = MAIN_PATH + 'elo/region_elo.json'
teamDB.flush(TEAM_ELO_PATH)
regionDB.flush(REGION_ELO_PATH)
