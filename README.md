# PostgreSQL database migrator for Golang tests

## Overview

The library creates a test database and runs migrations or fixtures on it.  
The newly database name is provided together with a cleanup function.  
In case test run fails, the cleanup function could not be run, thus the database could not deleted and be available for more analysis.
The database naming convention is "t" plus unix time seconds plus the test name.

The migration info is written in a table which name can be overriden.

## How to use

### Migration files

Migration files naming convention should be "V0001__name.sql".  
This naming convention can be overridden by passing the desired regex.

### Fixture files

On top of migration files, fixture files can be applied by concentrating their paths in `MigrationFilePaths`.

### Template Fixture files

On top of normal fixture files, template fixtures can be used. This should come handy when specific IDs should be used across fixture files.  
The template fixture files should be added in `TemplateFilePaths` together with a render function.  
Any render engine could be used as long as the render function signature is met.

### How to integrate with your code

Database could be created as in the makefile target.  
For files and render function see external test file example.  
