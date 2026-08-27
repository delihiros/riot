package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestLoadAPIKeyFromConfigFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "riot.yaml")
	if err := os.WriteFile(path, []byte("riot:\n  api_key: from-file\n"), 0600); err != nil {
		t.Fatal(err)
	}

	got, err := loadAPIKey(path, "")
	if err != nil {
		t.Fatal(err)
	}
	if got != "from-file" {
		t.Fatalf("api key = %q", got)
	}
}

func TestLoadAPIKeyFlagOverridesUnreadableConfig(t *testing.T) {
	got, err := loadAPIKey(filepath.Join(t.TempDir(), "missing.yaml"), "from-flag")
	if err != nil {
		t.Fatal(err)
	}
	if got != "from-flag" {
		t.Fatalf("api key = %q", got)
	}
}

func TestRootPreRunLoadsAPIKeyForSubcommands(t *testing.T) {
	path := filepath.Join(t.TempDir(), "riot.yaml")
	if err := os.WriteFile(path, []byte("riot:\n  api_key: from-prerun\n"), 0600); err != nil {
		t.Fatal(err)
	}

	oldConfigFile, oldAPIKey := cfgFile, apiKey
	defer func() {
		cfgFile, apiKey = oldConfigFile, oldAPIKey
	}()
	cfgFile, apiKey = path, ""
	if rootCmd.PersistentPreRunE == nil {
		t.Fatal("root command has no persistent config loader")
	}
	if err := rootCmd.PersistentPreRunE(rootCmd, nil); err != nil {
		t.Fatal(err)
	}
	if apiKey != "from-prerun" {
		t.Fatalf("api key = %q", apiKey)
	}
}

func TestValorantRecentMatchesCommandIsRegistered(t *testing.T) {
	command, _, err := rootCmd.Find([]string{"valorant", "match", "recent-matches"})
	if err != nil {
		t.Fatal(err)
	}
	if command.Name() != "recent-matches" {
		t.Fatalf("found command %q", command.Name())
	}
}

func TestContentLocaleDefaultsToUnfiltered(t *testing.T) {
	flag := getContentCmd.PersistentFlags().Lookup("locale")
	if flag == nil {
		t.Fatal("locale flag is not registered")
	}
	if flag.DefValue != "" {
		t.Fatalf("locale default = %q, want empty", flag.DefValue)
	}
}

func TestValorantCommandsReturnErrorsThroughCobra(t *testing.T) {
	oldAPIKey, oldRegion := apiKey, region
	oldLocale, oldMatchID, oldPUUID := locale, matchID, puuid
	oldQueue, oldActID := queue, actID
	defer func() {
		apiKey, region = oldAPIKey, oldRegion
		locale, matchID, puuid = oldLocale, oldMatchID, oldPUUID
		queue, actID = oldQueue, oldActID
	}()

	apiKey, region = "test-key", "invalid/region"
	matchID, puuid, queue, actID = "match", "player", "competitive", "act"

	commands := []*cobra.Command{
		getContentCmd,
		matchByIDCmd,
		matchListByIDCmd,
		recentMatchesCmd,
		leaderboardCmd,
		getStatusCmd,
	}
	for _, command := range commands {
		if command.Run != nil {
			t.Errorf("%s uses Run; want RunE", command.Name())
			continue
		}
		if command.RunE == nil {
			t.Errorf("%s has no RunE", command.Name())
			continue
		}
		if err := command.RunE(command, nil); err == nil {
			t.Errorf("%s returned nil error", command.Name())
		}
	}
}
