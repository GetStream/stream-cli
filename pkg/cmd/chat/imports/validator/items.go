package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	streamchat "github.com/GetStream/stream-chat-go/v5"
)

var (
	validUserIDRe       = regexp.MustCompile(`^[@\w-]*$`)
	validChannelIDRe    = regexp.MustCompile(`^[\w-]*$`)
	validReactionTypeRe = regexp.MustCompile(`^[\w-+:.]*$`)
)

type Item interface {
	validateFields() error
	index(*index) error
	validateReferences(*index) error
}

func newItem(rawItem *rawItem) (Item, error) {
	switch rawItem.Type {
	case "user":
		return newUserItem(rawItem.Item)
	case "channel":
		return newChannelItem(rawItem.Item)
	case "member":
		return newMemberItem(rawItem.Item)
	case "message":
		return newMessageItem(rawItem.Item)
	case "reaction":
		return newReactionItem(rawItem.Item)
	case "device":
		return newDeviceItem(rawItem.Item)
	default:
		return nil, fmt.Errorf("invalid item type %q", rawItem.Type)
	}
}

type rawItem struct {
	Type string          `json:"type" validate:"required"`
	Item json.RawMessage `json:"item" validate:"required"`
}

type extraFields map[string]interface{}

var extraFieldsType = reflect.TypeOf(extraFields{})

func unmarshalItem(data []byte, v interface{}) error {
	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	// extra fields...
	m := make(map[string]interface{})
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	vType := reflect.TypeOf(v).Elem()
	vValue := reflect.ValueOf(v).Elem()
	for i := 0; i < vType.NumField(); i++ {
		name := strings.Split(vType.Field(i).Tag.Get("json"), ",")[0]
		delete(m, name)

		field := vValue.Field(i)
		if field.Type() == extraFieldsType {
			field.Set(reflect.ValueOf(m))
		}
	}
	return nil
}

var userReservedFields = []string{
	"app_pk", "online", "pk", "user_pk", "last_active",
}

