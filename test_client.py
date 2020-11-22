import socket

client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
client.connect(('127.0.0.1', 7456))

tmpi = 0
while True:
    from_server = client.recv(4096)
    if from_server != b'':
        print(from_server)
        client.send(b"{\"success\": true}")
        # tmpi = 0
    # else:
    #     tmpi += 1
    #     if tmpi > 20:
    #         client.send(b"Something tells me you're closed. Is it true?")

client.close()
