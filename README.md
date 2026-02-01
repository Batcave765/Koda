# Koda 🎵

A full-stack application built with a **Go** backend and **Next.js** frontend for analyzing music files to detect **BPM** and **Musical Key**.

## Features

- **Format Support**: Accepts MP3, WAV, and MP4 files.
- **BPM Detection**: Leverages the `benjojo/bpm` algorithm for accurate beat detection.
- **Key Detection**: Implements a custom Krumhansl-Schmuckler algorithm using FFT and Chroma mapping.
- **Modern UI**: A clean, responsive dashboard built with Next.js, Tailwind CSS, and Roboto typography.
- **Dockerized Backend**: Seamless environment setup with FFmpeg pre-installed in the container.

## Project Structure

- `backend/`: Go server using the Gin framework and Docker.
- `frontend/`: Next.js 15+ application using Bun and Tailwind CSS.

## Prerequisites

- [Docker & Docker Compose](https://www.docker.com/products/docker-desktop/)
- [Node.js](https://nodejs.org/) & [Bun](https://bun.sh/)

## Getting Started

### 1. Start the Backend (via Docker)
The backend container handles all audio processing and FFmpeg dependencies.

```bash
cd backend
docker-compose up --build
```
The server will be available at `http://localhost:8080`.

### 2. Start the Frontend
The frontend proxies API requests to the backend container.

```bash
cd frontend
bun install
bun run dev
```
Open [http://localhost:3000](http://localhost:3000) in your browser.

## API Usage

If you wish to use the backend API directly:

- **Endpoint**: `POST /analyze`
- **Payload**: `multipart/form-data`
- **Field**: `file` (your audio file)

Example:
```bash
curl -F "file=@your-song.mp3" http://localhost:8080/analyze
```

## Troubleshooting

- **FFmpeg Errors**: Ensure the backend container is running, as it contains all necessary audio codecs.
- **CORS/Proxy Issues**: The frontend uses `next.config.ts` rewrites to proxy `/api` requests to port `8080`. Ensure both are running simultaneously.
