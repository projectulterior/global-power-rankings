from elods import EloDBQuery, EloDB, SortBy, EloGenerator
from util import getFile, writeFile
import math

eloGen = EloGenerator()

teamDB = EloDB()
regionDB = EloDB()

INIT_K = 4
MIN_K = 4
DECAY_RATE = math.e
DEFAULT_ELO = 1500

MAIN_PATH = './data/'
TOURNAMENT_PATH = MAIN_PATH + 'tournaments.json'
TEAM_PATH = MAIN_PATH + 'teams.json'
GAME_PATH = MAIN_PATH + 'games.json'
REGION_PATH = MAIN_PATH + 'team_region.json'

REGION_K = 4

games = getFile(GAME_PATH)
# TODO: change back to end_time
games.sort(key=lambda x: x['tournament_endDate'])
teams = getFile(TEAM_PATH)
regions = getFile(REGION_PATH)

gameCount = {}
for game in games:
    redTeam = game['red_id']
    blueTeam = game['blue_id']
    winner = game['winner']

    kFactor_RED = (INIT_K - MIN_K) * DECAY_RATE ** (0 if redTeam not in gameCount else -gameCount[redTeam] / 100) + MIN_K
    kFactor_BLUE = (INIT_K - MIN_K) * DECAY_RATE ** (0 if blueTeam not in gameCount else -gameCount[blueTeam] / 100) + MIN_K

    redRegion = regions[redTeam] if redTeam in regions else 'Unknown'
    blueRegion = regions[blueTeam] if blueTeam in regions else 'Unknown'
    redRegionELO = regionDB.getCurrent(redRegion)
    redRegionELO = redRegionELO if redRegionELO is not None else DEFAULT_ELO
    blueRegionELO = regionDB.getCurrent(blueRegion)
    blueRegionELO = blueRegionELO if blueRegionELO is not None else DEFAULT_ELO
    minELO = min(regionDB.getAll().values(), key=lambda x: x['elo'])['elo'] if len(regionDB.getAll()) > 0 else DEFAULT_ELO
    blueRegionBuff = blueRegionELO - minELO
    redRegionBuff = redRegionELO - minELO

    currentRedElo = (teamDB.getCurrent(redTeam)['elo'] if redTeam in teamDB.eloDB else DEFAULT_ELO)
    currentBlueElo = (teamDB.getCurrent(blueTeam)['elo'] if blueTeam in teamDB.eloDB else DEFAULT_ELO)
    # NOTE: only the opponent's region buff is added to the current elo
    # this is because if we update both, the elo will be inflated since we will be adding every single time
    newBlueELO = eloGen.generate(currentBlueElo, currentRedElo + redRegionBuff, winner == 'blue', kFactor_BLUE)
    newRedELO = eloGen.generate(currentRedElo, currentBlueElo + blueRegionBuff, winner == 'red', kFactor_RED)

    # TODO: change back to end_time
    teamDB.insert(redTeam, game['tournament_endDate'], newRedELO)
    teamDB.insert(blueTeam, game['tournament_endDate'], newBlueELO)

    # newBlueRegionELO = eloGen.generate(blueRegionELO, redRegionELO, winner is 'blue', REGION_K)
    # newRedRegionELO = eloGen.generate(redRegionELO, blueRegionELO, winner is 'red', REGION_K)
    # regionDB.insert(redRegion, game['end_time'], newRedRegionELO)
    # regionDB.insert(blueRegion, game['end_time'], newBlueRegionELO)

    gameCount[blueTeam] = gameCount[blueTeam] + 1 if blueTeam in gameCount else 1
    gameCount[redTeam] = gameCount[redTeam] + 1 if redTeam in gameCount else 1

TEAM_ELO_PATH = MAIN_PATH + 'elo/elo.json'
REGION_ELO_PATH = MAIN_PATH + 'elo/region_elo.json'
teamDB.flush(TEAM_ELO_PATH)
regionDB.flush(REGION_ELO_PATH)

