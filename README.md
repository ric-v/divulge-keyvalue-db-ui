<div align="center">
<img src="ui/public/logo-120px.png" />
<h1> Divulge - Golang Key Value Pair DB Web UI (WIP) </h1>
<h3> Divulge "makes known" Key-value DB data. It's yet another Golang service with a simple UI for managing and operating multiple Key Value pair DBs written in Golang. </h3>

[![Go](https://github.com/ric-v/divulge-keyvalue-db-ui/actions/workflows/go.yml/badge.svg)](https://github.com/ric-v/divulge-keyvalue-db-ui/actions/workflows/go.yml)
[![CodeQL](https://github.com/ric-v/divulge-keyvalue-db-ui/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/ric-v/divulge-keyvalue-db-ui/actions/workflows/codeql-analysis.yml)

</div>

---

Simple DB CRUD operations service. Supports Key-Value pair file based Databases. Upload the DB create a local copy and modify/view data and download the updates.

<!-- ![in-action-gif](https://github.com/ric-v/divulge-keyvalue-db-ui/blob/main/public/assets/screenshots/in-action.gif) -->

## Features

- Upload existing DB
- View Buckets in boltDB
- View Key-Value pairs under a boltDB bucket
- Add new bucket
- Add new Key-Value pair under a bucket
- Remove a bucket
- Move/Copy Key-Value pair under a bucket to another bucket (WIP)

## Usage

Download the latest [release from here](https://github.com/ric-v/divulge-keyvalue-db-ui/releases)

### windows

- Unzip / Untar the release
- Open the folder
- Run the .exe file (in the pop up screen, click more info > Run anyway > allow firewall access if the pop-up comes)
- The service will be available at http://localhost:8080/

### linux

- Unzip / Untar the release
- Open the folder in terminal
- Run commands:

```bash
$ chmod +x divulge-viewer-*-amd64
$ ./divulge-viewer-*-amd64
```

- The service will be available at http://localhost:8080/

## Supported DB

- [BoltDB](https://github.com/boltdb/bolt)
- [BuntDB](https://github.com/tidwall/buntdb) (WIP)

## Technologies used

- Golang 1.17.x
- ReactJS 17.x
- Material UI v5.x
- [fasthttp](https://github.com/valyala/fasthttp) ([fasthttp mux](https://github.com/fasthttp/router))
- [BoltDB](https://github.com/boltdb/bolt)

_Code.Share.Prosper_
