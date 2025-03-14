# BakeTyping

A terminal-based typing practice application built with Go and the Charm Bubble Tea framework.

## Overview

BakeTyping is a simple, interactive typing practice tool that runs in your terminal. It offers multiple difficulty levels and tracks your typing accuracy and speed.

## Features

- Four difficulty levels: Beginner, Medium, Advanced, and Guru
- 30-second timer for each typing challenge
- Real-time feedback with color-coded text (correct/incorrect characters)
- Performance metrics showing typing errors
- Visual timer that changes color when time is running low

## Controls

- **Up/Down Arrow Keys**: Navigate through difficulty levels
- **Enter**: Select difficulty level and start typing
- **A**: Restart after completing or timing out
- **Ctrl+C**: Quit the application

## Installation

Ensure you have Go installed on your system, then:

```bash
# Clone the repository
git clone https://github.com/yourusername/bubble-typer.git
cd bubble-typer

# Build the application
go build

# Run the application
./typer
