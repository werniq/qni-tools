import json
import os
import subprocess
import getpass

class Config:
    def __init__(self, target_ip, target_dir):
        self.targetIp = target_ip
        self.targetDir = target_dir


def create_json(config, json_file):
    with open(json_file, 'w') as file:
        json.dump(config.__dict__, file, indent=2)


def initialize():
    current_user = getpass.getuser()

    if os.getuid() == 0:
        chat_folder = os.path.expanduser("~") + "/.chat"

        if not os.path.exists(chat_folder):
            subprocess.run(["sudo", "chatting", "stop"])
            print("[*] Chatting service stopped")

            config = Config(
                target_ip=input("[I] Target IP: "),
                target_dir=input("[I] Target Directory: ")
            )

            with open("/etc/.chat/syncList", 'a') as sync_list_file:
                sync_list_file.write(os.getcwd() + "\n")

            os.makedirs(chat_folder, exist_ok=True)
            json_file_path = os.path.join(chat_folder, "config.json")
            create_json(config, json_file_path)

            print(f"[*] Added '{os.getcwd()}' in '/etc/.chat/syncList'")
            print(f"[*] Created {chat_folder}")
            print(f"[*] Created {json_file_path}")
            print(" [*] Chatting Service Started Again")

            subprocess.run(["sudo", "chatting", "start"])

        else:
            print(f"[!] {os.getcwd()} is already a [chat] watched file")

    else:
        print("[!] You need to be root to execute this command")

