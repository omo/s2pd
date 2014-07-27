
from fabric.api import local, run, cd, settings, sudo, env, put
import os

PROJECT_DIR   = "/home/ubuntu/work/s2pd"

def init_app():
    run("mkdir -p " + PROJECT_DIR)
    with cd(PROJECT_DIR):
        run("virtualenv --distribute venv")

def rebuild():
    with settings(warn_only=True):
        local("rm s2pd")
    local("go build -o s2pd && go test")

def update():
    rebuild()
    with cd(PROJECT_DIR):
        put("s2pd", "s2pd")
