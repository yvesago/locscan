# locscan

a simple golang cli to quickly catch URL info, files and IPs locations with [Maxmind GeoIP](https://www.maxmind.com/)

## Usage

```
$ ./locscan
    __
   / /___  ___________________  ____
  / / __ \/ ___/ ___/ ___/ __ \/ __ \
 / / /_/ / /__(__  ) /__/ /_/ / / / /
/_/\____/\___/____/\___/\__,_/_/ /_/

Version: 1.0.0-pub, build: 2023-03-25T11:43:36+0100
Usage of ./locscan:
  -b	black & white, no color
  -cd string
    	config: config dir (default "~/.scan/")
  -g	GeoIP versions
  -p	read form pipe, cat file or Ctrl+D
  -pc
    	pipe: add compact result
  -q	quiet: no banner, no emoji
  -scan string
    	Mandatory url https://somesite/dir
  -sf string
    	scan: checks in config dir yaml file (default "test.yaml")
  -si
    	scan: ip, ssl info
  -sy
    	scan: force yes
  -t int
    	Timeout (default 500)
  -v	verbose
  -vc
    	view city name
  -vd
    	view reverse dns
  -vs
    	view short: no asn
```
### URL info

```
$ ./locscan -scan https://www.google.com -si
locscan -scan https://www.google.com -si
    __
   / /___  ___________________  ____
  / / __ \/ ___/ ___/ ___/ __ \/ __ \
 / / /_/ / /__(__  ) /__/ /_/ / / / /
/_/\____/\___/____/\___/\__,_/_/ /_/

Version: 1.0.0-pub, build: 2023-03-25T11:43:36+0100
 <――――――――――――――――――――――――
URL: https://www.google.com
SSL info:
  DNS Names: [www.google.com]
  Issuer Name: CN=GTS CA 1C3,O=Google Trust Services LLC,C=US
  Created: 02/03/2023
  Expiry: 25/05/2023
Domain:
  google.com registrar: MarkMonitor Inc.
  creation: 1997-09-15T04:00:00Z
IPs:
  www.google.com. IN A 216.58.213.68
   loc ->  216.58.213.68 (US)🇺🇸 AS15169 - GOOGLE
  NetName: GOOGLE
  NetName: Google LLC (GOGL)
  www.google.com. IN A 2a00:1450:4007:806::2004
   loc ->  2a00:1450:4007:806::2004 (IE)🇮🇪 AS15169 - GOOGLE
  NetName: IE-GOOGLE-2a00-1450-4000-1
  ――――――――――――――――――――――――――――――――>

```
### IP location parser

```
$ head ip-list.txt | build/linux/locscan -p
    __
   / /___  ___________________  ____
  / / __ \/ ___/ ___/ ___/ __ \/ __ \
 / / /_/ / /__(__  ) /__/ /_/ / / / /
/_/\____/\___/____/\___/\__,_/_/ /_/

Version: 1.0.0-pub, build: 2023-03-25T11:43:36+0100
pipe a file or Ctrl+D to start process ...
―――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
88.253.25.252
  └>  88.253.25.252 (TR)🇹🇷 AS47331 - Turk Telekom
88.240.179.182
  └>  88.240.179.182 (TR)🇹🇷 AS47331 - Turk Telekom
85.105.227.62
  └>  85.105.227.62 (TR)🇹🇷 AS47331 - Turk Telekom
85.105.174.76
  └>  85.105.174.76 (TR)🇹🇷 AS47331 - Turk Telekom
85.105.145.229
  └>  85.105.145.229 (TR)🇹🇷 AS47331 - Turk Telekom
80.23.223.57
  └>  80.23.223.57 (IT)🇮🇹 AS3269 - Telecom Italia
77.220.53.12
  └>  77.220.53.12 (RU)🇷🇺 AS30881 - MTS PJSC
75.117.65.14
  └>  75.117.65.14 (US)🇺🇸 AS7029 - WINDSTREAM
74.116.129.188
  └>  74.116.129.188 (US)🇺🇸 AS53266 - DFNSYSRSWL
66.189.23.94
  └>  66.189.23.94 (US)🇺🇸 AS20115 - CHARTER-20115

―――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
```

### ``loc.sh`` wrapper

```
$ sudo tail /var/log/auth.log | ./loc.sh
Mar 22 16:56:17 localhost sshd[3530395]: Failed password for root from 46.114.94.29 port 22520 ssh2
  └>  46.114.94.29 (DE) AS12638 - E-Plus Mobilfunk GmbH
Mar 22 16:56:18 localhost sshd[3530397]: pam_unix(sshd:auth): authentication failure; logname= uid=0 euid=0 tty=ssh ruser= rhost=124.43.136.72  user=root
  └>  124.43.136.72 (LK) AS9329 - Sri Lanka Telecom Internet
Mar 22 16:56:19 localhost sshd[3530395]: Received disconnect from 46.114.94.29 port 22520:11: Bye Bye [preauth]
  └>  46.114.94.29 (DE) AS12638 - E-Plus Mobilfunk GmbH
```



## Build


```bash
$ git clone https://github.com/yvesago/locscan.git
$ cd locscan
$ go test
$ make
```

``make all`` for cross compiling

## Licence

MIT License

Copyright (c) 2023 Yves Agostini

<yves+github@yvesago.net>
