# stream-cli feeds import-status

Checks the status of a feeds import operation.

You can optionally watch for completion with the --watch flag.

## Usage

```bash
stream-cli feeds import-status [import-id]
```

## Arguments

- `import-id`: The ID of the import operation to check

## Flags

- `-w, --watch`: [optional] Watch import until completion
- `-o, --output-format string`: [optional] Output format. Can be json or tree (default "json")

## Examples

```bash
# Check import status
$ stream-cli feeds import-status dcb6e366-93ec-4e52-af6f-b0c030ad5272

# Watch import until completion
$ stream-cli feeds import-status dcb6e366-93ec-4e52-af6f-b0c030ad5272 --watch
```

## Output Formats

- `json`: Output in JSON format (default)
- `tree`: Output in a browsable tree format

## Watch Mode

When using the `--watch` flag, the command will continuously poll the import status every 5 seconds until the import is completed or failed.

## Output

The command will output the current status of the import task, including details such as:

- Import ID
- Status (pending, running, completed, failed)
- Progress information
- Error details (if any)
