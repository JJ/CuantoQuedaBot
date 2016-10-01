#!/usr/bin/env python
from fabric.api import env, run

env.hosts = [ '159.100.248.62' ]
env.user = "root"
env.release_path= "CuantoQuedaBot"

def uptime():
    run('uptime')
  
  
def build(goroot="/usr", gopath="/home/%s/lib/Go" % env.user):
    run("export GOROOT=%(goroot)s;export GOPATH=%(gopath)s;export GOBIN=$GOPATH/bin;cd %(release_path)s;echo $GOBIN; echo $GOPATH; $GOBIN/go get" % { 'goroot': goroot, 'gopath': gopath, 'release_path':env.release_path} )
