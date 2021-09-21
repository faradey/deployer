# deployer
Deploy your project on the server

## Installation
1. Clone this repo
```
git clone https://github.com/faradey/deployer.git
```
2. Setup configuration and run server
    * Go to deployer folder
    * Rename the deployer-config.sample file to deployer-config
    * Fill in the deployer-config file with the required options and attributes. [Description of options and attributes](./docs/DEPLOYERCONFIG.md)
    * Start the server with the command ./deployer
    
3. Usage
   * Configure an action that will make an http request for your domain and PORT, the PATH specified in the deployer-config
   ```
   Example URL
   http://example.com:8083/deploy/123456789/abcdefg
   ```
   * When the action is triggered, the server will accept the request and execute the commands specified in the deployer-config