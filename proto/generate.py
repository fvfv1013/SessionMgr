import os
import subprocess

PROTO_DIR = "./proto"
OUT_DIR = "../../"

if not os.path.exists(OUT_DIR):
    os.makedirs(OUT_DIR)

for proto_file in os.listdir(PROTO_DIR):
    if proto_file.endswith(".proto"):
        subprocess.run(["protoc", f"-I={PROTO_DIR}", f"--go_out={OUT_DIR}", os.path.join(PROTO_DIR, proto_file)])