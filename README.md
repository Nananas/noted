# Noted

> Minimal self-hosting notebook web-app in Go.

Noted is a minimal notebook web-app that can be self-hosted. 
It has no database, all notes are saved as files.
Sessions areA SHA256 salted hash is used for password protection, just in case.

## Usage
A precompiled binary for Linux amd64 is available to [download](https://github.com/Nananas/noted/releases).

```
./noted -h
```

### Account setup

First, create a file called `.accounts` in the working directory, containing `[user]|[pass]` lines:
```
test|testing
user|pass123
```

### Account hashing
Next, run the following command to convert this file to salt and hash the passwords:
```
./salthash
```

This will create a `accounts` file and a `users` file, containing the obscured account information and a list of the recognised users respectively.
(The `users` file functions as a backup/reminder). If everything went well, the `.accounts` file will then be removed.

### Running
Running:

```
./noted
```
or, for printing debug to stdout instead of `logfile.log`:

```
./noted -d
```

## Building from source
Dependencies:
	- [ymake](https://github.com/Nananas/ymake) (Optional)

Getting the code:
```
go get github.com/nananas/noted
# or
git clone https://github.com/nananas/noted.git
```

Building, creating a test account, and running:
```
ymake
```
