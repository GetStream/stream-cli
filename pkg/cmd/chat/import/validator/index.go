package validator

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"

	streamchat "github.com/GetStream/stream-chat-go/v5"
)

var (
	hash    = sha512.New()
	hashSep = []byte{','}
)

type Bits256 [sha512.Size256]byte

func hashValues(values ...string) Bits256 {
	var sum256 Bits256
	hash.Reset()
	for i := range values {
		_, _ = hash.Write([]byte(values[i]))
		_, _ = hash.Write(hashSep)
	}
	sum := hash.Sum(nil)
	copy(sum256[:], sum[:sha512.Size256])
	return sum256
}

type Set map[string]struct{}

func NewSet(ls ...string) Set {
	m := make(Set, len(ls))
	for _, s := range ls {
		m[s] = struct{}{}
	}
	return m
}

type (
	roleMap        map[string]*streamchat.Role
	channelTypeMap map[string]*streamchat.ChannelType
)

type index struct {
	// users
	userRoles      roleMap
	userPKs        map[Bits256]struct{}
	userPKsToTeams map[Bits256]Set

	// channels
	channelTypes          channelTypeMap
	channelPKs            map[Bits256]struct{}
	channelPKsToTeam      map[Bits256]string
	channelPKsToMemberSet map[Bits256]Set
	channelPKsToMembers   map[Bits256][]string

	// members
	memberPKs map[Bits256]struct{}

	// messages
	messagePKs             map[Bits256]struct{}
	messagePKsWithReaction map[Bits256]struct{}
	messagePKsReplies      map[Bits256]struct{}

	// reactions
	reactionPKs map[Bits256]struct{}
}

func newIndex(roles map[string]*streamchat.Role, channelTypes map[string]*streamchat.ChannelType) *index {
	return &index{
		userRoles:              roles,
		userPKs:                make(map[Bits256]struct{}),
		userPKsToTeams:         make(map[Bits256]Set),
		channelTypes:           channelTypes,
		channelPKs:             make(map[Bits256]struct{}),
		channelPKsToTeam:       make(map[Bits256]string),
		channelPKsToMemberSet:  make(map[Bits256]Set),
		channelPKsToMembers:    make(map[Bits256][]string),
		memberPKs:              make(map[Bits256]struct{}),
		messagePKs:             make(map[Bits256]struct{}),
		messagePKsWithReaction: make(map[Bits256]struct{}),
		messagePKsReplies:      make(map[Bits256]struct{}),
		reactionPKs:            make(map[Bits256]struct{}),
	}
}

func (i *index) stats() map[string]int {
	return map[string]int{
		"users":     len(i.userPKs),
		"channels":  len(i.channelPKs),
		"members":   len(i.memberPKs),
		"messages":  len(i.messagePKs),
		"reactions": len(i.reactionPKs),
	}
}

func (i *index) validateChannelMembers() error {
	errs := new(multiError)
	for channelID := range i.channelPKsToMemberSet {
		if len(i.channelPKsToMemberSet[channelID]) != 0 {
			members := make([]string, 0, len(i.channelPKsToMemberSet[channelID]))
			for member := range i.channelPKsToMemberSet[channelID] {
				members = append(members, member)
			}
			errs.add(fmt.Errorf("distinct channel: %q is missing members: %q. Please include all members as separate member entries",
				i.channelPKsToMembers[channelID], members))
		}
	}
	if errs.hasErrors() {
		return errs
	}
	return nil
}

func (i *index) roleExist(role string) bool {
	_, ok := i.userRoles[role]
	return ok
}

func getUserPK(userID string) Bits256 {
	return hashValues(userID)
}

func (i *index) userExist(userID string) bool {
	pk := getUserPK(userID)
	_, ok := i.userPKs[pk]
	return ok
}

func (i *index) addUser(userID string, teams []string) error {
	if i.userExist(userID) {
		return fmt.Errorf("duplicate user found %q", userID)
	}

	pk := getUserPK(userID)
	i.userPKs[pk] = struct{}{}
	if len(teams) > 0 {
		i.userPKsToTeams[pk] = NewSet(teams...)
	}
	return nil
}

const (
	DistinctChannelPrefix   = "!members-"
	DefaultMaxMessageLength = 20_000
)

func (i *index) channelTypeExist(channelType string) bool {
	_, ok := i.channelTypes[channelType]
	return ok
}

func (i *index) maxMessageLength(channelType string) int {
	if ct := i.channelTypes[channelType]; ct != nil {
		return ct.MaxMessageLength
	}
	return DefaultMaxMessageLength
}

func getChannelPK(channelType, channelID string) Bits256 {
	return hashValues(channelType, channelID)
}

