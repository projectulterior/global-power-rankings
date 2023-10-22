import os
import json

def getFile(path):
    fd = os.open(path, os.O_RDONLY)
    data = os.read(fd, os.path.getsize(path))
    data = json.loads(data)
    os.close(fd)
    return data

def writeFile(path, data):
    os.remove(path)
    data_bytes = json.dumps(data).encode('utf-8')
    fd = os.open(path, os.O_WRONLY | os.O_CREAT)
    os.write(fd, data_bytes)
    os.close(fd)