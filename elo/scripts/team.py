from elods import EloDBQuery, EloDB, SortBy, EloGenerator
from util import getFile, writeFile
import math

eloGen = EloGenerator()

teamDB = EloDB()
regionDB = EloDB()

INIT_K = 8
MIN_K = 4
DECAY_RATE = math.e
DEFAULT_ELO = 1500

MAIN_PATH = '../../data/'
TOURNAMENT_PATH = MAIN_PATH + 'tournaments.json'
TEAM_PATH = MAIN_PATH + 'teams.json'
GAME_PATH = MAIN_PATH + 'games.json'
REGION_PATH = MAIN_PATH + 'team_regions.json'

REGION_K = 4

games = getFile(GAME_PATH)
games.sort(key=lambda x: x['end_time'])
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
    minELO = min(regionDB.getAll().values(), key=lambda x: x['elo'])['elo']
    blueRegionBuff = blueRegionELO - minELO
    redRegionBuff = redRegionELO - minELO

    currentRedElo = (teamDB.getCurrent(redTeam) if redTeam in teamDB.eloDB else DEFAULT_ELO)
    currentBlueElo = (teamDB.getCurrent(blueTeam) if blueTeam in teamDB.eloDB else DEFAULT_ELO)
    # NOTE: only the opponent's region buff is added to the current elo
    # this is because if we update both, the elo will be inflated since we will be adding every single time
    newBlueELO = eloGen.generate(currentBlueElo, currentRedElo + redRegionBuff, winner is 'blue', kFactor_BLUE)
    newRedELO = eloGen.generate(currentRedElo, currentBlueElo + blueRegionBuff, winner is 'red', kFactor_RED)

    teamDB.insert(redTeam, game['end_time'], newRedELO)
    teamDB.insert(blueTeam, game['end_time'], newBlueELO)

    newBlueRegionELO = eloGen.generate(blueRegionELO, redRegionELO, winner is 'blue', REGION_K)
    newRedRegionELO = eloGen.generate(redRegionELO, blueRegionELO, winner is 'red', REGION_K)
    regionDB.insert(redRegion, game['end_time'], newRedRegionELO)
    regionDB.insert(blueRegion, game['end_time'], newBlueRegionELO)

TEAM_ELO_PATH = MAIN_PATH + 'elo.json'
REGION_ELO_PATH = MAIN_PATH + 'region_elo.json'
teamDB.flush(TEAM_ELO_PATH)
regionDB.flush(REGION_ELO_PATH)

