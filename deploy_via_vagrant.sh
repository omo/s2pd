#!/bin/sh

vagrant ssh -c "cd /vagrant && source ./venv/bin/activate && fab deploy"
