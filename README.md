# golagraphite

![build status](https://travis-ci.org/mavlyutov/golagraphite.svg)

A tool that can be used to collect Windows Performance Counters and/or SQL Metrics and send them over to the Graphite server.

## Features

* Sends Metrics to Graphite's Carbon daemon
* Can collect Windows Performance Counters
* Can collect values by using T-SQL queries against MS SQL databases
* Converts time to UTC on sending
* All configuration can be done from a simple YAML file
* Allows you to override the hostname in Windows Performance Counters before sending on to Graphite
* Executable can be installed to run as a service
* Supports Hosted Graphite (https://www.hostedgraphite.com)

## Installation

The simplest way is to download latest release from Github and run executable.

### Modifying the Configuration File

The configuration file is fairly self-explanatory, but here is a description for each of the values.

#### Graphite Configuration Section

Configuration Name | Description
--- | ---
hosts | The server name and the port number where Carbon is running. The Carbon daemon is usually running on the Graphite server.

#### Performance Counters Configuration Section

Configuration Name | Description
--- | ---
metric_prefix | The path of the Performance Counters metric you want to be sent to the server
interval | The interval to send metrics to Carbon. I recommend 5 seconds or greater. The more queries you are running the longer it takes to send them to the Graphite server.
counters | This section lists the performance counters you want the machine to send to Graphite.

You can get counters from Performance Monitor (perfmon.exe) or by using the command `typeperf -qx` in a command prompt.
I have included some basic performance counters in the configuration file. Asterisks can be used as a wildcard.
Here are some other examples:

* `\Network Interface(*)\Bytes Received/sec`
* `\Network Interface(*)\Bytes Sent/sec`
* `\PhysicalDisk(*)\Avg. Disk Write Queue Length`

#### MSSQLMetrics Configuration Section

This section allows you to configure a list of SQL servers and the queries that will be run against those servers. You can add as many queries or servers as required.

`<sql_servers>` Configuration Values | Description
--- | ---
metric_prefix | The Graphite metric name to use for this SQL server.
connection_string | The connection string to connect to SQL with using SQL Authentication. Leaving the *password* option blank will make the script use Windows Authentication against the SQL Server.

The next part of the configuration allows you to add a list of the T-SQL queries that will be run against the SQL server. You can add as many queries or servers as required.

`<queries>` Configuration Values | Description
--- | ---
interval | The interval to send metrics to Carbon. I recommend 5 seconds or greater. The more queries you are running the longer it takes to send them to the Graphite server.
tsql | The T-SQL query to run against the SQL Server.

There are a few important things to keep in mind when using this feature.

* The T-SQL query should be returned with named columns which will be used as metric suffixes. You can return as many columns as you want.
* If you provide the SQL **Username** and **Password** options, they is stored in plain text in the configuration file. If you do not provide a username and password, the windows account that the golagraphite is running under will be used against the SQL Server. This is a good way to protect the credentials.
* There is no verification that the SQL command in the configuration file is not destructive. Be sure to use a low privilege account to authenticate against SQL so that any malicious T-SQL queries don't destroy your data.
* If your T-SQL query returns multiple results, only the first one will be sent over to Graphite.

## Installing as a Service

Once you have edited the configuration file and verified everything is functioning correctly, you might want to install golagraphite as a service.

The easiest way to achieve this is using NSSM - the Non-Sucking Service Manager.

1. Download nssm from [nssm.cc](http://nssm.cc)
2. Open up an Administrative command prompt and run `nssm install golagraphite`. (You can call the service whatever you want).
3. A dialog will pop up allowing you to enter in settings for the new service. The following two tables below contains the settings to use.

![NSSM Dialog](http://i.imgur.com/6xjVkkT.png "NSSM Dialog")

4. Click *Install Service*
5. Make sure the service is started and it is set to Automatic
6. Check your Graphite server and make sure the metrics are coming in
