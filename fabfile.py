#!/usr/bin/env python
from fabric.api import env, run

env.hosts = [ '159.100.248.62' ]
env.user = "root"

def uptime():
  run('uptime')
