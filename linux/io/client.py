import socket

s = socket.socket()
s.connect(('127.0.0.1', 8000))
s.send("hello".encode())
s.close()