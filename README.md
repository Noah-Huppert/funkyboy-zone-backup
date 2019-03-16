# Mountain Backup
File backup tool.

# Table Of Contents
- [Overview](#overview)
- [Configure](#configure)

# Overview
Creates a GZip-ed tar ball and uploads it to a S3 compatible object 
storage service.

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

Configuration:

```toml
[Files.XXXXX]
# List of files to backup
Files = [ ... ]
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
