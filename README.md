# Shario

Shario: self-hosted local file sharing in your browser.
Inspired by Apple's Airdrop.
Based on [Snapdrop](https://github.com/RobinLinus/snapdrop)

**Shario is built with the following awesome technologies:**

* Vanilla HTML5 / ES6 / CSS3 frontend
* [WebRTC](http://webrtc.org/) / [WebSockets](http://www.websocket.org/)
* [Go](https://go.dev/) backend
* [Progressive Web App](https://de.wikipedia.org/wiki/Progressive_Web_App)

Have any questions? Read our [FAQ](/docs/faq.md).

## Instalation

### Build server

```bash
    cd server
    go build
```

### Server Usage

```bash
    shario
    # or
    shario --addr localhost:5000 # (default port: 3000)
```
