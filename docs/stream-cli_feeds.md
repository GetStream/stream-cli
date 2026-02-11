# stream-cli feeds

Allows you to interact with your Feeds applications

## Commands

### stream-cli feeds import-validate

Validates a JSON file for feeds import.

This command checks if the file is valid JSON format.

**Usage:**

```bash
stream-cli feeds import-validate [filename]
```

**Examples:**

```bash
# Validates a JSON feeds import file
$ stream-cli feeds import-validate feeds-data.json
```

### stream-cli feeds import

Imports feeds data from a JSON file.

This command uploads the file to S3 and initiates the import process.

**Usage:**

```bash
stream-cli feeds import [filename] --apikey [api-key]
```

**Flags:**

- `-k, --apikey string`: [required] API key for authentication
- `-m, --mode string`: [optional] Import mode. Can be upsert or insert (default "upsert")

**Examples:**

```bash
# Import feeds data with API key
$ stream-cli feeds import feeds-data.json --apikey your-api-key

# Import feeds data with custom mode
$ stream-cli feeds import feeds-data.json --apikey your-api-key --mode insert
```

### stream-cli feeds import-status

Checks the status of a feeds import operation.

You can optionally watch for completion with the --watch flag.

**Usage:**

```bash
stream-cli feeds import-status [import-id]
```

**Flags:**

- `-w, --watch`: [optional] Watch import until completion
- `-o, --output-format string`: [optional] Output format. Can be json or tree (default "json")

**Examples:**

```bash
# Check import status
$ stream-cli feeds import-status dcb6e366-93ec-4e52-af6f-b0c030ad5272

# Watch import until completion
$ stream-cli feeds import-status dcb6e366-93ec-4e52-af6f-b0c030ad5272 --watch
```
