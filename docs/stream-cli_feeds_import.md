# stream-cli feeds import

Imports feeds data from a JSON file.

This command uploads the file to S3 and initiates the import process.

## Usage

```bash
stream-cli feeds import [filename] --apikey [api-key]
```

## Arguments

- `filename`: Path to the JSON file to import

## Flags

- `-k, --apikey string`: [required] API key for authentication
- `-m, --mode string`: [optional] Import mode. Can be upsert or insert (default "upsert")

## Examples

```bash
# Import feeds data with API key
$ stream-cli feeds import feeds-data.json --apikey your-api-key

# Import feeds data with custom mode
$ stream-cli feeds import feeds-data.json --apikey your-api-key --mode insert
```

## Import Modes

- `upsert`: Updates existing feeds and creates new ones (default)
- `insert`: Only creates new feeds, fails if feed already exists

## Output

The command will output the import task details including the import ID:

```
✅ Import started successfully
Import ID: dcb6e366-93ec-4e52-af6f-b0c030ad5272
```

The full import task object will also be printed in the specified output format.
