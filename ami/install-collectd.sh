#!/bin/bash
set -e

cd /tmp

# Install dependencies
yum install -y gcc rpm-build yum-utils rpmdevtools epel-release

# Fetch source
curl -O http://collectd.org/files/collectd-5.5.0.tar.bz2

# Extract source
tar xvf collectd-5.5.0.tar.bz2

# Use our Amazon Linux friendly RPM spec
mv /home/ec2-user/collectd.spec.amzn1 collectd-5.5.0/contrib/redhat/collectd.spec

# Create build directory
mkdir -p /root/rpmbuild/SOURCES/

# Copy source
cp collectd-5.5.0.tar.bz2 /root/rpmbuild/SOURCES/

# Install dependencies of the spec file
yum-builddep -y collectd-5.5.0/contrib/redhat/collectd.spec

# Create RPMS
rpmbuild -bb collectd-5.5.0/contrib/redhat/collectd.spec

# Install RPMS
rpm -i /root/rpmbuild/RPMS/x86_64/collectd-5.5.0-1.amzn1.x86_64.rpm
rpm -i /root/rpmbuild/RPMS/x86_64/collectd-write_http-5.5.0-1.amzn1.x86_64.rpm
rpm -i /root/rpmbuild/RPMS/x86_64/collectd-python-5.5.0-1.amzn1.x86_64.rpm

# Enable collectd
chkconfig collectd on

# Cleanup
rm -rf /root/rpmbuild
rm -rf collectd-5.5.0*

