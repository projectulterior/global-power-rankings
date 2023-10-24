from elods import EloDBQuery, EloDB, SortBy, EloGenerator
from util import getFile, writeFile
from team import TEAM_ELO_PATH

teamELOS = getFile(TEAM_ELO_PATH)
teamELODB = EloDB(teamELOS)

sortedELOS = teamELODB.getCurrentAll()
sortedELOS.sort(key=lambda x: x['elo'], reverse=True)

writeFile('./data/elo/sorted_elo.json', sortedELOS)