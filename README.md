# PostgreSQL database migrator for Golang tests

## Overview

The library creates a test database and runs migrations or fixtures on it.  
In case test run fails, the database is not deleted. The database naming convention is "t" plus unix time seconds plus the test name.

The migration info is written in a table which name can be overriden.

## How to use

### Prepare migration files

Migration files naming convention should be "V0001__name.sql".  
This naming convention can be overridden by passing the desired regex.

### How to integrate with your code

Database could be created as in the makefile target.  
For tests see external test file.  
