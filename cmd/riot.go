package cmd

import (
	"github.com/spf13/cobra"
)

var (
	cfgFile    string
	apiKey     string
	region     string
	locale     string
	matchID    string
	puuid      string
	actID      string
	size       int
	startIndex int

	rootCmd = &cobra.Command{
		Use:   "riot",
		Short: "Riot is a command line interface for RiotGames API",
		Long:  "Riot is a command line interface for RiotGames API",
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.riot.yaml)")
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "riot api key")
	rootCmd.MarkPersistentFlagRequired("api-key")
	rootCmd.PersistentFlags().StringVar(&region, "region", "ap", "region")

	getContentCmd.PersistentFlags().StringVar(&locale, "locale", "ja-JP", "locale")
	matchByIDCmd.PersistentFlags().StringVar(&matchID, "match-id", "", "match ID")
	matchByIDCmd.MarkPersistentFlagRequired("match-id")
	matchListByIDCmd.PersistentFlags().StringVar(&puuid, "puuid", "", "puuid")
	matchListByIDCmd.MarkPersistentFlagRequired("puuid")
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
	valorantCmd.AddCommand(rankedCmd)
	rankedCmd.AddCommand(leaderboardCmd)
	valorantCmd.AddCommand(statusCmd)
	statusCmd.AddCommand(getStatusCmd)
}
