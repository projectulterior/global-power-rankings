import json
from util import getFile, writeFile, regularize, sigmaRegularization, winProb
import math

GAME_PATH = './data/analysis/games_kda.json'
SORTED_CURRENT_ELO_PATH = './data/elo/sorted_elo.json'
INITIAL_ELO_PATH = './data/initial_elo.json'

WEIGHTS = {
    'worlds': 50,
    'msi': 30,
    'easternRegionGamesTier1': 15,
    'westernRegionGamesTier1': 10, 
    'easternRegionWinsTier1': 15, 
    'westernRegionWinsTier1': 10, 
    'easternRegionGamesTier2': 3, 
    'westernRegionGamesTier2': 2,
    'easternRegionWinsTier2': 3,
    'westernRegionWinsTier2': 2,
    'topKDA': 0,
    'jungleKDA': 0,
    'midKDA': 0,
    'adcKDA': 0,
    'supportKDA': 0,
    'teamGold': 0,
}

worlds = set(["108998961191900167", "106926282333089592", "104841804583318464"])
msi = set(["110198981276611770", "107905284983470149", "105873410870441926"])
westernRegionTier1 = set(["LCK", "LPL"]) #4
westernRegionTier2 = set(["PCS", "VCS", "LJL"]) #1
easternRegionTier1 = set(["LCS", "LEC"]) #3
easternRegionTier2 = set(["CBLOL", "LLA"]) #1
semileague = ["CHALLENGER", "MASTERS", "ACADEMY"]
positions = ["top", "jungle", "mid", "adc", "support"]


def getAppearances():
    toRet = {}

    games = getFile(GAME_PATH)
    for game in games:
        redTeamID = game['red_id']
        blueTeamID = game['blue_id']
        winner = game['winner']

        toAddRed = {
            'worlds': 0,
            'msi': 0,
            'easternRegionGamesTier1': 0,
            'westernRegionGamesTier1': 0,
            'easternRegionWinsTier1': 0,
            'westernRegionWinsTier1': 0,
            'easternRegionGamesTier2': 0,
            'westernRegionGamesTier2': 0,
            'easternRegionWinsTier2': 0,
            'westernRegionWinsTier2': 0,
            'topKDA': 0,
            'jungleKDA': 0,
            'midKDA': 0,
            'adcKDA': 0,
            'supportKDA': 0,
            'teamGold': 0,
            'updates': 0,

        } if redTeamID not in toRet else toRet[redTeamID]

        toAddBlue = {
            'worlds': 0,
            'msi': 0,
            'easternRegionGamesTier1': 0,
            'westernRegionGamesTier1': 0,
            'easternRegionWinsTier1': 0,
            'westernRegionWinsTier1': 0,
            'easternRegionGamesTier2': 0,
            'westernRegionGamesTier2': 0,
            'easternRegionWinsTier2': 0,
            'westernRegionWinsTier2': 0,
            'topKDA': 0,
            'jungleKDA': 0,
            'midKDA': 0,
            'adcKDA': 0,
            'supportKDA': 0,
            'teamGold': 0,
            'updates': 0
            
        } if blueTeamID not in toRet else toRet[blueTeamID]

        # tournament appearances
        tournamentID = game['tournament_id']
        if tournamentID in worlds:
            toAddRed['worlds'] += 1
            toAddBlue['worlds'] += 1
        elif tournamentID in msi:
            toAddRed['msi'] += 1
            toAddBlue['msi'] += 1

        # avoid challenger and academy leagues
        tournamentSlug = str(game['tournament_slug'])
        league = tournamentSlug.split('_')[0].upper()
        flag = False
        for avoid in semileague:
            if avoid in league.upper():
                flag = True
        if flag:
            continue

        # games and wins
        def updateRegionGameStats(variable: str, winner: str):
            toAddRed[variable] += 1
            toAddBlue[variable] += 1
            if winner == 'red':
                toAddRed[variable] += 1
            else:
                toAddBlue[variable] += 1

        if league in westernRegionTier1:
            updateRegionGameStats('westernRegionGamesTier1', winner)
        if league in westernRegionTier2:
            updateRegionGameStats('westernRegionGamesTier2', winner)
        if league in easternRegionTier1:
            updateRegionGameStats('easternRegionGamesTier1', winner)
        if league in easternRegionTier2:
            updateRegionGameStats('easternRegionGamesTier2', winner)

        # kdas
        # TODO: average it out instead
        toAddRed['updates'] += 1
        toAddBlue['updates'] += 1

        def updateKDAStats(position: str, team: str, kda: float):
            if team == 'red':
                toAddRed[position + 'KDA'] += (kda - toAddRed[position + 'KDA']) / toAddRed['updates']
            else:
                toAddBlue[position + 'KDA'] += (kda - toAddBlue[position + 'KDA']) / toAddBlue['updates']
        
        # TODO: average it out instead
        def updateTeamGoldStats(team: str, gold: int):
            if team == 'red':
                toAddRed['teamGold'] += (gold - toAddRed['teamGold']) / toAddRed['updates']
            else:
                toAddBlue['teamGold'] += (gold - toAddBlue['teamGold']) / toAddBlue['updates']


        def getTeamKDA(team: str):
            kills = 0
            deaths = 0
            for pos in positions:
                currentKills = game[team + '_' + pos + '_kills'] if (team + '_' + pos + '_kills') in game is not None else 0
                currentAssists = game[team + '_' + pos + '_assists'] if (team + '_' + pos + '_assists') in game is not None else 0
                currentDeaths = game[team + '_' + pos + '_deaths'] if (team + '_' + pos + '_deaths') in game is not None else 0
                kills += currentKills + currentAssists
                deaths += currentDeaths 
            return kills / deaths if deaths != 0 else kills


        for pos in positions:
            redKDA = getTeamKDA('red')
            blueKDA = getTeamKDA('blue')

            # if redKDA != 0 and blueKDA != 0:
            #     print(redKDA, blueKDA)
            updateKDAStats(pos, 'red', redKDA)
            updateKDAStats(pos, 'blue', blueKDA)
            updateTeamGoldStats('red', game['red_gold'] if 'red_gold' in game else 0)
            updateTeamGoldStats('blue', game['blue_gold'] if 'blue_gold' in game else 0)

        toRet[redTeamID] = toAddRed
        toRet[blueTeamID] = toAddBlue
    return toRet

def getInitialScore(stats): 
    score = 0
    for key, weight in WEIGHTS.items():
        score += weight * stats[key]
    return score

def createInitialScores():
    toRet = {}

    appearances = getAppearances()

    # sortedElo = getFile(SORTED_CURRENT_ELO_PATH)
    maxElo = 170 #sortedElo[0]['elo']
    minElo = 10 #sortedElo[-1]['elo']

    # normalize
    minScore = math.inf
    maxScore = -math.inf
    for _, stats in appearances.items():
        score = getInitialScore(stats)
        if score < minScore:
            minScore = score
        if score > maxScore:
            maxScore = score
    
    diff = maxScore - minScore
    for teamID, stats in appearances.items():
        score = getInitialScore(stats)
        score = (score - minScore) / diff if diff != 0 else 1
        toRet[teamID] = score
    
    for teamID in toRet:
        toRet[teamID] = minElo + (maxElo - minElo) * toRet[teamID]
    
    return toRet

scoreMap = createInitialScores()
writeFile(INITIAL_ELO_PATH, scoreMap)
        

        

        


