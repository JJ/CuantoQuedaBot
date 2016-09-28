#!/usr/bin/env python
from fabric.api import local

def uptime():
  local('uptime')
