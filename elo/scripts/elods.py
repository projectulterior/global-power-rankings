import os
import json
import math
from enum import Enum


class SortBy(Enum):
    ASC = 'asc'
    DESC = 'desc'

class EloDBQuery:
    def __init__(self):
        self.ID = None
        self.startDate = None
        self.endDate = None
        self.minELO = None
        self.maxELO = None
        self.sortBy = SortBy.ASC
        self.limit = None
    
    def setTeamID(self, ID):
        self.ID = ID
        return self

    def setStartDate(self, startDate):
        self.startDate = startDate
        return self
    
    def setEndDate(self, endDate):
        self.endDate = endDate
        return self
    
    def setMinELO(self, minELO):
        self.minELO = minELO
        return self
    
    def setMaxELO(self, maxELO):
        self.maxELO = maxELO
        return self
    
    def setSortBy(self, sortBy):
        self.sortBy = sortBy
        return self

class EloDB:
    def __init__(self, ):
        self.eloDB = {}
        self.gameCount = {}
        self.decayFn = lambda x: 0 if x not in self.gameCount else -self.gameCount[x] / 100

    def get(self, ELODBQuery):
        if ELODBQuery.teamID is None:
            raise Exception('Team ID must be provided')

        if ELODBQuery.teamID not in self.eloDB:
            return []
        
        eloHistory =  self.eloDB[ELODBQuery.teamID]

        # binary search to find first occurrence of startDate, or first occurrence after startDate
        startIndex = 0
        if ELODBQuery.startDate is not None:
            left = 0
            right = len(eloHistory) - 1
            while left <= right:
                mid = (left + right) // 2
                if eloHistory[mid]['timestamp'] < ELODBQuery.startDate:
                    left = mid + 1
                else:
                    right = mid - 1
            startIndex = left

        # binary search to find first occurrence of endDate, or first occurrence before endDate
        endIndex = len(eloHistory) - 1
        if ELODBQuery.endDate is not None:
            left = 0
            right = len(eloHistory) - 1
            while left <= right:
                mid = (left + right) // 2
                if eloHistory[mid]['timestamp'] < ELODBQuery.endDate:
                    left = mid + 1
                else:
                    right = mid - 1
            endIndex = right
        
        # filter by minELO
        if ELODBQuery.minELO is not None:
            eloHistory = list(filter(lambda x: x['elo'] >= ELODBQuery.minELO, eloHistory))
        
        # filter by maxELO
        if ELODBQuery.maxELO is not None:
            eloHistory = list(filter(lambda x: x['elo'] <= ELODBQuery.maxELO, eloHistory))
        
        # sort by timestamp
        eloHistory.sort(key=lambda x: x['timestamp'], reverse=ELODBQuery.sortBy == SortBy.DESC)

        # limit
        if ELODBQuery.limit is not None:
            eloHistory = eloHistory[:ELODBQuery.limit]

        return eloHistory[startIndex:endIndex + 1]
    
    def getCurrent(self, teamID):
        if teamID not in self.eloDB:
            return None
        
        return self.eloDB[teamID][-1]
    
    def getAll(self):
        return self.eloDB

    def insert(self, ID, timestamp, elo, metadata=None):
        if ID not in self.eloDB:
            self.eloDB[ID] = []
        
        if len(self.eloDB[ID]) > 0 and timestamp < self.eloDB[ID][-1]['timestamp']:
            raise Exception('Timestamp must be greater than the last timestamp')

        self.eloDB[ID].append({
            'timestamp': timestamp,
            'elo': elo,
            'metadata': metadata
        })

    def flush(self, path):
        data_bytes = json.dumps(self.eloDB).encode('utf-8')
        fd = os.open(path, os.O_WRONLY | os.O_CREAT)
        os.write(fd, data_bytes)
        os.close(fd)

class Team(Enum):
    RED = 'red'
    BLUE = 'blue'

class EloGenerator:
    def generate(self, player, opponent, win=True, kFactor=None):
        expectedPlayer = 1 / (1 + 10 ** ((player - opponent) / 400))

        newPlayerElo = 0
        if win:
            newPlayerElo = player + kFactor * (1 - expectedPlayer)
        else:
            newPlayerElo = player + kFactor * (0 - expectedPlayer)
        
        return newPlayerElo

# elo data format
# {
#     'NA': [
#         {
#             'timestamp': 01/01/2001,
#             'elo': 1500
#             'metadata: anything
#         },
#         {
#             'timestamp': 01/01/2001, 
#             'elo': 1500
#         }
#     ]
# }

# team elo format
# {
#     'TEAMID': [
#         {
#             'date': 01/01/2001,
#             'elo': 1500
#         },
#         {
#             'date': 01/01/2001,
#             'elo': 1500
#         }
#     ]
# }
