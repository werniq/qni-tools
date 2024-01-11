import hashlib
import socket
import subprocess
import os
import shutil

host = '0.0.0.0'
def_port = 1248

def get_hash(path):
    hasher = hashlib.sha256()
    with open(path, 'rb') as f:
        hasher.update(f.read())
    return hasher.hexdigest()


def listen_to_client(addr: str = host, port: int = def_port):
    server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server.bind((addr, port))
    server.listen()
    print(f"Listening on {addr}:{port}.")
    return server


def server(action: bool):
    if action:
        # Start the server
        print("Starting server")
        server_socket = listen_to_client()
        try:
            while True:
                conn, addr = server_socket.accept()
                handle_connection(conn)
        except KeyboardInterrupt:
            server_socket.close()

    else:
        subprocess.run(["sudo", "pkill", "chatting"])


def handle_connection(conn):
    remote_addr = conn.getpeername()[0]
    print(f"[*] Client connected from {remote_addr}")

    file_loc = conn.recv(1024).decode().strip()
    data_file = os.path.join(file_loc, "data")
    conn.send(get_hash(data_file).encode() + b'\n')

    incoming_data = b""
    while True:
        data = conn.recv(1024)
        if not data:
            break
        incoming_data += data

    if incoming_data:
        with open(data_file, 'wb') as f:
            f.write(incoming_data)

    print(f"Client at {remote_addr} disconnected.")
