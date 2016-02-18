#cloud-config

# attempt to disable everything except ssh and the initial yum update

cloud_init_modules:
  - users-groups
  - ssh

cloud_config_modules:
  - yum-configure
  - yum-add-repo
  - package-update-upgrade-install

packages: []

cloud_final_modules: []

bootcmd: []

runcmd: []

