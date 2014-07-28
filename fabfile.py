
from fabric.api import local, run, cd, settings, sudo, env, put
import os

env.hosts = ['ubuntu@s2p.flakiness.es']

PROJECT_DIR   = "/home/ubuntu/work/s2pd"

def init_app():
    pass

def rebuild():
    with settings(warn_only=True):
        local("rm s2pd")
    local("go build -o s2pd && go test")

def update():
    rebuild()
    run("mkdir -p " + PROJECT_DIR)
    with cd(PROJECT_DIR):
        put("s2pd", "s2pd", mode=0755)

def reload_daemons():
    with cd(PROJECT_DIR):
        with settings(warn_only=True):
            sudo("stop s2pd")
        run("mkdir -p " + PROJECT_DIR + "/logs")
        put("confs/nginx.conf", "nginx.conf")
        put("confs/upstart.conf", "upstart.conf")
        sudo("cp nginx.conf /etc/nginx/sites-enabled/s2pd.conf")
        sudo("cp upstart.conf /etc/init/s2pd.conf")
        sudo("/etc/init.d/nginx restart")
        sudo("start --verbose s2pd")

def deploy():
    update()
    reload_daemons()
