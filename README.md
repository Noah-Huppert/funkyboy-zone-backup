# Mountain Backup
File backup tool.

# Table Of Contents
- [Overview](#overview)
- [Configure](#configure)

# Overview
Creates a GZip-ed tar ball and uploads it to a S3 compatible object 
storage service.

**Why "Mountain Backup"?**  
A reference to 
["Steel Mountain"](https://mrrobot.fandom.com/wiki/Steel_Mountain) in the show
Mr. Robot.  

# Configure
The tool's behavior is specified in a TOML configuration file.  

## Upload Configuration
The `Upload` section of the file defines where backups will be stored.  

Configuration:

```toml
[Upload]
# Storage API host
Host = "..."

# Storage API Key ID
KeyID = "..."

# Storage API Secret Access Key
SecretAccessKey = "..."

# (Optional) Backup name format without file extension, can use strftime symbols
# Defaults to value below
Format = "backup-%Y-%m-%d-%H:%M:%S"
```

## Metrics Configuration
The tool can push metrics to Prometheus about the backup process. To push metrics 
[Prometheus Push Gateway](https://github.com/prometheus/pushgateway) must be accessible to mountain backup.

Configuration:

```toml
[Metrics]
# Host which Prometheus Push Gateway can be accessed
PushGatewayHost = "localhost:9091"

# Value of `host` label in metrics
Host = "foobar"
```

## Backup Configuration
The mountain backup tools provides different modules to handle unique 
backup scenarios. These modules are configured by creating sub sections in the 
configuration file under the modules name's.  

The names of these sub-sections do not matter. The only constraint is that they 
should be unique in that section.  

For example to configure a module named `ExampleModule` one could create a 
configuration section named `ExampleModule.Foo` or `ExampleModule.Bar`.

### Files
The `Files` module backs up normal files.  

All configuration parameters can include shell globs.

Configuration:

```toml
[Files.XXXXX]
# List of files to backup
Files = [ "..." ]

# Files / directories to exclude from backup
Exclude = [ "..." ]
```

### Prometheus
The `Prometheus` module makes a snapshot of a Prometheus database via the 
admin API and backs it up.  

The snapshot files (`${DataDirectory}/data/snapshots/xxxx`) will be backed up 
as if they were located in the main data directory (`${DataDirectory}/data`).

This way Prometheus will use the data from the backed up snapshot. Instead of 
simply placing the snapshot files in the snapshot directory but starting with 
an empty database.

Configuration:

```toml
[Prometheus.XXXXX]
# Admin API host
AdminAPIHost = "localhost:9090"

# Directory in which Prometheus data is stored
DataDirectory = "/var/lib/prometheus"
```
