

# Deployment notes:

Basic idea:

 * The binary is built using local Vagrant box.
 * Then the binary is pushed to prod using Fabric.
 * The prod has nginx, logrotate, upstart as given.

Setup Vagrant:

 * `Vagrantfile` does it for you.
 * Don't forget to enable ssh forwarding agent.
   * http://wildlyinaccurate.com/using-ssh-agent-forwarding-with-vagrant

Setup prod:

 * Not clear :-(
 * Needs nginx, logrotate, upstart at least. (Hopefully that's all)

Push:

 * `vagrant up`.
 * run `deploy_via_vagrant.sh`
