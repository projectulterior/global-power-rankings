import sys
from file import read, write
import os
import json
import pandas as pd

FILES = [
    'leagues.json',
    'player.json',
    'teams.json',
    'tournaments.json',
]

OUTPUT = "data/games.json"

leagues = read('data/leagues.json')
players = read('data/player.json')
teams = read('data/teams.json')
tournaments = read('data/tournaments.json')

def update(dest, obj, prefix):
    for key, value in obj.items():
        dest[prefix+key] = value
    return dest

def parseTeam(team):
    t = {}
    # players
    for player in team["players"]:
        t[player["role"]] = player["id"]
    return t

def copy(obj, excludes=[]):
    ret = obj.copy()
    for exclude in excludes:
        del ret[exclude]
    return ret

# games
flat = []
skipped = 0
for tournament in tournaments:
    t = copy(tournament, ['stages'])

    stages = tournament['stages']
    for stage in stages:
        s = copy(stage, ['sections'])

        sections = stage['sections']
        for section in sections:
            sec = copy(section, ['matches', 'rankings'])

            matches = section['matches']
            for match in matches:
                m = copy(match, ['games', 'teams', "strategy"])
                update(m, match['strategy'], 'strategy_')

                games = match['games']
                for game in games:
                    if game['state'] != 'completed':
                        skipped += 1
                        continue

                    ## collect data on each game ##

                    update(game, t, "tournament_")
                    update(game, s, "stage_")
                    update(game, sec, "section_")
                    update(game, m, "match_")


                    id = game["id"]
                    state = game["state"]

                    # total number of games in this match
                    game["total"] = len(games)

                    teams = game["teams"]
                    assert len(teams) == 2

                    # red/blue teams ids
                    red, blue = None, None
                    for team in teams:
                        if team["side"] == "red":
                            game["red_id"] = team["id"]
                            red = team
                        elif team["side"] == "blue":
                            game["blue_id"] = team["id"]
                            blue = team
                        else:
                            print("unknown side", id, state)
                    
                    # team data
                    for team in match["teams"]:
                        if red and team["id"] == red["id"]:
                            update(game, parseTeam(team), "red_")
                        if blue and team["id"] == blue["id"]:
                            update(game, parseTeam(team), "blue_")

                    # winner
                    winner, isForfeit = None, False
                    if game["state"] == "completed":
                        if red["result"]["outcome"] == "win" and (blue["result"]["outcome"] == "loss" or blue["result"]["outcome"] == "forfeit"):
                            winner = "red"
                            if blue["result"]["outcome"] == "forfeit":
                                isForfeit = True
                        elif (red["result"] ["outcome"] == "loss" or red["result"]["outcome"] == "forfeit") and blue["result"]["outcome"] == "win": 
                            winner = "blue"
                            if red["result"]["outcome"] == "forfeit":
                                isForfeit = True
                        else:
                            print("unknown winner", id, state)
                    game["winner"] = winner
                    game["is_forfeit"] = isForfeit                    

                    # append game
                    del game["teams"]
                    flat.append(game)

print("parsed:", len(flat))
print("skipped:", skipped)
print("total:", len(flat)+skipped)       

write(OUTPUT, flat)