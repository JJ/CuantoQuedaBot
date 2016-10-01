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
          bot_token=os.environ('BOT_TOKEN'),
          papertrail_host=os.environ('PAPERTRAIL_HOST'),
          papertrail_port=os.environ('PAPERTRAIL_PORT')):
    
    with cd(env.release_path):
        with shell_env(GOROOT=goroot,GOPATH=gopath,
                       GOBIN=gobin,BOT_TOKEN=bot_token,
                       PAPERTRAIL_HOST=papertrail_host,
                       PAPERTRAIL_PORT=papertrail_port):
            run(" $GOBIN/go run CuantoQuedaBot" %  variables )
            
def stop():
    run("pkill CuantoQuedaBot")
    
