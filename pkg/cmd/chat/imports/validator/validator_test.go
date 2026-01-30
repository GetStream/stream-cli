package validator

import (
	"errors"
	"fmt"
	"os"
	"testing"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/stretchr/testify/require"
)

func TestValidator_Validate(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     *Results

		lighterChanIDValidation bool
	}{
		{name: "Valid data", filename: "valid-data.json", want: &Results{
			Stats:  map[string]int{"channels": 3, "devices": 2, "members": 4, "messages": 3, "reactions": 3, "users": 4},
			Errors: nil,
		}},
		{name: "Valid channels with lighter channel ID validation", filename: "valid-channels-with-light.json", lighterChanIDValidation: true, want: &Results{
			Stats:  map[string]int{"channels": 2, "devices": 0, "members": 2, "messages": 1, "reactions": 0, "users": 2},
			Errors: nil,
		}},
		{name: "Invalid users", filename: "invalid-users.json", want: &Results{
			Stats: map[string]int{"channels": 0, "devices": 0, "members": 0, "messages": 0, "reactions": 0, "users": 1},
			Errors: []error{
				errors.New(`validation error: user.id required`),
				errors.New(`validation error: user.id max length exceeded (255)`),
				errors.New(`validation error: user.id invalid ("^[@\w-]*$" allowed)`),
				errors.New(`validation error: user.online is a reserved field`),
				errors.New(`duplicate user found "userA"`),
				errors.New(`reference error: user.role "admin" doesn't exist (user "userA")`),
			},
		}},
		{name: "Invalid channels", filename: "invalid-channels.json", want: &Results{
			Stats: map[string]int{"channels": 2, "devices": 0, "members": 0, "messages": 0, "reactions": 0, "users": 1},
			Errors: []error{
				errors.New(`validation error: either channel.id or channel.member_ids required`),
				errors.New(`validation error: either channel.id or channel.member_ids required`),
				errors.New(`validation error: channel.id max length exceeded (64)`),
				errors.New(`validation error: channel.type required`),
				errors.New(`validation error: channel.id "channelA@abc" invalid ("^[\w-]*$" allowed)`),
				errors.New(`validation error: channel.created_by required`),
				errors.New(`validation error: channel.cid is a reserved field`),
				errors.New(`reference error: channel.type "" doesn't exist (channel ":channelA")`),
				errors.New(`reference error: channel.created_by "" doesn't exist (channel "messaging:channelA")`),
				errors.New(`reference error: channel.type "livestream" doesn't exist (channel "livestream:[userA]")`),
				errors.New(`reference error: channel.created_by "userB" doesn't exist (channel "messaging:channelB")`),
				errors.New(`validation error: distinct channel: ["userA"] is missing members: ["userA"]. Please include all members as separate member entries`),
			},
		}},
		{name: "Invalid channel mutes", filename: "invalid-channel-mutes.json", want: &Results{
			Stats: map[string]int{"channels": 2, "devices": 0, "members": 0, "messages": 0, "reactions": 0, "users": 2},
			Errors: []error{
				errors.New(`reference error: muted channel "messaging:456" by user "user1" doesn't exist`),
			},
		}},
		{name: "Invalid user mutes", filename: "invalid-user-mutes.json", want: &Results{
			Stats: map[string]int{"channels": 0, "devices": 0, "members": 0, "messages": 0, "reactions": 0, "users": 2},
			Errors: []error{
				errors.New(`reference error: muted user "missing_user" by user "user1" doesn't exist`),
			},
		}},
		{name: "Invalid members", filename: "invalid-members.json", want: &Results{
			Stats: map[string]int{"channels": 4, "devices": 0, "members": 5, "messages": 0, "reactions": 0, "users": 3},
			Errors: []error{
				errors.New(`validation error: member.user_id required`),
				errors.New(`validation error: member.channel_type required`),
				errors.New(`validation error: either member.channel_id or member.channel_member_ids required`),
				errors.New(`reference error: user "" doesn't exist`),
				errors.New(`reference error: channel ":channelA" doesn't exist`),
				errors.New(`reference error: distinct channel with type "messaging" and members:[] doesn't exist`),
				errors.New(`reference error: user "userA" with teams map[] cannot be a member of channel messaging:channelB with team "teamB"`),
				errors.New(`reference error: user "userC" specified as channel member but not present in channel_members_ids: [userA userB]`),
				errors.New(`validation error: distinct channel: ["userA" "userB"] is missing members: ["userB"]. Please include all members as separate member entries`),
			},
		}},
		{name: "Invalid messages", filename: "invalid-messages.json", want: &Results{
			Stats: map[string]int{"channels": 2, "devices": 0, "members": 2, "messages": 1, "reactions": 0, "users": 1},
			Errors: []error{
				errors.New(`validation error: message.id max length exceeded (255)`),
				errors.New(`validation error: message.channel_type required`),
				errors.New(`validation error: either message.channel_id or message.channel_member_ids required`),
				errors.New(`validation error: message.user required`),
				errors.New(`validation error: message.type invalid ("regular", "deleted", "system" and "reply" allowed)`),
				errors.New(`validation error: message.type "deleted" while message.deleted_at is null`),
				errors.New(`validation error: message.deleted_at "2022-02-14T12:34:30Z" while message.type is "regular"`),
				errors.New(`reference error: channel ":channelA" doesn't exist`),
				errors.New(`reference error: distinct channel with type "messaging" and members:[] doesn't exist`),
				errors.New(`reference error: user "" doesn't exist (message_id messageA)`),
				errors.New(`reference error: message parent_id "parentID" doesn't exist`),
			},
		}},
		{name: "Invalid devices", filename: "invalid-devices.json", want: &Results{
			Stats: map[string]int{"channels": 0, "devices": 2, "members": 0, "messages": 0, "reactions": 0, "users": 1},
			Errors: []error{
				errors.New(`validation error: device.id max length exceeded (255)`),
				errors.New(`validation error: device.id required`),
				errors.New(`validation error: device.user_id required`),
				errors.New(`validation error: device.push_provider_type invalid, available options are: firebase,huawei,apn,xiaomi`),
				errors.New(`validation error: device.created_at required`),
				errors.New(`duplicate device id:duplicate`),
				errors.New(`reference error: device.user_id "userB" doesn't exist`),
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open("testdata/" + tt.filename)
			require.NoError(t, err)
			defer func() {
				_ = f.Close()
			}()

			var options []Options
			if tt.lighterChanIDValidation {
				options = append(options, LighterValidationChannelID())
			}
			v := New(f, []*stream.Role{{Name: "user"}}, map[string]*stream.ChannelType{"messaging": nil}, options...)

			got := v.Validate()

			require.Equal(t, tt.want.Stats, got.Stats)
			require.Equal(t, len(tt.want.Errors), len(got.Errors), got.Errors)
			for i := range tt.want.Errors {
				require.Equal(t, tt.want.Errors[i].Error(), got.Errors[i].Error(), fmt.Sprintf(`errors #%d doesn't match`, i))
			}
		})
	}
}
