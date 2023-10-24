import os
import json
import math
import trueskill
from scipy.stats import norm

def getFile(path):
    fd = os.open(path, os.O_RDONLY)
    data = os.read(fd, os.path.getsize(path))
    data = json.loads(data)
    os.close(fd)
    return data

def writeFile(path, data):
    if os.path.exists(path):
        os.remove(path)
    
    data_bytes = json.dumps(data).encode('utf-8')
    fd = os.open(path, os.O_WRONLY | os.O_CREAT)
    os.write(fd, data_bytes)
    os.close(fd)

# omega is how fast the decay is
def regularize(trueSkill, numMatches, avg, tau, omega):
    # exponential decay for numMatches
    ex = -math.pow(omega, numMatches)
    term = 1 - math.exp(ex / tau)
    return (1 - term) * avg + term * trueSkill

def sigmaRegularization(mu, sigma, avg, tau):
    ex = -math.pow(sigma, 2)
    term = math.exp(ex / tau)
    return (1 - term) * avg + term * mu


def winProb(team1: trueskill.Rating, team2: trueskill.Rating):
    deltaMu = team1.mu - team2.mu
    sumSigma = team1.sigma ** 2 + team2.sigma ** 2
    denom = math.sqrt(2 * (trueskill.BETA * trueskill.BETA) + sumSigma)
    return norm.cdf(deltaMu / denom)