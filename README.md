# Discord OperaGX nitro generator

## Requirements
- Go 1.21+ (https://go.dev/dl/)
- The Tor Expert Bundle (https://www.torproject.org/download/tor/)
- A connection to the internet

## How to use
You need to download the Tor Expert Bundle from https://www.torproject.org/download/tor/
Unzip it into C:\Tor (That's my preference, you can put it anywhere you want)
You need to add the tor.exe inside the tor folder to your PATH, you can do that by following this tutorial: https://www.architectryan.com/2018/03/17/add-to-the-path-on-windows-10/
Then you run the tor.exe and wait for it to connect to the Tor network
and now you can run the program

## How to avoid getting banned
Use a rotating Tor configuration

Put this in a file named `torrc` on Windows at `%APPDATA%\tor\ `
```
CircuitBuildTimeout 10
LearnCircuitBuildTimeout 0
MaxCircuitDirtiness 30
```
This will assure a new circuit every 30 seconds

If you put it under 30 seconds it may not be able to find a new circuit and will delay the rotation process thus prolonging the time on a single IP 