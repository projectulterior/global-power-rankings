import json
from util import getFile, writeFile, regularize, sigmaRegularization, winProb
from team import GAME_PATH, SORTED_CURRENT_ELO_PATH, INITIAL_ELO_PATH
import math

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
    'topKDA': 1,
    'jungleKDA': 2,
    'midKDA': 3,
    'adcKDA': 4,
    'supportKDA': 5,
    'teamGold': 6,

}

worlds = set(["108998961191900167", "106926282333089592", "104841804583318464"])
msi = set(["110198981276611770", "107905284983470149", "105873410870441926"])
westernRegionTier1 = set(["LCK", "LPL"]) #4
westernRegionTier2 = set(["PCS", "VCS", "LJL"]) #1
easternRegionTier1 = set(["LCS", "LEC"]) #3
easternRegionTier2 = set(["CBLOL", "LLA"]) #1
semileague = ["CHALLENGER", "MASTERS", "ACADEMY"]


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
            'redGold': 0,

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
            'blueGold': 0,
            
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
    print(appearances)

    sortedElo = getFile(SORTED_CURRENT_ELO_PATH)
    maxElo = sortedElo[0]['elo']
    minElo = sortedElo[-1]['elo']

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
print(scoreMap)
writeFile(INITIAL_ELO_PATH, scoreMap)
        

        

        


