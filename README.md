## A DBA tool set for manipulating backup files. ##

The package 'dblist/v3' is used in databases backups files management.

Dblist/v3 deals with the task of selecting last backup files in each file group, according to a config file.

Sometimes a DBA needs to keep a directory with different databases backup files rotated.
Dblist/v3 package contains API that is used by utilities that work with such files or with rotation configuration file.

Filenames must obey naming scheme.
Database name in file name and file suffix define a file group of a database backup.
Every file group may have its most recent files and outdated files.
