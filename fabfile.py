#!/usr/bin/env python
from fabric.api import env, run
from fabric.context_managers import shell_env, cd
import os

env.hosts = [ '159.100.248.62' ]
env.user = "bot"
env.release_path= "CuantoQuedaBot"

def uptime():
    run('uptime')
  
  
def build(gobin="/home/%s/lib/go/bin"% env.user, gopath="/home/%s/lib/go" % env.user):
    with cd(env.release_path):
        with shell_env(GOPATH=gopath,GOBIN=gobin):
            run("git pull" )
            run("go get")

def start(goroot="/usr/lib/go", gobin="/usr/bin",
          gopath="/home/%s/lib/go" % env.user,
          bot_token=os.environ['BOT_TOKEN'],
          papertrail_host=os.getenv('PAPERTRAIL_HOST',''),
          papertrail_port=os.getenv('PAPERTRAIL_PORT',''),
          logz_host=os.getenv('LOGZ_HOST',""),
          logz_token=os.getenv('LOGZ_TOKEN',"")):
    
    with cd(env.release_path):
        with shell_env(GOROOT=goroot,GOPATH=gopath,
                       GOBIN=gobin,BOT_TOKEN=bot_token,
                       PAPERTRAIL_HOST=papertrail_host,
                       PAPERTRAIL_PORT=papertrail_port,
                       LOGZ_HOST=logz_host,
                       LOGZ_TOKEN=logz_token):
            run("supervisorctl start CuantoQuedaBot" )
            
def stop():
    run("supervisorctl stop CuantoQuedaBot")
    
