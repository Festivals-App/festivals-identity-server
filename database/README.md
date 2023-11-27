The database contains all users that have acces to the FestivalsApp backend or parts of the FestivalApp bakend. The default user is called administrator with the password set to 'we4711. It alwo contains all service node keys.


# Local development macOS

First you need to [install](https://www.novicedev.com/blog/how-install-mysql-macos-homebrew) mysql on your development machine.


Staring and logging into mysql
```
brew services start mysql
mysql -uroot
```

Logout and stopping mysql
```
> EXIT;
brew services stop mysql
```

### MYSQL cheatsheet
```
SHOW DATABASES;

```