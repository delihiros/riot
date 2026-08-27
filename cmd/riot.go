package cmd

import (
	"fmt"

	"github.com/delihiros/riot/config"
	"github.com/spf13/cobra"
)

var (
	cfgFile    string
	apiKey     string
	region     string
	locale     string
	matchID    string
	puuid      string
	queue      string
	actID      string
	size       int
	startIndex int

	rootCmd = &cobra.Command{
		Use:   "riot",
		Short: "Riot is a command line interface for RiotGames API",
		Long:  "Riot is a command line interface for RiotGames API",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			key, err := loadAPIKey(cfgFile, apiKey)
			if err != nil {
				return err
			}
			apiKey = key
			return nil
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func loadAPIKey(configFile string, flagValue string) (string, error) {
	if flagValue != "" {
		return flagValue, nil
	}
	if err := config.ReadConfig(configFile); err != nil {
		return "", err
	}
	key := config.GetConfig().GetString("riot.api_key")
	if key == "" {
		return "", fmt.Errorf("riot API key is required; set RIOT_API_KEY, riot.api_key in the config file, or --api-key")
	}
	return key, nil
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.riot.yaml)")
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "riot api key override")
	rootCmd.PersistentFlags().StringVar(&region, "region", "ap", "region")

	getContentCmd.PersistentFlags().StringVar(&locale, "locale", "", "optional locale")
	matchByIDCmd.PersistentFlags().StringVar(&matchID, "match-id", "", "match ID")
	matchByIDCmd.MarkPersistentFlagRequired("match-id")
	matchListByIDCmd.PersistentFlags().StringVar(&puuid, "puuid", "", "puuid")
	matchListByIDCmd.MarkPersistentFlagRequired("puuid")
	recentMatchesCmd.PersistentFlags().StringVar(&queue, "queue", "", "queue")
	recentMatchesCmd.MarkPersistentFlagRequired("queue")
	leaderboardCmd.PersistentFlags().StringVar(&actID, "act-id", "", "act IDs can be found using the content API")
	leaderboardCmd.MarkPersistentFlagRequired("act-id")
	leaderboardCmd.PersistentFlags().IntVar(&size, "size", 200, "Defaults to 200. Valid values: 1 to 200.")
	leaderboardCmd.PersistentFlags().IntVar(&startIndex, "start-index", 0, "Defaults to 0.")

	rootCmd.AddCommand(valorantCmd)
	valorantCmd.AddCommand(contentCmd)
	contentCmd.AddCommand(getContentCmd)
	valorantCmd.AddCommand(matchCmd)
	matchCmd.AddCommand(matchByIDCmd)
	matchCmd.AddCommand(matchListByIDCmd)
	matchCmd.AddCommand(recentMatchesCmd)
	valorantCmd.AddCommand(rankedCmd)
	rankedCmd.AddCommand(leaderboardCmd)
	valorantCmd.AddCommand(statusCmd)
	statusCmd.AddCommand(getStatusCmd)
}
