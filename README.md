Very simple webapp intended to run on a Windows machine to control local network light devices primarily
in a very specific arcade/gaming room setup.

This application has the ability to do screen capture to compute the screen color and therefore must run locally and not inside
of a Docker container.

# Development

```powershell
go run .\cmd\main.go
```

# Build/Package

```powershell
 .\build\package\Build.bat
```

# Config / Environment Variables
This application will leverage a .env file to load environment variables.  The following is an example of the .env file:

```powershell
PORT=8000

# Valid values are: [AVERAGE, SQUARED_AVERAGE, MEDIAN, MODE]
COLOR_ALGO=AVERAGE

# How often to capture the screen (should generally be greater than 50ms due to screen capture latency)
CAPTURE_INTERVAL=80ms

# Valid values are: [LIFX]
LIGHT_TYPE=LIFX

# Name of LIFX group to control
LIGHT_GROUP_NAME=ARCADE

# Adjust maximum brightness of the lights between 0 and 1. 1 is full brightness. (makes screen flashes or white screens quite bright)
MAX_BRIGHTNESS=0.45

# Adjust minimum brightness of the lights between 0 and 1. 0 is the light turned off.
MIN_BRIGHTNESS=0

# Adjust PIXEL_GRID_SIZE to increase performance or accuracy. Lower values are slower but more accurate. 1 being the most accurate.
PIXEL_GRID_SIZE=5

# Adjust SCREEN_NUMBER to target a different screen. 0 is the primary screen.
SCREEN_NUMBER=0
```