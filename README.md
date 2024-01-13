# Pprof visualizer

## Overview
Pprof visualizer is a fork of (to be completed) aimed at providing
a web interface for loading and visualize pprof reports.

## Features
* [WIP] PProf reports as Tree: Easily upload pprof reports and visualize
them in a user friendly web interface
* [WIP] Compare Pprof reports: Upload two pprof reports and highlight
difference

## Getting Started

### Prerequisites
* Go (1.21.4)
* Docker for linting
* Make

### Installation and Setup
1. Clone the repository
```
git clone https://github.com/artemys/pprof-visualizer.git
```

2. Install dependencies
```
make deps
```

3. Environment setup
* Create a `.env` based on `.env.example`

4. Run the project
```
make run
```

4.1 Run on IDEA
```
Configuration > Run king > Package
Configuration > Package path : github.com/artemys/pprof-visualizer
Configuration > Program arguments: api
EnvFile > .env
```

### Cleaning up
To clean the project, use:
```
make clean
```