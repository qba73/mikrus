[![Go Reference](https://pkg.go.dev/badge/github.com/qba73/mikrus.svg)](https://pkg.go.dev/github.com/qba73/mikrus)
[![Go Report Card](https://goreportcard.com/badge/github.com/qba73/mikrus)](https://goreportcard.com/report/github.com/qba73/mikrus)
[![Tests](https://github.com/qba73/mikrus/actions/workflows/go.yml/badge.svg)](https://github.com/qba73/mikrus/actions/workflows/go.yml)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/qba73/mikrus)
![GitHub](https://img.shields.io/github/license/qba73/mikrus)

# mikrus

`mikrus` is a CLI app for [MIKRUS VPS](https://mikr.us). It allows you to:

- show information about your server
- list provisioned servers
- restart your server
- check last log messages
- boost your server performance by turning on Amfetamina functionality
- show config for your DB(s) (Postgres, MySQL)
- execute remote commands on your server
- show usage / resource utilization statistics
- show ports assigned to your server
- show cloud functions assigned to your account (including stats)
- add / change domain assigned to your server

## Installing the command-line client

Install the client binary

```shell
go install github.com/qba73/mikrus/cmd/mikctl@latest
```

## Using the command-line client

To see help on using the client, run:

```shell
mikctl -h
```

## Setting your API key and Server ID

To use the client with your Mikrus account, you will need the API Key and Server ID provisioned in your Mikrus account. Go to the [Mikrus page](https://mikr.us/#pricing), sign up for the service. When your account is ready, go to the panel page and get your `server ID` and corresponding `API key`.

There are three ways to pass your API key to the client: in a config file, in an environment variable, or on the command line.

### In a config file

The `mikctl` client will read a config file named `.mikrus.yaml` (or `.mikrus.json`, or any other extension that Viper supports) in your home directory, or in the current directory.

For example, you can put your API key and sever ID in the file `$HOME/.mikrus.yaml`, and `mikctl` will find and read them automatically (replace `XXX` with your own API key, and `YYY` with your server ID):

```yaml
apiKey: XXX
srvID: YYY
```

### In an environment variable

`mikctl` will look for the API key and server ID in an environment variable named MIKRUS_API_KEY and MIKRUS_SRV_ID:

```shell
export MIKRUS_API_KEY=XXX
export MIKRUS_SRV_ID=YYY
mikctl ...
```

### On the command line

You can also pass your API key and server ID to the `mikctl` client using the `--apiKey` and `--srvID` flags like this:

```shell
mikctl --apiKey XXX --srvID YYY
```

## Testing your configuration

To test that your API key is correct and `mikctl` is reading it properly, run:

```shell
mikctl server
```

or

```shell
 mikctl --srvID YYY --apiKey XXX server
```

## Getting server info

The `mikctl server` command will list information about your server:

```shell
mikctl --srvID XXX --apiKey YYY server
ServerID: j230
Server name:
Expiration date: 2026-06-08 00:00:00
Cytrus expiration date:
Storage expiration date:
RAM size: 1024
Disk size: 10
Last log time: 2024-06-07 09:06:35
Is Pro service: nie
```

## Listing servers

The `mikctl servers` command will list basic information about your provisioned servers:

```shell
mikctl --srvID XXX --apiKey YYY  servers

Server ID: a135
Server name:
Expiration date: 2025-06-05 00:00:00
RAM size: 1024
ParamDisk: 10

Server ID: j330
Server name:
Expiration date: 2026-06-08 00:00:00
RAM size: 1024
ParamDisk: 10
```

## Listing logs

The `mikctl logs` command will list last ten (max) log messages:

```shell
mikctl --srvID j230 --apiKey XXX logs

ID: 3756
Server ID: j230
Task: sshkey
Created: 2024-06-07 09:06:58
Done: 2024-06-07 09:07:01
Output: Uploaded SSH key

ID: 3751
Server ID: j230
Task: restart
Created: 2024-06-05 09:57:54
Done: 2024-06-05 09:58:07
Output: OK

ID: 3749
Server ID: j230
Task: password
Created: 2024-06-05 09:16:50
Done: 2024-06-05 09:17:02
Output: OK

ID: 3748
Server ID: j230
Task: upgrade
Created: 2024-06-05 08:59:28
Done: 2024-06-05 09:00:04
Output: === Aktualne parametry: 768 RAM / 10 DYSK 2 / 20 Dodaje: 256MB RAM oraz 0GB dysku Po zmianie: 1024 MB / 10 GB [succes] GOTOWE!
```

## Bugs and feature requests

If you find a bug in the `mikrus` client, please [open an issue](https://github.com/qba73/mikrus/issues). Similarly, if you'd like a feature added or improved, let me know via an issue.

Not all the functionality of the [Mikrus API](https://api.mikr.us) is implemented yet.

Pull requests welcome!

