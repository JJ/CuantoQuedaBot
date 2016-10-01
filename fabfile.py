#!/usr/bin/env python
from fabric.api import env, run

env.hosts = [ '159.100.248.62' ]
env.user = "root"
env.release_path= "CuantoQuedaBot"

def uptime():
    run('uptime')
  
  
def build(goroot="/usr/lib/go", gobin="/usr/bin", gopath="/home/%s/lib/go" % env.user):
    variables = { 'goroot': goroot, 'gobin':gobin,'gopath': gopath, 'release_path':env.release_path} 
    run("cd %(release_path)s;git pull" % variables )
    run("export GOROOT=%(goroot)s;export GOPATH=%(gopath)s;export GOBIN=%(gobin)s;cd %(release_path)s; $GOBIN/go get" %  variables )
