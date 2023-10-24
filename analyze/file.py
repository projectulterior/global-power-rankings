import json
import os

def read(filename):
    f = open(filename)
    data = json.load(f)
    f.close()
    return data

def write(filename, data):
    b = json.dumps(data, indent=4).encode('utf-8')
    fd = os.open(filename, os.O_WRONLY | os.O_CREAT)
    os.write(fd, b)
    os.close(fd)