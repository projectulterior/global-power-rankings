from file import read, write
from elo import EloRating

games = read('../data/analysis/games.json')

def calculate_elos(k, default=1500):
    elos = {} # {[], []}
    for game in games:
        red = game['red_id']
        blue = game['blue_id']
        winner = game['winner']
        timestamp = game['start_time'] if 'end_time' in game else game['tournament_endDate'] +  'T17:40:29.537Z'

        redElo = elos[red][-1]['elo'] if red in elos else default
        blueElo = elos[blue][-1]['elo'] if blue in elos else default

        redElo, blueElo = EloRating(redElo, blueElo, k, 1 if winner == 'red' else 2)

        if red not in elos:
            elos[red] = [{"elo": redElo, "timestamp": timestamp}]
        else:
            elos[red].append({"elo": redElo, "timestamp": timestamp})

        if blue not in elos:
            elos[blue] = [{"elo": blueElo, "timestamp": timestamp}]
        else:
            elos[blue].append({"elo": blueElo, "timestamp": timestamp})
    return elos

write('data/elos.json', calculate_elos(4))
