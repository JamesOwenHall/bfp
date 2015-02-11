# Core

The core of the system is a server that listens for connections from the libraries and tracks the requests.

## Configuration file

You configure the core using a JSON configuration file.  By default, core looks for a file in the working directory named "config.json", but you can specify the path to the configuration file with the `-c` command line option.  There is an example configuration file included in this repository.  The following configuration options are available.

- `"listen type"` expects a string of either `"unix"` or `"tcp"`.  This specifies which type of socket to create.  Unix sockets generally perform better but are not available on Windows systems.
- `"listen address"` expects a string that is the address on which the server will listen.  In the case of a Unix listen type, this is the path of the socket (e.g. `"/tmp/bfp.sock"`).  In the case of a TCP listen type, you specify the (optional) address and port (e.g. `":4567"`).
- `"directions"` expects an array of objects that represent the different directions to track.  Each direction should have all of the following fields
	- `"name"` is the string of the name of the direction (e.g. `"password"`).
	- `"type"` is the string that describes the type of data.  The available options for type are `"string"` and `"int32"`.
	- `"window size"` is the positive number of the time (in seconds) that we're tracking values.
	- `"max hits"` is the positive number that limits how many hits we allow within the observation window.
	- `"clean up time"` is the positive integer that sets how many seconds should pass between clean up runs.
