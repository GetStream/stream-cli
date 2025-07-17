package feeds

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestNewCmds(t *testing.T) {
	cmds := NewCmds()
	if len(cmds) != 3 {
		t.Errorf("Expected 3 commands, got %d", len(cmds))
	}

	expectedCommands := []string{"import-validate", "import", "import-status"}
	for i, expected := range expectedCommands {
		if cmds[i].Name() != expected {
			t.Errorf("Expected command %d to be %s, got %s", i, expected, cmds[i].Name())
		}
	}
}

func TestNewRootCmd(t *testing.T) {
	cmd := NewRootCmd()
	if cmd.Use != "feeds" {
		t.Errorf("Expected root command to be 'feeds', got %s", cmd.Use)
	}

	if len(cmd.Commands()) != 3 {
		t.Errorf("Expected 3 subcommands, got %d", len(cmd.Commands()))
	}
}

func TestValidateFeedsFile(t *testing.T) {
	cmd := &cobra.Command{}

	// Change to project root directory for test files
	projectRoot, err := findProjectRoot()
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}

	// Test with nonexistent file
	err = validateFeedsFile(cmd, "nonexistent.json")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}

	// Test with valid complete file
	err = validateFeedsFile(cmd, filepath.Join(projectRoot, "test/feeds-sample.json"))
	if err != nil {
		t.Errorf("Expected no error for valid file, got: %v", err)
	}

	// Test with users only
	err = validateFeedsFile(cmd, filepath.Join(projectRoot, "test/users-only.json"))
	if err != nil {
		t.Errorf("Expected no error for users-only file, got: %v", err)
	}

	// Test with feeds only
	err = validateFeedsFile(cmd, filepath.Join(projectRoot, "test/feeds-only.json"))
	if err != nil {
		t.Errorf("Expected no error for feeds-only file, got: %v", err)
	}

	// Test with activities only
	err = validateFeedsFile(cmd, filepath.Join(projectRoot, "test/activities-only.json"))
	if err != nil {
		t.Errorf("Expected no error for activities-only file, got: %v", err)
	}

	// Test with reactions only
	err = validateFeedsFile(cmd, filepath.Join(projectRoot, "test/reactions-only.json"))
	if err != nil {
		t.Errorf("Expected no error for reactions-only file, got: %v", err)
	}

	// Test with invalid data
	err = validateFeedsFile(cmd, filepath.Join(projectRoot, "test/invalid-feeds.json"))
	if err == nil {
		t.Error("Expected error for invalid file")
	}
}

// findProjectRoot finds the project root directory by looking for go.mod
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find go.mod in any parent directory")
		}
		dir = parent
	}
}

func TestValidateUsers(t *testing.T) {
	// Test valid users
	users := []User{
		{ID: "user:1", Role: "user"},
		{ID: "user:2", Role: "admin"},
	}
	err := validateUsers(users)
	if err != nil {
		t.Errorf("Expected no error for valid users, got: %v", err)
	}

	// Test invalid user (empty ID)
	invalidUsers := []User{
		{ID: "", Role: "user"},
	}
	err = validateUsers(invalidUsers)
	if err == nil {
		t.Error("Expected error for user with empty ID")
	}
}

func TestValidateFeeds(t *testing.T) {
	// Test valid feeds
	feeds := []Feed{
		{ID: "feed:1", GroupID: "group:1", FID: "fid:1", Name: "user:1"},
		{ID: "feed:2", GroupID: "group:2", FID: "fid:2", Name: "user:2"},
	}
	err := validateFeeds(feeds)
	if err != nil {
		t.Errorf("Expected no error for valid feeds, got: %v", err)
	}

	// Test invalid feed (empty ID)
	invalidFeeds := []Feed{
		{ID: "", GroupID: "group:1", FID: "fid:1", Name: "user:1"},
	}
	err = validateFeeds(invalidFeeds)
	if err == nil {
		t.Error("Expected error for feed with empty ID")
	}

	// Test invalid feed (empty group_id)
	invalidFeeds2 := []Feed{
		{ID: "feed:1", GroupID: "", FID: "fid:1", Name: "user:1"},
	}
	err = validateFeeds(invalidFeeds2)
	if err == nil {
		t.Error("Expected error for feed with empty group_id")
	}

	// Test invalid feed (empty fid)
	invalidFeeds3 := []Feed{
		{ID: "feed:1", GroupID: "group:1", FID: "", Name: "user:1"},
	}
	err = validateFeeds(invalidFeeds3)
	if err == nil {
		t.Error("Expected error for feed with empty fid")
	}

	// Test invalid feed (empty name)
	invalidFeeds4 := []Feed{
		{ID: "feed:1", GroupID: "group:1", FID: "fid:1", Name: ""},
	}
	err = validateFeeds(invalidFeeds4)
	if err == nil {
		t.Error("Expected error for feed with empty name")
	}
}