func newUserItem(itemBody json.RawMessage) (*userItem, error) {
	user := userItem{
		Role: "user",
	}
	if err := unmarshalItem(itemBody, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

type userItem struct {
	ID                string           `json:"id"`
	Role              string           `json:"role"`
	Invisible         bool             `json:"invisible"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
	DeletedAt         *time.Time       `json:"deleted_at"`
	DeactivatedAt     *time.Time       `json:"deactivated_at"`
	Teams             []string         `json:"teams"`
	PushNotifications pushNotification `json:"push_notifications"`
	Custom            extraFields
}

type pushNotification struct {
	Disabled      bool       `json:"disabled"`
	DisabledUntil *time.Time `json:"disabled_until"`
}

func (u *userItem) validateFields() error {
	if u.ID == "" {
		return errors.New("user.id required")
	}

	if len(u.ID) > 255 {
		return errors.New("user.id max length exceeded (255)")
	}

	if !validUserIDRe.MatchString(u.ID) {
		return fmt.Errorf(`user.id invalid ("%s" allowed)`, validUserIDRe)
	}

	for _, field := range userReservedFields {
		if _, found := u.Custom[field]; found {
			return fmt.Errorf("user.%s is a reserved field", field)
		}
	}
	return nil
}

func (u *userItem) index(idx *index) error {
	return idx.addUser(u.ID, u.Teams)
}

func (u *userItem) validateReferences(idx *index) error {
	if !idx.roleExist(u.Role) {
		return fmt.Errorf("user.role %q doesn't exist (user %q)", u.Role, u.ID)
	}
	return nil
}

type deviceItem struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	CreatedAt        time.Time `json:"created_at"`
	Disabled         bool      `json:"disabled"`
	DisabledReason   string    `json:"disabled_reason"`
	PushProviderType string    `json:"push_provider_type"`
	PushProviderName string    `json:"push_provider_name"`
}

func newDeviceItem(itemBody json.RawMessage) (*deviceItem, error) {
	var device deviceItem
	if err := unmarshalItem(itemBody, &device); err != nil {
		return nil, err
	}
	return &device, nil
}

var pushProviders = []string{
	streamchat.PushProviderFirebase,
	streamchat.PushProviderHuawei,
	streamchat.PushProviderAPNS,
	streamchat.PushProviderXiaomi,
}

func (d *deviceItem) validateFields() error {
	if d.ID == "" {
		return errors.New("device.id required")
	}
	if len(d.ID) > 255 {
		return errors.New("device.id max length exceeded (255)")
	}

	if d.UserID == "" {
		return errors.New("device.user_id required")
	}
	if len(d.UserID) > 255 {
		return errors.New("device.user_id max length exceeded (255)")
	}

	if d.CreatedAt.IsZero() {
		return errors.New("device.created_at required")
	}

	if d.PushProviderType == "" {
		return errors.New("device.push_provider_type required")
	}
	var found bool
	for _, p := range pushProviders {
		if d.PushProviderType == p {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("device.push_provider_type invalid, available options are: %s", strings.Join(pushProviders, ","))
	}

	return nil
}

func (d *deviceItem) index(i *index) error {
	return i.addDevice(d.ID)
}

func (d *deviceItem) validateReferences(i *index) error {
	if d.UserID != "" && !i.userExist(d.UserID) {
		return fmt.Errorf("device.user_id %q doesn't exist", d.UserID)
	}
	return nil
}

var channelReservedFields = []string{
	"last_message_at", "cid", "created_by_pk", "members", "config", "app_pk", "pk",
}

func newChannelItem(itemBody json.RawMessage) (*channelItem, error) {
	var channel channelItem
	if err := unmarshalItem(itemBody, &channel); err != nil {
		return nil, err
	}
	return &channel, nil
}

type channelItem struct {
	ID        string    `json:"id"`
	MemberIDs []string  `json:"member_ids"`
	Type      string    `json:"type"`
	CreatedBy string    `json:"created_by"`
	Frozen    bool      `json:"frozen"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Team      string    `json:"team"`
	Custom    extraFields
}

func (c *channelItem) validateFields() error {
	if (c.ID == "" && len(c.MemberIDs) == 0) || (c.ID != "" && len(c.MemberIDs) > 0) {
		return errors.New("either channel.id or channel.member_ids required")
	}

	if len(c.ID) > 64 {
		return errors.New("channel.id max length exceeded (64)")
	}

	if !validChannelIDRe.MatchString(c.ID) {
		return fmt.Errorf(`channel.id invalid ("%s" allowed)`, validChannelIDRe)
	}

	if c.Type == "" {
		return errors.New("channel.type required")
	}

	if c.CreatedBy == "" {
		return errors.New("channel.created_by required")
	}

	for _, field := range channelReservedFields {
		if _, found := c.Custom[field]; found {
			return fmt.Errorf("channel.%s is a reserved field", field)
		}
	}
	return nil
}

func (c *channelItem) index(idx *index) error {
	return idx.addChannel(c.Type, c.ID, c.MemberIDs, c.Team)
}

func (c *channelItem) validateReferences(idx *index) error {
	if !idx.channelTypeExist(c.Type) {
		if c.ID == "" {
			return fmt.Errorf(`channel.type %q doesn't exist (channel "%s:%v")`, c.Type, c.Type, c.MemberIDs)
		}
		return fmt.Errorf(`channel.type %q doesn't exist (channel "%s:%s")`, c.Type, c.Type, c.ID)
	}

	if !idx.userExist(c.CreatedBy) {
		if c.ID == "" {
			return fmt.Errorf(`channel.created_by %q doesn't exist (channel "%s:%v")`, c.CreatedBy, c.Type, c.MemberIDs)
		}
		return fmt.Errorf(`channel.created_by %q doesn't exist (channel "%s:%s")`, c.CreatedBy, c.Type, c.ID)
	}
	return nil
}

func newMemberItem(itemBody json.RawMessage) (*memberItem, error) {
	member := memberItem{
		CreatedAt: time.Now(),
	}
	if err := unmarshalItem(itemBody, &member); err != nil {
		return nil, err
	}
	return &member, nil
}

type memberItem struct {
	ChannelType        string    `json:"channel_type"`
	ChannelID          string    `json:"channel_id"`
	ChannelMemberIDs   []string  `json:"channel_member_ids"`
	UserID             string    `json:"user_id"`
	IsModerator        bool      `json:"is_moderator"`
	Invited            bool      `json:"invited"`
	LastRead           time.Time `json:"last_read"`
	CreatedAt          time.Time `json:"created_at"`
	InviteAcceptedAt   time.Time `json:"invite_accepted_at"`
	InviteRejectedAt   time.Time `json:"invite_rejected_at"`
	HideChannel        bool      `json:"hide_channel"`
	HideMessagesBefore time.Time `json:"hide_messages_before"`
}

func (m *memberItem) validateFields() error {
	if m.ChannelType == "" {
		return errors.New("member.channel_type required")
	}

	if (m.ChannelID == "" && len(m.ChannelMemberIDs) == 0) || (m.ChannelID != "" && len(m.ChannelMemberIDs) > 0) {
		return errors.New("either member.channel_id or member.channel_member_ids required")
	}

	if m.UserID == "" {
		return errors.New("member.user_id required")
	}

	return nil
}

func (m *memberItem) index(idx *index) error {
	return idx.addMember(m.ChannelType, m.ChannelID, m.ChannelMemberIDs, m.UserID)
}

func (m *memberItem) validateReferences(idx *index) error {
	channelID, isDistinct := getChannelID(m.ChannelID, m.ChannelMemberIDs)
	if !idx.channelExist(m.ChannelType, channelID) {
		if isDistinct {
			return fmt.Errorf("distinct channel with type %q and members:%v doesn't exist", m.ChannelType, m.ChannelMemberIDs)
		}
		return fmt.Errorf(`channel "%s:%s" doesn't exist`, m.ChannelType, m.ChannelID)
	}
	if !idx.userExist(m.UserID) {
		return fmt.Errorf("user %q doesn't exist", m.UserID)
	}

	if err := idx.sameTeam(m.UserID, m.ChannelType, channelID); err != nil {
		return err
	}

	channelPK := getChannelPK(m.ChannelType, channelID)
	if len(m.ChannelMemberIDs) > 0 {
		if memberSet, ok := idx.channelPKsToMemberSet[channelPK]; ok {
			if _, found := memberSet[m.UserID]; found {
				delete(memberSet, m.UserID)
			} else {
				return fmt.Errorf("user %q specified as channel member but not present in channel_members_ids: %v",
					m.UserID,
					m.ChannelMemberIDs,
				)
			}
		} else {
			return fmt.Errorf("user %q specified as a distinct channel member, but no distinct channel declared with members: %v",
				m.UserID,
				m.ChannelMemberIDs,
			)
		}
	}

	return nil
}

var messageReservedFields = []string{
	"tmp_id",
	"app_pk",
	"latest_reactions",
	"own_reactions",
	"reactions_count",
	"reply_count",
	"command",
	"pk",
	"quoted_message",
	"pinned_by",
}

const (
	TypeMessageRegular string = "regular"
	TypeMessageDeleted string = "deleted"
	TypeMessageSystem  string = "system"
	TypeMessageReply   string = "reply"
)

func newMessageItem(itemBody json.RawMessage) (*messageItem, error) {
	member := messageItem{
		Type: TypeMessageRegular,
	}
	if err := unmarshalItem(itemBody, &member); err != nil {
		return nil, err
	}
	return &member, nil
}

type messageItem struct {
	ID               string            `json:"id"`
	ParentID         string            `json:"parent_id"`
	ChannelType      string            `json:"channel_type"`
	ChannelID        string            `json:"channel_id"`
	ChannelMemberIDs []string          `json:"channel_member_ids"`
	Text             string            `json:"text"`
	HTML             string            `json:"html"`
	Attachments      []*attachmentItem `json:"attachments"`
	User             string            `json:"user"`
	Type             string            `json:"type"`
	ShowInChannel    bool              `json:"show_in_channel"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	DeletedAt        time.Time         `json:"deleted_at"`
	MentionedUserIDs []string          `json:"mentioned_users_ids"`
	QuotedMessageID  string            `json:"quoted_message_id"`
	PinnedAt         *time.Time        `json:"pinned_at"`
	PinnedByID       string            `json:"pinned_by_id"`
	PinExpires       *time.Time        `json:"pin_expires"`
	Custom           extraFields
}

func (m *messageItem) validateFields() error {
	if len(m.ID) > 255 {
		return errors.New("message.id max length exceeded (255)")
	}

	if m.ChannelType == "" {
		return errors.New("message.channel_type required")
	}

	if (m.ChannelID == "" && len(m.ChannelMemberIDs) == 0) || (m.ChannelID != "" && len(m.ChannelMemberIDs) > 0) {
		return errors.New("either message.channel_id or message.channel_member_ids required")
	}

	if m.User == "" {
		return errors.New("message.user required")
	}

	switch m.Type {
	case TypeMessageRegular, TypeMessageDeleted, TypeMessageSystem, TypeMessageReply:
	default:
		return fmt.Errorf("message.type invalid (%q, %q, %q and %q allowed)", TypeMessageRegular, TypeMessageDeleted, TypeMessageSystem, TypeMessageReply)
	}

	if m.Type == TypeMessageDeleted && m.DeletedAt.IsZero() {
		return fmt.Errorf("message.type %q while message.deleted_at is null", TypeMessageDeleted)
	}
	if !m.DeletedAt.IsZero() && m.Type != TypeMessageDeleted {
		return fmt.Errorf("message.deleted_at %q while message.type is %q", m.DeletedAt.Format(time.RFC3339), m.Type)
	}

	if m.PinnedByID != "" && m.PinnedAt == nil {
		return fmt.Errorf("message.pinned_by_id %q while pinned_at is null", m.PinnedByID)
	}
	if m.PinnedByID == "" && m.PinnedAt != nil {
		return fmt.Errorf("message.pinned_at %q while pinned_by_id is null", m.PinnedAt)
	}
	if m.PinExpires != nil && (m.PinnedAt == nil || m.PinnedByID == "") {
		return fmt.Errorf("message.pin_expires %q while pinned_at/pinned_by_id is null", m.PinExpires)
	}

	for i := range m.Attachments {
		if err := m.Attachments[i].validateFields(); err != nil {
			return err
		}
	}

	for _, field := range messageReservedFields {
		if _, found := m.Custom[field]; found {
			return fmt.Errorf("message.%s is a reserved field", field)
		}
	}
	return nil
}

func (m *messageItem) index(idx *index) error {
	return idx.addMessage(m.ID, m.ParentID)
}

func (m *messageItem) validateReferences(idx *index) error {
	if m.ParentID != "" {
		if idx.isReply(m.ParentID) {
			return errors.New("only one level thread is supported")
		}
	}
	channelID, isDistinct := getChannelID(m.ChannelID, m.ChannelMemberIDs)
	if !idx.channelExist(m.ChannelType, channelID) {
		if isDistinct {
			return fmt.Errorf(
				"distinct channel with type %q and members:%v doesn't exist",
				m.ChannelType,
				m.ChannelMemberIDs,
			)
		}
		return fmt.Errorf(`channel "%s:%s" doesn't exist`, m.ChannelType, m.ChannelID)
	}

	if !idx.userExist(m.User) {
		return fmt.Errorf("user %q doesn't exist (message_id %s)", m.User, m.ID)
	}

	for _, userID := range m.MentionedUserIDs {
		if !idx.userExist(userID) {
			return fmt.Errorf("mentioned user %q doesn't exist (message_id %q)", userID, m.ID)
		}
	}

	if maxMessageLength := idx.maxMessageLength(m.ChannelType); utf8.RuneCountInString(m.Text) > maxMessageLength {
		return fmt.Errorf("message.text max length exceeded (%d)", maxMessageLength)
	}

	return nil
}

var attachmentReservedFields = []string{
	"actions", "fields",
}

type attachmentItem struct {
	Type             string `json:"type"`
	Fallback         string `json:"fallback"`
	Color            string `json:"color"`
	Pretext          string `json:"pretext"`
	AuthorName       string `json:"author_name"`
	AuthorLink       string `json:"author_link"`
	AuthorIcon       string `json:"author_icon"`
	Title            string `json:"title"`
	TitleLink        string `json:"title_link"`
	Text             string `json:"text"`
	ImageURL         string `json:"image_url"`
	ThumbURL         string `json:"thumb_url"`
	Footer           string `json:"footer"`
	FooterIcon       string `json:"footer_icon"`
	AssetURL         string `json:"asset_url"`
	OGScrapeURL      string `json:"og_scrape_url"`
	OriginalWidth    int    `json:"original_width"`
	OriginalHeight   int    `json:"original_height"`
	Size             int    `json:"size"`
	MigrateResources bool   `json:"migrate_resources"`
	Custom           extraFields
}

func (a *attachmentItem) validateFields() error {
	if a.Type == "" {
		return errors.New("attachment.type required")
	}

	for _, field := range attachmentReservedFields {
		if _, found := a.Custom[field]; found {
			return fmt.Errorf("attachment.%s is a reserved field", field)
		}
	}
	return nil
}

var reactionReservedFields = []string{
	"app_pk", "pk", "user", "user_pk",
}

func newReactionItem(itemBody json.RawMessage) (*reactionItem, error) {
	var reaction reactionItem
	if err := unmarshalItem(itemBody, &reaction); err != nil {
		return nil, err
	}
	return &reaction, nil
}

type reactionItem struct {
	MessageID string    `json:"message_id" validate:"required"`
	Type      string    `json:"type" validate:"required,max=30,reactionType"`
	UserID    string    `json:"user_id" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	Custom    extraFields
}

func (r *reactionItem) validateFields() error {
	if r.MessageID == "" {
		return errors.New("reaction.message_id required")
	}

	if r.Type == "" {
		return errors.New("reaction.type required")
	}

	if len(r.Type) > 30 {
		return errors.New("reaction.type max length exceeded (30)")
	}

	if !validReactionTypeRe.MatchString(r.Type) {
		return fmt.Errorf(`reaction.type invalid ("%s" allowed)`, validReactionTypeRe)
	}

	if r.UserID == "" {
		return errors.New("reaction.user_id required")
	}

	if r.CreatedAt.IsZero() {
		return errors.New("reaction.created_at required")
	}

	for _, field := range reactionReservedFields {
		if _, found := r.Custom[field]; found {
			return fmt.Errorf("reaction.%s is a reserved field", field)
		}
	}
	return nil
}

func (r *reactionItem) index(idx *index) error {
	return idx.addReaction(r.MessageID, r.Type, r.UserID)
}

func (r *reactionItem) validateReferences(idx *index) error {
	if !idx.messageExist(r.MessageID) {
		return fmt.Errorf("message_id %q doesn't exist", r.MessageID)
	}
	if !idx.userExist(r.UserID) {
		return fmt.Errorf("user %q doesn't exist", r.UserID)
	}

	return nil
}
