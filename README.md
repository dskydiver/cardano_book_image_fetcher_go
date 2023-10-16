# Cardano Book Image Fetcher (Golang Version)

This application fetches high-resolution cover images of books based on a provided Cardano policy ID and saves them to a specified directory.

## Description

The command-line application accepts two arguments: a `Cardano policy ID` and a `path to an output directory`.

The actions performed by the application are as follows:

1. Verifies if the policy ID is a valid Book.io policy ID.
2. Downloads high-resolution cover images of up to 10 books related to the policy ID.
3. Saves these images in the specified output directory.

The application ensures idempotency. If interrupted, it resumes from where it left off, avoiding the re-download of already fetched images.

## Prerequisites

- Golang version 1.20 or later.
- An active internet connection.

## Blockfrost Configuration

Set `BLOCKFROST_PROJECT_ID` in `.env` file.

```bash
BLOCKFROST_PROJECT_ID=<your blockfrost project id>
```

## Usage

To run this command line application, clone the repo first.

```bash
git clone https://github.com/dskydiver/cardano_book_image_fetcher_go.git
cd cardano_book_image_fetcher_go
```

1. Build the application

```bash
go build
```

2. Run executable

```bash
go run main.go fetch --policy-id <policy_id> --output-dir <output_dir>
```

- example

```bash
go run main.go fetch --output-dir ./output --policy-id e7514e65f977ee4b84a8e62e7d97ea2e5c11682dfe1444d8a14e74db
```

## Development

You are welcome to contribute providing new features, fixing bugs, suggesting improvements. Please create a new issue before submitting a Pull Request.
