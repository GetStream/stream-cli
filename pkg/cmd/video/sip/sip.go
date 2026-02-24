package sip

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		listSIPInboundRulesCmd(),
		createSIPInboundRuleCmd(),
		deleteSIPInboundRuleCmd(),
		updateSIPInboundRuleCmd(),
		listSIPTrunksCmd(),
		createSIPTrunkCmd(),
		deleteSIPTrunkCmd(),
		updateSIPTrunkCmd(),
		resolveSipInboundCmd(),
	}
}

func getHTTPClient(cmd *cobra.Command) (*utils.HTTPClient, error) {
	appCfg, err := config.GetConfig(cmd).GetAppConfig(cmd)
	if err != nil {
		return nil, err
	}
	return utils.NewHTTPClient(appCfg.AccessKey, appCfg.AccessSecretKey, appCfg.ChatURL), nil
}

func doJSON(cmd *cobra.Command, method, path string, body interface{}) error {
	h, err := getHTTPClient(cmd)
	if err != nil {
		return err
	}
	resp, err := h.DoRequest(cmd.Context(), method, path, body)
	if err != nil {
		return err
	}
	var result interface{}
	_ = json.Unmarshal(resp, &result)
	return utils.PrintObject(cmd, result)
}

func listSIPInboundRulesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-sip-inbound-rules",
		Short: "List SIP inbound routing rules",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doJSON(cmd, "GET", "video/sip/inbound_routing_rules", nil)
		},
	}
	cmd.Flags().StringP("output-format", "o", "json", "[optional] Output format")
	return cmd
}

func createSIPInboundRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-sip-inbound-rule --properties [json]",
		Short: "Create a SIP inbound routing rule",
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", "video/sip/inbound_routing_rules", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Rule properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func deleteSIPInboundRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-sip-inbound-rule [id]",
		Short: "Delete a SIP inbound routing rule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}
			_, err = h.DoRequest(cmd.Context(), "DELETE", "video/sip/inbound_routing_rules/"+args[0], nil)
			if err != nil {
				return err
			}
			cmd.Printf("Successfully deleted SIP inbound rule [%s]\n", args[0])
			return nil
		},
	}
	return cmd
}

func updateSIPInboundRuleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-sip-inbound-rule --id [id] --properties [json]",
		Short: "Update a SIP inbound routing rule",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "PUT", "video/sip/inbound_routing_rules/"+id, body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("id", "i", "", "[required] Rule ID")
	fl.StringP("properties", "p", "", "[required] Rule properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func listSIPTrunksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-sip-trunks",
		Short: "List SIP trunks",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doJSON(cmd, "GET", "video/sip/inbound_trunks", nil)
		},
	}
	cmd.Flags().StringP("output-format", "o", "json", "[optional] Output format")
	return cmd
}

func createSIPTrunkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-sip-trunk --properties [json]",
		Short: "Create a SIP trunk",
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", "video/sip/inbound_trunks", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Trunk properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func deleteSIPTrunkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-sip-trunk [id]",
		Short: "Delete a SIP trunk",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}
			_, err = h.DoRequest(cmd.Context(), "DELETE", "video/sip/inbound_trunks/"+args[0], nil)
			if err != nil {
				return err
			}
			cmd.Printf("Successfully deleted SIP trunk [%s]\n", args[0])
			return nil
		},
	}
	return cmd
}

func updateSIPTrunkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-sip-trunk --id [id] --properties [json]",
		Short: "Update a SIP trunk",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "PUT", "video/sip/inbound_trunks/"+id, body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("id", "i", "", "[required] Trunk ID")
	fl.StringP("properties", "p", "", "[required] Trunk properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func resolveSipInboundCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resolve-sip-inbound --properties [json]",
		Short: "Resolve SIP inbound routing",
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", "video/sip/resolve", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Resolve properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}
