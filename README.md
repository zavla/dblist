## A DBA backup files tool set. ##

The package 'dblist' is used in databases backups management.
Sometimes a DBA needs to keep a directory with different backup files rotated.
Dblist/v2 package contains API that is needed when utilities work with config file. Config file controls the backup files management.

Files must obey naming scheme.
Database name in file name and file suffix define a file group.
Every file group may have its most recent files and outdated files.
Dblist deals with the task of selecting last backup files in each file group, according to a config file.
