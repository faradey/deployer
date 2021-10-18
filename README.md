# deployer
Deploy your project on the server

## Description
This little program is designed to help you deploy a project to your server.
The program is implemented as an http server that listens for requests to a specific port and path. Upon receiving the request, the program reads the configuration file and executes all the commands specified in it.

### Key Features
* Running asynchronous commands.
* Running asynchronous groups of commands.
* Each command or group of commands can be run from a specific user on the system.
* Each command or group of commands can be run in a specific directory.

## Installation
1. Clone this repo
```
git clone git@github.com:faradey/deployer.git
```
2. Setup configuration and run server
    * Go to `deployer` folder
    * Rename the `deployer-config.sample` file to `deployer-config`
    * Fill in the `deployer-config` file with the required options and attributes. [Description of options and attributes](./docs/DEPLOYERCONFIG.md)
    * Start the server with the command `./deployer` as root or using sudo
    
## Usage
   * Configure an action that will make an http request for your domain and `PORT`, `PATH` specified in the `deployer-config`
   ```
   Example URL
   http://example.com:8083/deploy/123456789/abcdefg
   ```
   * When the action is triggered, the server will accept the request and execute the commands specified in the `deployer-config`
