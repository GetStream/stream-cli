# stream-cli feeds import-validate

Validates a JSON file for feeds import.

This command checks if the file is valid JSON format.

## Usage

```bash
stream-cli feeds import-validate [filename]
```

## Arguments

- `filename`: Path to the JSON file to validate

## Examples

```bash
# Validates a JSON feeds import file
$ stream-cli feeds import-validate feeds-data.json
```

## Output

The command will output a success message if the file is valid JSON:

```
✅ File 'feeds-data.json' is valid JSON
```

If the file is invalid or doesn't exist, an error message will be displayed.
