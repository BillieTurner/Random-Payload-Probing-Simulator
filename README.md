# A Random Payload Probing Simulator.

According to the suggestion from ["A practical guide to defend against the GFW's latest active probing"](https://gfw.report/blog/ss_advise/en/), <em>Probe your implementations</em> is a good practice to make your tool more resistant to censorship.

This is a probing simulator that generates random payloas and analyzes the response time and response message.

## How to use
* Step 1. Collect data

for QUIC server

```
go run main.go -addr=<ip>:<port> -alpn=h3 -o=log.csv -maxlen=2000 -quic
```

for TCP connection

```
go run main.go -addr=<ip>:<port> -o=log.csv -maxlen=2000
```

* Step 2. Run analyzer
```
python3 log.csv
```