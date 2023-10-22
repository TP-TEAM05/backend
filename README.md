# 🚗 ReCo web backend 🌐 ☁️

## 📖 Prerequisites

###  GOLANG

[Download](https://go.dev/learn/)

###  GO AIR Package

Used for live reload while developing.

```bash
go install github.com/cosmtrek/air@latest
```

### Delve

Used for debugging.

```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

## First time setup

### 📦 Dependencies

To install all dependencies run the command below.

```bash
go build .
```

## 📅 Daily business

Use the command below to start the project with live reload.

```bash
air
```

## 🐞 Debugging

Debugging is done with Delve using VSCode

- Open the project in VSCode
- Start the project with `air`
- Open the `Run and Debug` tab in VSCode
- Select `Attach to Air` from the dropdown
- You are ready to debug 🐞
