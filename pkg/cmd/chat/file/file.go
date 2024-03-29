package file

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/cmd/chat/utils"
	"github.com/GetStream/stream-cli/pkg/config"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		uploadFileCmd(),
		uploadImageCmd(),
		deleteFileCmd(),
		deleteImageCmd(),
	}
}

func uploadFileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload-file --channel-type [channel-type] --channel-id [channel-id] --user-id [user-id] --file [file]",
		Short: "Upload a file",
		Long: heredoc.Doc(`
			Stream will not block any file types from uploading, however, different
			clients may handle different types differently or not at all.
			You can set a more restrictive list for your application if needed.
			The maximum file size is 100MB.
			Stream will allow any file extension. If you want to be more restrictive
			for an application, this is can be set via API or by logging into your dashboard.
		`),
		Example: heredoc.Doc(`
			# Uploads a file to 'redteam' channel of 'messaging' channel type
			$ stream-cli chat upload-file --channel-type messaging --channel-id redteam --user-id "user-1" --file "./snippet.txt"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			chType, _ := cmd.Flags().GetString("channel-type")
			chID, _ := cmd.Flags().GetString("channel-id")
			user, _ := cmd.Flags().GetString("user-id")
			file, _ := cmd.Flags().GetString("file")

			path, err := utils.UploadFile(c, cmd, chType, chID, user, file)
			if err != nil {
				return err
			}

			cmd.Printf("Successfully uploaded file: %s\n", path)
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("channel-type", "t", "", "[required] Channel type to interact with")
	fl.StringP("channel-id", "i", "", "[required] Channel id to interact with")
	fl.StringP("user-id", "u", "", "[required] User id")
	fl.StringP("file", "f", "", "[required] File path")
	_ = cmd.MarkFlagRequired("channel-type")
	_ = cmd.MarkFlagRequired("channel-id")
	_ = cmd.MarkFlagRequired("user-id")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func uploadImageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload-image --channel-type [channel-type] --channel-id [channel-id] --user-id [user-id] --file [file] --content-type [content-type]",
		Short: "Upload an image",
		Long: heredoc.Doc(`
			Stream supported image types are: image/bmp, image/gif, image/jpeg, image/png,
			image/webp, image/heic, image/heic-sequence, image/heif, image/heif-sequence,
			image/svg+xml.
			You can set a more restrictive list for your application if needed.
			The maximum file size is 100MB.
			Stream will allow any file extension. If you want to be more restrictive
			for an application, this is can be set via API or by logging into your dashboard.
		`),
		Example: heredoc.Doc(`
			# Uploads an image to 'redteam' channel of 'messaging' channel type
			$ stream-cli chat upload-image --channel-type messaging --channel-id redteam --user-id "user-1" --file "./picture.png" --content-type "image/png"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			chType, _ := cmd.Flags().GetString("channel-type")
			chID, _ := cmd.Flags().GetString("channel-id")
			user, _ := cmd.Flags().GetString("user-id")
			file, _ := cmd.Flags().GetString("file")

			path, err := utils.UploadImage(c, cmd, chType, chID, user, file)
			if err != nil {
				return err
			}

			cmd.Printf("Successfully uploaded image: %s\n", path)
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("channel-type", "t", "", "[required] Channel type to interact with")
	fl.StringP("channel-id", "i", "", "[required] Channel id to interact with")
	fl.StringP("user-id", "u", "", "[required] User id")
	fl.StringP("file", "f", "", "[required] Image file path")
	_ = cmd.MarkFlagRequired("channel-type")
	_ = cmd.MarkFlagRequired("channel-id")
	_ = cmd.MarkFlagRequired("user-id")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func deleteFileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-file --channel-type [channel-type] --channel-id [channel-id] --file-url [file-url]",
		Short: "Delete a file",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			chType, _ := cmd.Flags().GetString("channel-type")
			chID, _ := cmd.Flags().GetString("channel-id")
			fileUrl, _ := cmd.Flags().GetString("file-url")

			_, err = c.Channel(chType, chID).DeleteFile(cmd.Context(), fileUrl)
			if err != nil {
				return err
			}

			cmd.Println("Successfully deleted file")
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("channel-type", "t", "", "[required] Channel type to interact with")
	fl.StringP("channel-id", "i", "", "[required] Channel id to interact with")
	fl.StringP("file-url", "u", "", "[required] URL of the file to delete")
	_ = cmd.MarkFlagRequired("channel-type")
	_ = cmd.MarkFlagRequired("channel-id")
	_ = cmd.MarkFlagRequired("file-url")

	return cmd
}

func deleteImageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-image --channel-type [channel-type] --channel-id [channel-id] --image-url [image-url]",
		Short: "Delete an image",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			chType, _ := cmd.Flags().GetString("channel-type")
			chID, _ := cmd.Flags().GetString("channel-id")
			imageURL, _ := cmd.Flags().GetString("image-url")

			_, err = c.Channel(chType, chID).DeleteImage(cmd.Context(), imageURL)
			if err != nil {
				return err
			}

			cmd.Println("Successfully deleted image")
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("channel-type", "t", "", "[required] Channel type to interact with")
	fl.StringP("channel-id", "i", "", "[required] Channel id to interact with")
	fl.StringP("image-url", "u", "", "[required] URL of the image to delete")
	_ = cmd.MarkFlagRequired("channel-type")
	_ = cmd.MarkFlagRequired("channel-id")
	_ = cmd.MarkFlagRequired("image-url")

	return cmd
}
