# Running the festivals identity server locally on you mac

This guide provides instructions for setting up and running the identity server on your macOS machine. Whether you're a new developer or setting up a fresh environment, you'll find everything needed to install dependencies, configure the project, and start development efficiently.  

Before proceeding, ensure you have the required tools installed and follow the steps below to get your local environment up and running smoothly.  

## Prerequisites

As all festivalsapp services communicate based on DNS names you need to add some entries to your `/etc/hosts` file.

```ini
# local development on this machine
127.0.0.1       identity.festivalsapp.dev
127.0.0.1       discovery.festivalsapp.dev
```

The identity server stores data inside a database. At the moment that database needs to be running on the same host and be reachable at `localhost:3306`.
Let's use [homebrew](https://brew.sh/) to install mysql and start it.

```bash
brew install mysql
brew services start mysql
```

Now we can create the database and the database users. The creation script will automatically create a default admin user, service key and API key and those we will use those throughout local development.

```bash
mysql -e "source <path/to/project/folder/festivals-identity-server/database/create_database.sql>"
mysql -e "CREATE USER 'festivals.identity.writer'@'localhost' IDENTIFIED BY 'we4711';"
mysql -e "GRANT SELECT, INSERT, UPDATE, DELETE ON festivals_identity_database.* TO 'festivals.identity.writer'@'localhost';"
mysql -e "CREATE USER 'festivals.identity.backup'@'localhost' IDENTIFIED BY 'we4711';"
mysql -e "GRANT ALL PRIVILEGES ON festivals_identity_database.* TO 'festivals.identity.backup'@'localhost';"
mysql -e "FLUSH PRIVILEGES;"
```

## Running the festivals identity server

This project uses Make to streamline local setup and execution. The Makefile includes commands for installing dependencies, configuring the environment, and running the service. Using Make ensures a consistent workflow and simplifies common tasks.

1. First you need to build the binary for local development using the `build` command.

    ```bash
    make build
    ```

2. By invoking the `install` command Make will install the newly build binary and all files it needs to run. The default install path is a folder inside your users container folder at `~/Library/Containers/org.festivalsapp.project`, this is so you don't need to use `sudo` to install and run the website node.

    ```bash
    make install
    ```

3. Now you can run the binary by issuing the `run` command. This will run the binary with the `--container="~/Library/Containers/org.festivalsapp.project"` option, telling the binary that the config file will be located at `~/Library/Containers/org.festivalsapp.project/usr/local/festivals-website-node` instead of the default `/usr/local/festivals-website-node`.

    ```bash
    make run
    ```

4. To prevent annoying error messages you should run the [FestivalsApp Gateway](https://github.com/Festivals-App/festivals-gateway) service. You can do that with the `run-env` command but in order for the command to work you need to run the `install` commmand at least once for the gateway service. To stop the gateway service you can use the `stop-env` command.

```bash
make run-env
make stop-env
```

## Testing

The festivals identity server is now reachable on your machine at `https://identity.festivalsapp.dev:22580`.

Lets login as the default admin user using the client certificate and get the server info:

```bash
curl -H "Api-Key: TEST_API_KEY_001" -u "admin@email.com:we4711" --cert /opt/homebrew/etc/pki/issued/client.crt --key /opt/homebrew/etc/pki/private/client.key --cacert /opt/homebrew/etc/pki/ca.crt https://identity.festivalsapp.dev:22580/users/login
```

This should return a JWT Token `<Header.<Payload>.<Signatur>`

  > eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.
  > eyJVc2VySUQiOiIxIiwiVXNlclJvbGUiOjQyLCJVc2VyRmVzdGl2YWxzIjpbXSwiVXNlckFydGlzdHMiOltdLCJVc2VyTG9jYXRpb25zIjpbXSwiVXNlckV2ZW50cyI6W10sIlVzZXJMaW5rcyI6W10sIlVzZXJQbGFjZXMiOltdLCJVc2VySW1hZ2VzIjpbXSwiVXNlclRhZ3MiOltdLCJpc3MiOiJpZGVudGl0eS0wLmZlc3RpdmFsc2FwcC5ob21lIiwiZXhwIjoxNzQwMjMxMTQ4fQ.
  > geBq1pxEvqwjnKA5YTHQ8IjJc9mwkpsQIRy1kGc63oNXzyAhPrPJsepICXxr2yVmB0E8oDECXLn4Cy5V_p4UAduWXnc0r8S05ijV8NCfmsEcJg-RRO8POkGykiC2mrn-XR8Nf8OF0fLp7Mhsb0_aqBoTOLdtB9V7IV49-JjWwX5gHl3HuXGOOhe4n_epumc8w8yDxYakWeaBFtEtaRmhFXK_yttexYOLP6Z1BBTL005hBGhO58qVW0cfgf_t5VWBpUnz3zqdC-GFegItqJQbKZ2pmfmXNz_AoJf2JmPtCzpJ4lG6QeSslvdFuwaCdYpDQPOvnMSIORwrAq_FL2m7qw

Use this token to make authorized calls to the identity server:

```bash
curl -H "Authorization: Bearer <JWT>" --cert /opt/homebrew/etc/pki/issued/client.crt --key /opt/homebrew/etc/pki/private/client.key --cacert /opt/homebrew/etc/pki/ca.crt https://identity.festivalsapp.dev:22580/info
