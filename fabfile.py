
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
        put("s2pd", "s2pd.fresh", mode=0755)

def reload_daemons():
    with cd(PROJECT_DIR):
        with settings(warn_only=True):
            sudo("stop s2pd")
        run("cp s2pd.fresh s2pd")
        run("mkdir -p " + PROJECT_DIR + "/logs")
        put("confs/nginx.conf", "/etc/nginx/sites-enabled/s2pd.conf", use_sudo=True)
        put("confs/upstart.conf", "/etc/init/s2pd.conf", use_sudo=True)
        put("confs/logrotate.conf", "/etc/logrotate.d/s2pd", use_sudo=True)
        sudo("/etc/init.d/nginx restart")
        sudo("start --verbose s2pd")

def deploy():
    update()
    reload_daemons()
