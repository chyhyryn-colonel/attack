# 💣 🇷🇺

[![asciicast](https://asciinema.org/a/yz2ef8kaT6RFFSeZTOiykShkX.svg)](https://asciinema.org/a/yz2ef8kaT6RFFSeZTOiykShkX?t=20)

## How to run?

Prebuilt binaries are in `bin/` directory in this repo. Here are links:

- [Windows](https://github.com/chyhyryn-colonel/attack/blob/main/bin/win_amd64.exe?raw=true)
- [Linux](https://github.com/chyhyryn-colonel/attack/blob/main/bin/linux_amd64?raw=true)
- [macOS](https://github.com/chyhyryn-colonel/attack/blob/main/bin/mac_amd64?raw=true)
- [macOS ARM](https://github.com/chyhyryn-colonel/attack/blob/main/bin/mac_arm64?raw=true)

Just download and run. It will automatically refresh list of URLs to attack every 1 minute.

## How to install Docker?

- [Windows](https://docs.docker.com/desktop/windows/install)
- [Linux](https://docs.docker.com/engine/install/ubuntu/)
- [Mac](https://docs.docker.com/desktop/mac/install)

## How to run with Docker?

As simple as:
```
docker run -it --rm chyhyryncolonel/attack:latest
```

> NOTE! If you see that most of the websites have a **very low success rate** or **your computer freezes** consider reducing parallelism with "-p" flag. Default value (125) might be too much for Mac OS/Windows users.

```
docker run -it --rm chyhyryncolonel/attack:latest app -p 50
```

## VPN locations

Here's a list of preferred locations to VPN when performing an attack. This is based on my analysis of resource availability.
> NOTE! **Higher is better!** So if you can get VPN in the country with highter availability **go for it**.

|Country	| Availability  |
|---------------|---------------|
| Russia	| 0.70 		|
| Moldova	| 0.66 		|
| United states | 0.64 		|
| Kazakhstan	| 0.64 		|
| Austria	| 0.64 		|
| Portugal	| 0.63 		|
| Italy		| 0.63 		|
| Germany	| 0.63 		|
| Canada	| 0.63 		|
| Netherlands	| 0.62 		|
| Lithuania	| 0.61 		|
| Turkey	| 0.61 		|
| Switzerland	| 0.61 		|
| Australia	| 0.59 		|
| Hong kong	| 0.58 		|
| France	| 0.58 		|
| Iran		| 0.56 		|
| Ukraine	| 0.52 		|


Availability of 0.7 in Russia means that 70% of the URLs are **still available** from Russian IPs.

Availability of 0.52 in Ukraine means that almost 50% of all the requests are **blocked** if they're coming from Ukrainian IP, hence use VPN!

## Slava Ukraini!
