#!/bin/bash

VERSION=$1

dnf install -y epel-release
dnf install -y mock

tar -xf warewulf-${VERSION}.tar.gz warewulf-${VERSION}/warewulf.spec

mock -r epel-8-x86_64 --rebuild --spec=warewulf-${VERSION}/warewulf.spec --sources=.
mv /var/lib/mock/epel-8-x86_64/result/warewulf-${VERSION}-*.el8.x86_64.rpm .

mock -r epel-7-x86_64 --rebuild --spec=warewulf-${VERSION}/warewulf.spec --sources=.
mv /var/lib/mock/epel-7-x86_64/result/warewulf-${VERSION}-*.el7.x86_64.rpm .