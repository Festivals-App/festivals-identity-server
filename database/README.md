# Database deployment

The database contains all users that have access to the FestivalsApp backend or parts of the FestivalApp backend. The default user is called administrator with the password set to 'we4711'. It also contains all service keys and API keys.

## Local development macOS

First you need to [install](https://www.novicedev.com/blog/how-install-mysql-macos-homebrew) and configure mysql on your development machine.

```bash
brew install mysql
mysql_secure_installation
```

Staring and logging into mysql

```mysql
brew services start mysql
mysql -uroot -p
```

Logout and stopping mysql

```mysql
exit;
brew services stop mysql
```

## Server deployment

The install script will install and secure the database.

### MYSQL cheatsheet

```mysql
SHOW DATABASES;
USE database;
SELECT * FROM table
```
