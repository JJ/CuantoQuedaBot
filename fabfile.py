#!/usr/bin/env python
from fabric.api import env, run
from fabric.context_managers import shell_env, cd
import os

env.hosts = [ '159.100.248.62' ]
env.user = "root"
env.release_path= "CuantoQuedaBot"

def uptime():
    run('uptime')
  
  
def build(goroot="/usr/lib/go", gobin="/usr/bin", gopath="/home/%s/lib/go" % env.user):
    with cd(env.release_path):
        with shell_env(GOROOT=goroot,GOPATH=gopath,GOBIN=gobin):
            run("git pull" )
            run("$GOBIN/go get")

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
            run("echo $GOPATH; nohup $GOBIN/go run CuantoQuedaBot.go >& bot.log < /dev/null &", pty=False )
            
def stop():
    run("pkill CuantoQuedaBot")
    
