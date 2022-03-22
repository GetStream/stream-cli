package imports

import (
	"context"
	"os"

	stream "github.com/GetStream/stream-chat-go/v5"

	"github.com/GetStream/stream-cli/pkg/cmd/chat/import/validator"
)

func ValidateFile(ctx context.Context, client *stream.Client, filename string) (*validator.Results, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rolesResp, err := client.Permissions().ListRoles(ctx)
	if err != nil {
		return nil, err
	}

	channelTypesResp, err := client.ListChannelTypes(ctx)
	if err != nil {
		return nil, err
	}

	v := validator.New(f, rolesResp.Roles, channelTypesResp.ChannelTypes)
	return v.Validate(), nil
}
