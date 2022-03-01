package cli

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v2"
)

func PrintMessage(ctx *cli.Context, message string) {
	ctx.App.Writer.Write([]byte(" > " + message + "\n"))
}

func PrintMessageFormat(ctx *cli.Context, message string, args ...interface{}) {
	PrintMessage(ctx, fmt.Sprintf(message, args...))
}

func PrintHappyMessage(ctx *cli.Context, message string) {
	PrintMessage(ctx, message+" ✅")
}

func PrintHappyMessageFormatted(ctx *cli.Context, message string, args ...interface{}) {
	PrintHappyMessage(ctx, fmt.Sprintf(message, args...))
}

func PrintSadMessage(ctx *cli.Context, message string) {
	PrintMessage(ctx, message+" ❌")
}

func PrintRawJson(ctx *cli.Context, data interface{}) {
	res, _ := json.MarshalIndent(data, "", "  ")
	ctx.App.Writer.Write(res)
}
