<div align="center">
<img src="ui/public/logo-120px.png" />
<h1> Divulge - Golang Key Value Pair DB Web UI (WIP) </h1>
<h3> Divulge "makes known" Key-value DB data. It's yet another Golang service with a simple UI for managing and operating multiple Key Value pair DBs written in Golang. </h3>

[![Go](https://github.com/ric-v/divulge-keyvalue-db-ui/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/ric-v/divulge-keyvalue-db-ui/actions/workflows/go.yml)
[![CodeQL](https://github.com/ric-v/divulge-keyvalue-db-ui/actions/workflows/codeql-analysis.yml/badge.svg?branch=main)](https://github.com/ric-v/divulge-keyvalue-db-ui/actions/workflows/codeql-analysis.yml)
[![Maintained](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://img.shields.io/badge/Maintained%3F-yes-green.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/ric-v/divulge-keyvalue-db-ui)](https://goreportcard.com/report/github.com/ric-v/divulge-keyvalue-db-ui)
[![CodeFactor](https://www.codefactor.io/repository/github/ric-v/divulge-keyvalue-db-ui/badge)](https://www.codefactor.io/repository/github/ric-v/divulge-keyvalue-db-ui)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ric-v_divulge-keyvalue-db-ui&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=ric-v_divulge-keyvalue-db-ui)
[![GoDoc](https://godoc.org/github.com/narqo/go-badge?status.svg)](https://pkg.go.dev/github.com/ric-v/divulge-keyvalue-db-ui/database#)


</div>

---

Simple DB CRUD operations service. Supports some golang Key-Value pair file based Databases. Upload the DB create a local copy and modify/view data and download the updates.

<!-- ![in-action-gif](https://github.com/ric-v/divulge-keyvalue-db-ui/blob/main/public/assets/screenshots/in-action.gif) -->

## Features

- [x] ~~Upload existing DB~~
- [x] ~~View Key-Value pairs~~
- [x] ~~Add new Key-Value pair~~
- [x] ~~Remove Key-Value pair~~
- [x] ~~Update Key-Value pair~~
- [x] ~~Download updated file~~
- [x] View Buckets in boltDB
- [x] Add / remove bucket
- [ ] Move/Copy Key-Value pair under a bucket to another bucket

## Usage

Download the latest [release from here](https://github.com/ric-v/divulge-keyvalue-db-ui/releases)

- ### windows

  - Unzip / Untar the release
  - Open the folder
  - Run the .exe file (in the pop up screen, click more info > Run anyway > allow firewall access if the pop-up comes)
  - The service will be available at <http://localhost:8080/>

- ### linux

  - Unzip / Untar the release
  - Open the folder in terminal
  - Run commands:

  ```bash
  chmod +x divulge-viewer-*-amd64
  ./divulge-viewer-*-amd64
  ```

  - The service will be available at <http://localhost:8080/>

## Screeshots

**Home screen:**![Home Light](ui/public/screenshots/home-light.png)

**Home screen Dark side:** ![Home Dark](ui/public/screenshots/home-dark.png)

**Upload/Create new DB:** ![Upload / Create DB](ui/public/screenshots/upload-create-db.png)

**BuntDB paginated view:** ![BuntDB Paginated](ui/public/screenshots/buntdb-paginated.png)

**BoltDB view:** ![BoltDB View](ui/public/screenshots/bolt-view.png)

**BoltDB view dark side:** ![BoltDB View Dark](ui/public/screenshots/bolt-view-dark.png)

**Manage buckets:** ![Manage Bolt Buckets](ui/public/screenshots/manage-buckets-bolt.png)

## Supported DB

- [BuntDB](https://github.com/tidwall/buntdb)
- [BoltDB](https://github.com/boltdb/bolt) (WIP)

## Technologies used

- [Golang 1.17.x](https://go.dev/)
- [ReactJS 17.x](https://reactjs.org/)
- [Material UI v5.x](https://mui.com/)
- [fasthttp](https://github.com/valyala/fasthttp) ([fasthttp mux](https://github.com/fasthttp/router))
- [BoltDB](https://github.com/boltdb/bolt)
- [BuntDB](https://github.com/tidwall/buntdb)

_Code.Share.Prosper_
