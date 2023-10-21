import json
import pandas as pd

FILES = [
    'leagues.json',
    'player.json',
    'teams.json',
    'tournaments.json',
]

def parse(filename):
    f = open(filename)
    data = json.load(f)
    return data

leagues = parse('data/leagues.json')
players = parse('data/player.json')
teams = parse('data/teams.json')
tournaments = parse('data/tournaments.json')

# print("leagues", leagues)
# print("players", players)
# print("teams", teams)
# print("tournaments", tournaments)

# for tournament in tournaments:
#     print(tournament["id"])

t = pd.read_json(open("data/tournaments.json"))
# print(t)

df = pd.json_normalize(tournaments, record_path=["stages"], meta=list(t.columns), meta_prefix="tournament.", record_prefix="stage.")
# print(df)

df = df.drop("tournament.stages", axis=1)
# print(df)


exploded = df.explode("stage.sections")
print(exploded)

# df = pd.json_normalize(exploded, record_path=[1, "stage.sections"])
# print(df)
# 0/0

for section in exploded.iterrows():
    row = section[1]
    # print(row)

    sections = row["stage.sections"]
    # print(sections)

    df = pd.json_normalize(sections)
    print(df)

    matches = df.explode("matches")
    print("hello", matches)

    for match in matches.iterrows():
        m = match[1]["matches"]

        teams = pd.json_normalize(m["teams"])
        # print("teams\n", teams)

        games = pd.json_normalize(m["games"])
        print("games\n", games)

        teams = games.explode("teams")
        print("teams\n", teams)

        


    0/0

    # print(len(section))
    # s = section[1]
    # print(s)

    # df = pd.DataFrame(s)
    # print(df)
    # sec = pd.DataFrame(section)
    # print(sec)

# sections = pd.json_normalize(exploded, record_path=["stage.sections"])
# print(sections)
