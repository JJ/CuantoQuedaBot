#!/usr/bin/env python
from fabric.api import env, run
from fabric.context_managers import shell_env, cd

env.hosts = [ '159.100.248.62' ]
env.user = "root"
env.release_path= "CuantoQuedaBot"

def uptime():
    run('uptime')
  
  
def build(goroot="/usr/lib/go", gobin="/usr/bin", gopath="/home/%s/lib/go" % env.user):
    variables = { 'goroot': goroot, 'gobin':gobin,'gopath': gopath, 'release_path':env.release_path} 
    with cd(env.release_path):
        with shell_env(GOROOT=goroot,GOPATH=gopath,GOBIN=gobin):
            run("git pull" )
            run("$GOBIN/go get")

def start(goroot="/usr/lib/go", gobin="/usr/bin", gopath="/home/%s/lib/go" % env.user):
    variables = { 'goroot': goroot, 'gobin':gobin,'gopath': gopath, 'release_path':env.release_path} 
    run("export GOROOT=%(goroot)s;export GOPATH=%(gopath)s;export GOBIN=%(gobin)s;cd %(release_path)s; $GOBIN/go run CuantoQuedaBot" %  variables )
    
