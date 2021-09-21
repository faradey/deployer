# deployer config

### Options
* `HOST` - *optional*. Leave blank or specify localhost or your domain
* `PORT` - *required*. Specify the port to which requests will be received
* `PATH` - *required*. Specify the URL path to which requests will be received
* `SHELL` - *optional*. Specify a command line interpreter. The default is `bash`
* `USER` - *optional*. Specify USERNAME and USER GROUP. For example `USER admin admin`. This option can be specified multiple times. All commands after this option will be executed from under the specified USERNAME and USER GROUP
* `CD` - *optional*. Specify the path relative to which the commands will be executed. The default is the server folder (`deployer`). Both absolute and relative paths are supported. For relative paths, the root is the `deployer` folder. This option can be specified multiple times. All subsequent commands will use this path.
* `TRY` - *optional*. Specify the number of retries for the command that returned the error. Deploy will return an error only when all attempts are unsuccessful. The default is 1.
* `RUN` - *required*. Specify the command to be executed. Multiple commands can be combined using `&&` `||` etc. This option can be specified multiple times.
* `ASYNC_RUN` - *optional*. This option runs the specified command asynchronously.
* `ASYNC_GROUP_START` - *optional*. A pointer to the beginning of a group of commands that will be executed asynchronously, but the group itself will be executed synchronously in the command queue. This option can be specified multiple times.
* `ASYNC_GROUP_END` - *condition*. Pointer to the end of a group of asynchronous commands. This option can be specified multiple times.

### Additional features

#### Changing these options requires a server restart
`HOST` `PORT` `PATH` `SHELL`