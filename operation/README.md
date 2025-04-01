# Operation

The `operation` directory contains all configuration templates and scripts to install and run the festvials-identity-server.

* `backup.sh` script to backup the database to run periodically from a cron job
* `install.sh` script to install the festivals-identity-server and the database on a VM
* `secure-mysql.sh` script to secure the intitial mysql installation
* `service_template.service` festivals-identity-server unit file for `systemctl`
* `ufw_app_profile` firewall app profile file for `ufw`
* `update.sh` script to update the festivals-identity-server

## Deployment

Follow the [**deployment guide**](DEPLOYMENT.md) for deploying the festivals-identity-server inside a virtual machine or the [**local deployment guide**](./local/README.md) for running it on your macOS developer machine.
