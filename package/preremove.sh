#!/bin/bash
#

systemctl stop broadtail
systemctl disable broadtail

rm -fr /var/lib/broadtail