func TestValidateActivities(t *testing.T) {
	// Test valid activities
	activities := []Activity{
		{ID: "activity:1", Type: "post", Feeds: []string{"feed:1"}},
		{ID: "activity:2", Type: "comment", Feeds: []string{"feed:2"}},
	}
	err := validateActivities(activities)
	if err != nil {
		t.Errorf("Expected no error for valid activities, got: %v", err)
	}

	// Test invalid activity (empty ID)
	invalidActivities := []Activity{
		{ID: "", Type: "post", Feeds: []string{"feed:1"}},
	}
	err = validateActivities(invalidActivities)
	if err == nil {
		t.Error("Expected error for activity with empty ID")
	}

	// Test invalid activity (empty type)
	invalidActivities2 := []Activity{
		{ID: "activity:1", Type: "", Feeds: []string{"feed:1"}},
	}
	err = validateActivities(invalidActivities2)
	if err == nil {
		t.Error("Expected error for activity with empty type")
	}

	// Test invalid activity (empty feeds)
	invalidActivities3 := []Activity{
		{ID: "activity:1", Type: "post", Feeds: []string{}},
	}
	err = validateActivities(invalidActivities3)
	if err == nil {
		t.Error("Expected error for activity with empty feeds")
	}
}

func TestValidateMembers(t *testing.T) {
	// Test valid members
	members := []Member{
		{UserID: "user:1", Feed: "feed:1"},
		{UserID: "user:2", Feed: "feed:2", Role: "member"},
	}
	err := validateMembers(members)
	if err != nil {
		t.Errorf("Expected no error for valid members, got: %v", err)
	}

	// Test invalid member (empty user_id)
	invalidMembers := []Member{
		{UserID: "", Feed: "feed:1"},
	}
	err = validateMembers(invalidMembers)
	if err == nil {
		t.Error("Expected error for member with empty user_id")
	}

	// Test invalid member (empty feed)
	invalidMembers2 := []Member{
		{UserID: "user:1", Feed: ""},
	}
	err = validateMembers(invalidMembers2)
	if err == nil {
		t.Error("Expected error for member with empty feed")
	}
}

func TestValidateReactions(t *testing.T) {
	// Test valid reactions
	reactions := []Reaction{
		{ActivityID: "activity:1", Type: "like", UserID: "user:1"},
		{ActivityID: "activity:2", CommentID: "comment:1", Type: "love", UserID: "user:2"},
	}
	err := validateReactions(reactions)
	if err != nil {
		t.Errorf("Expected no error for valid reactions, got: %v", err)
	}

	// Test invalid reaction (empty activity_id)
	invalidReactions := []Reaction{
		{ActivityID: "", Type: "like", UserID: "user:1"},
	}
	err = validateReactions(invalidReactions)
	if err == nil {
		t.Error("Expected error for reaction with empty activity_id")
	}

	// Test invalid reaction (empty type)
	invalidReactions2 := []Reaction{
		{ActivityID: "activity:1", Type: "", UserID: "user:1"},
	}
	err = validateReactions(invalidReactions2)
	if err == nil {
		t.Error("Expected error for reaction with empty type")
	}

	// Test invalid reaction (empty user_id)
	invalidReactions3 := []Reaction{
		{ActivityID: "activity:1", Type: "like", UserID: ""},
	}
	err = validateReactions(invalidReactions3)
	if err == nil {
		t.Error("Expected error for reaction with empty user_id")
	}
}

func TestImportDataStructure(t *testing.T) {
	// Test that ImportData can be marshaled and unmarshaled
	data := ImportData{
		Users: []User{
			{ID: "user:1", Role: "user"},
		},
		Feeds: []Feed{
			{ID: "feed:1", GroupID: "group:1", FID: "fid:1", Name: "user:1"},
		},
		Activities: []Activity{
			{ID: "activity:1", Type: "post", Feeds: []string{"feed:1"}},
		},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Errorf("Failed to marshal ImportData: %v", err)
	}

	// Unmarshal back
	var unmarshaledData ImportData
	err = json.Unmarshal(jsonData, &unmarshaledData)
	if err != nil {
		t.Errorf("Failed to unmarshal ImportData: %v", err)
	}

	// Verify data is preserved
	if len(unmarshaledData.Users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(unmarshaledData.Users))
	}
	if len(unmarshaledData.Feeds) != 1 {
		t.Errorf("Expected 1 feed, got %d", len(unmarshaledData.Feeds))
	}
	if len(unmarshaledData.Activities) != 1 {
		t.Errorf("Expected 1 activity, got %d", len(unmarshaledData.Activities))
	}
}
