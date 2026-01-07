# Dyslexic's Clock

This repo contains the code for a Word Clock capable of connecting to a web server and reading configs from there.

### Codes

The main branch is empty. you can find code for frontend in `server-fronted` , backend in `server-backend` and the microcontroller in `microcontroller-v2` 

## Features:
- Clock:
    - shows time in 5 minute intervals
    - capable of handling multiple alarms
    - snooze and stop logic
    - happy birthday logic
    - can ring and stop with web interfaces
    - has wifi portal to connect to any internet
    - has offline and online modes
    - customizable brightness
    - customizable coloring
    - customizable volume
- Backend and Frontend (web app):
    - user interacts with clock by these interfaces (alarm crud and every other config)
    - uses mqtt to send messages to clock
