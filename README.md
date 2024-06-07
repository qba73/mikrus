# mikrus

`mikrus` is a Go library and command-line client for [MIKRUS VPS](https://mikr.us) provider. It allows you to interact with provisioned servers and perform various tasks, for example:

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

To install the client binary, run:

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

## Bugs and feature requests

If you find a bug in the `mikrus` client or library, please [open an issue](https://github.com/qba73/mikrus/issues). Similarly, if you'd like a feature added or improved, let me know via an issue.

Not all the functionality of the [Mikrus API](https://api.mikr.us) is implemented yet.

Pull requests welcome!
