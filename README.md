# AssemblyAI server

This repository is a micro service which communicate to [AssemblyAI](https://www.assemblyai.com/) and get the transcripted text version of your submited audio file via the [client](https://github.com/praveenmahasena/texttospeechclient)

## Features
- works as a middle man between AssemblyAI and the client

## Getting Started

### Prerequisites
- [Go](https://go.dev/dl/) programming language
- [Key](https://www.assemblyai.com/app) API key to access AssemblyAI server on your own

## Environment Variables

Make sure you have configure your Environment Variables properly on your `.env` file. The `.env` file contains following properties
```
TOKEN=token to assemblyai api
MEDIAUPLOADLINK=https://api.assemblyai.com/v2/upload
```

### Installation
```
git clone git@github.com:praveenmahasena/texttospeechserver.git
```

### Start
```bash
make texttospeechserver && ./bin/texttospeechserver
```

