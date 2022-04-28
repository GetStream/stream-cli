package validator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_index_AddUser(t *testing.T) {
	i := newIndex(nil, nil)

	require.Zero(t, i.stats()["users"])
	require.NoError(t, i.addUser("john", nil))
	require.Equal(t, 1, i.stats()["users"])
	require.Error(t, i.addUser("john", nil))
	require.NoError(t, i.addUser("paul", nil))
	require.Equal(t, 2, i.stats()["users"])
}

func Test_getChannelID(t *testing.T) {
	channelID, isDistinct := getChannelID("xyz", nil)
	require.Equal(t, "xyz", channelID)
	require.False(t, isDistinct)

	channelID, isDistinct = getChannelID("", []string{"john", "paul"})
	require.True(t, strings.HasPrefix(channelID, DistinctChannelPrefix))
	require.True(t, isDistinct)
}

func Test_index_addChannel(t *testing.T) {
	i := newIndex(nil, nil)

	require.Zero(t, i.stats()["channels"])
	require.NoError(t, i.addChannel("messaging", "abc", nil, ""))
	require.Equal(t, 1, i.stats()["channels"])
	require.Error(t, i.addChannel("messaging", "abc", nil, ""))
	require.NoError(t, i.addChannel("messaging", "", []string{"john", "paul"}, ""))
	require.Equal(t, 2, i.stats()["channels"])
	require.Error(t, i.addChannel("messaging", "", []string{"john", "paul"}, ""))
}

func Test_index_addMember(t *testing.T) {
	i := newIndex(nil, nil)

	require.Zero(t, i.stats()["members"])
	require.NoError(t, i.addMember("messaging", "abc", nil, "john"))
	require.Equal(t, 1, i.stats()["members"])
	require.Error(t, i.addMember("messaging", "abc", nil, "john"))
	require.NoError(t, i.addMember("messaging", "", []string{"john", "paul"}, "paul"))
	require.Equal(t, 2, i.stats()["members"])
	require.Error(t, i.addMember("messaging", "", []string{"john", "paul"}, "paul"))
}

func Test_index_addMessage(t *testing.T) {
	i := newIndex(nil, nil)

	require.Zero(t, i.stats()["messages"])
	require.NoError(t, i.addMessage("msg1", ""))
	require.Equal(t, 1, i.stats()["messages"])
	require.Error(t, i.addMessage("msg1", ""))
	require.NoError(t, i.addMessage("msg2", ""))
	require.Equal(t, 2, i.stats()["messages"])
}

func Test_index_addReaction(t *testing.T) {
	i := newIndex(nil, nil)

	require.Zero(t, i.stats()["reactions"])
	require.NoError(t, i.addReaction("msg1", "like", "john"))
	require.Equal(t, 1, i.stats()["reactions"])
	require.Error(t, i.addReaction("msg1", "like", "john"))
	require.NoError(t, i.addReaction("msg1", "like", "paul"))
	require.Equal(t, 2, i.stats()["reactions"])
}

func Test_index_channelTypeExist(t *testing.T) {
	i := newIndex(nil, channelTypeMap{"messaging": nil})

	require.True(t, i.channelTypeExist("messaging"))
	require.False(t, i.channelTypeExist("livestream"))
}

func Test_index_isReply(t *testing.T) {
	i := newIndex(nil, nil)

	require.NoError(t, i.addMessage("msg1", ""))
	require.NoError(t, i.addMessage("msg2", "msg1"))
	require.True(t, i.isReply("msg2"))
	require.False(t, i.isReply("msg1"))
}

func Test_index_roleExist(t *testing.T) {
	i := newIndex(roleMap{"moderator": nil}, nil)

	require.True(t, i.roleExist("moderator"))
	require.False(t, i.roleExist("admin"))
}

func Test_index_sameTeam(t *testing.T) {
	i := newIndex(nil, nil)

	require.NoError(t, i.addUser("john", []string{"the beatles"}))
	require.NoError(t, i.addUser("mick", []string{"the rolling stones"}))
	require.NoError(t, i.addChannel("messaging", "liverpool", nil, "the beatles"))
	require.NoError(t, i.addChannel("messaging", "london", nil, "the rolling stones"))

	require.NoError(t, i.sameTeam("john", "messaging", "liverpool"))
	require.Error(t, i.sameTeam("john", "messaging", "london"))

	require.NoError(t, i.sameTeam("mick", "messaging", "london"))
	require.Error(t, i.sameTeam("mick", "messaging", "liverpool"))
}