func getChannelID(channelID string, memberIDs []string) (string, bool) {
	if channelID != "" {
		return channelID, false
	}
	sort.Strings(memberIDs)
	userString := strings.Join(memberIDs, ",")

	// Base64 encoded Sha512 takes up 43 chars, prefix adds 9, so total is 52 which fits in the column length of 60
	hasher := sha512.New512_256()
	_, _ = hasher.Write([]byte(userString))
	return DistinctChannelPrefix + base64.RawURLEncoding.EncodeToString(hasher.Sum(nil)), true
}

func (i *index) channelExist(channelType, channelID string) bool {
	pk := getChannelPK(channelType, channelID)
	_, ok := i.channelPKs[pk]
	return ok
}

func (i *index) addChannel(channelType, channelID string, memberIDs []string, team string) error {
	channelID, isDistinct := getChannelID(channelID, memberIDs)
	if i.channelExist(channelType, channelID) {
		if isDistinct {
			return fmt.Errorf("duplicate channel '%s:%v'", channelType, memberIDs)
		}
		return fmt.Errorf("duplicate channel '%s:%s'", channelType, channelID)
	}

	pk := getChannelPK(channelType, channelID)
	i.channelPKs[pk] = struct{}{}
	if team != "" {
		i.channelPKsToTeam[pk] = team
	}

	if len(memberIDs) > 0 {
		i.channelPKsToMemberSet[pk] = NewSet(memberIDs...)
		i.channelPKsToMembers[pk] = memberIDs
	}
	return nil
}

func (i *index) sameTeam(userID, channelType, channelID string) error {
	usersTeams, userHasTeams := i.userPKsToTeams[getUserPK(userID)]
	channelTeam, channelHasTeam := i.channelPKsToTeam[getChannelPK(channelType, channelID)]

	// no teams, they match
	if !userHasTeams && !channelHasTeam {
		return nil
	}
	if _, ok := usersTeams[channelTeam]; !ok {
		return fmt.Errorf("user %q with teams %v cannot be a member of channel %s:%s with team %q", userID, usersTeams, channelType, channelID, channelTeam)
	}
	return nil
}

func getMemberPK(channelType, channelID, userID string) Bits256 {
	return hashValues(channelType, channelID, userID)
}

func (i *index) memberExist(channelType, channelID, userID string) bool {
	pk := getMemberPK(channelType, channelID, userID)
	_, ok := i.memberPKs[pk]
	return ok
}

func (i *index) addMember(channelType, channelID string, channelMemberIDs []string, userID string) error {
	channelID, isDistinct := getChannelID(channelID, channelMemberIDs)
	if i.memberExist(channelType, channelID, userID) {
		if isDistinct {
			return fmt.Errorf(
				"duplicate member %q (channel with type %q and members:%v)",
				userID,
				channelType,
				channelMemberIDs,
			)
		}
		return fmt.Errorf(`duplicate member %q (channel "%s:%s")`, userID, channelType, channelID)

	}

	pk := getMemberPK(channelType, channelID, userID)
	i.memberPKs[pk] = struct{}{}
	return nil
}

func (i *index) isReply(msgID string) bool {
	_, ok := i.messagePKsReplies[getMessagePK(msgID)]
	return ok
}

func getMessagePK(messageID string) Bits256 {
	return hashValues(messageID)
}

func (i *index) messageExist(messageID string) bool {
	pk := getMessagePK(messageID)
	_, ok := i.messagePKs[pk]
	return ok
}

func (i *index) addMessage(messageID, parentID string) error {
	if i.messageExist(messageID) {
		return fmt.Errorf("duplicate message %q", messageID)
	}

	pk := getMessagePK(messageID)
	i.messagePKs[pk] = struct{}{}

	if parentID != "" {
		i.messagePKsReplies[pk] = struct{}{}
	}

	return nil
}

func getReactionPK(messageID, reactionType, userID string) Bits256 {
	return hashValues(messageID, reactionType, userID)
}

func (i *index) reactionExist(messageID, reactionType, userID string) bool {
	pk := getReactionPK(messageID, reactionType, userID)
	_, ok := i.reactionPKs[pk]
	return ok
}

func (i *index) addReaction(messageID, reactionType, userID string) error {
	if i.reactionExist(messageID, reactionType, userID) {
		return fmt.Errorf("duplicate reaction message:%s:type:%s:user:%s", messageID, reactionType, userID)
	}

	reactionPK := getReactionPK(messageID, reactionType, userID)
	i.reactionPKs[reactionPK] = struct{}{}

	messagePK := getMessagePK(messageID)
	i.messagePKsWithReaction[messagePK] = struct{}{}

	return nil
}
