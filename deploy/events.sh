#!/bin/bash -e
#########################################################################
# this script handles LIFECYCLE_EVENTs from aws code deploy
# (see appspec.yml hooks)
#########################################################################
event=${LIFECYCLE_EVENT:-""} #injected by code deploy
case $event in
  "AfterInstall") systemctl enable corkboard ;;
  "ApplicationStart") systemctl start corkboard ;;
  "ApplicationStop") systemctl stop corkboard ;;
  *) echo "Can not handle unknown event: ($event)" 1>&2 && exit 1 ;;
esac
