#! /usr/bin/env bash
#
# deploy-db.sh generates usda.db and syncs with remote
# for the rare occassion this is needed.

# TODO gen db
ssh grubprint.io -C 'sudo systemctl stop grubprint@daniel.service'
rsync -P -z usda.db grubprint.io:usda.db
ssh grubprint.io -C 'sudo systemctl start grubprint@daniel.service'
