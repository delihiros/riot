package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/delihiros/riot/pkg/client"
	content "github.com/delihiros/riot/pkg/valorant/content/v1"
	match "github.com/delihiros/riot/pkg/valorant/match/v1"
	ranked "github.com/delihiros/riot/pkg/valorant/ranked/v1"
	status "github.com/delihiros/riot/pkg/valorant/status/v1"

	"github.com/spf13/cobra"
)

var (
	valorantCmd = &cobra.Command{
		Use:   "valorant",
		Short: "subcommand for valorant",
		Long:  "subcommand for valorant",
	}

	contentCmd = &cobra.Command{
		Use:   "content",
		Short: "subcommand for content",
		Long:  "subcommand for content: https://developer.riotgames.com/apis#val-content-v1",
	}
	getContentCmd = &cobra.Command{
		Use:   "get-content",
		Short: "Get content optionally filtered by locale",
		Long:  "Get content optionally filtered by locale",
		RunE: func(cmd *cobra.Command, args []string) error {
			httpClient := client.New(apiKey)
			c := content.New(httpClient)
			dto, err := c.GetContent(region, locale)
			if err != nil {
				return err
			}
			bytes, err := json.Marshal(dto)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), string(bytes))
			return err
		},
	}

	matchCmd = &cobra.Command{
		Use:   "match",
		Short: "subcommand for match",
		Long:  "subcommand for match: https://developer.riotgames.com/apis#val-match-v1",
	}
	matchByIDCmd = &cobra.Command{
		Use:   "get-match",
		Short: "Get match by id",
		RunE: func(cmd *cobra.Command, args []string) error {
			httpClient := client.New(apiKey)
			c := match.New(httpClient)
			dto, err := c.GetMatchByID(region, matchID)
			if err != nil {
				return err
			}
			bytes, err := json.Marshal(dto)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), string(bytes))
			return err
		},
	}
	matchListByIDCmd = &cobra.Command{
		Use:   "matchlist",
		Short: "Get matchlist for games played by puuid",
		RunE: func(cmd *cobra.Command, args []string) error {
			httpClient := client.New(apiKey)
			c := match.New(httpClient)
			dto, err := c.GetMatchListByPUUID(region, puuid)
			if err != nil {
				return err
			}
			bytes, err := json.Marshal(dto)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), string(bytes))
			return err
		},
	}
	recentMatchesCmd = &cobra.Command{
		Use:   "recent-matches",
		Short: "Get recent matches for a queue",
		RunE: func(cmd *cobra.Command, args []string) error {
			httpClient := client.New(apiKey)
			c := match.New(httpClient)
			dto, err := c.GetRecentMatches(region, queue)
			if err != nil {
				return err
			}
			bytes, err := json.Marshal(dto)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), string(bytes))
			return err
		},
	}

	rankedCmd = &cobra.Command{
		Use:   "ranked",
		Short: "subcommand for ranked",
		Long:  "subcommand for ranked: https://developer.riotgames.com/apis#val-ranked-v1",
	}
	leaderboardCmd = &cobra.Command{
		Use:   "leaderboard",
		Short: "Get leaderboard for the competitive queue",
		RunE: func(cmd *cobra.Command, args []string) error {
			httpClient := client.New(apiKey)
			c := ranked.New(httpClient)
			dto, err := c.GetLeaderboard(region, actID, size, startIndex)
			if err != nil {
				return err
			}
			bytes, err := json.Marshal(dto)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), string(bytes))
			return err
		},
	}

	statusCmd = &cobra.Command{
		Use:   "status",
		Short: "subcommand for status",
		Long:  "subcommand for status: https://developer.riotgames.com/apis#val-status-v1",
	}
	getStatusCmd = &cobra.Command{
		Use:   "get-status",
		Short: "Get VALORANT status for the given platform.",
		RunE: func(cmd *cobra.Command, args []string) error {
			httpClient := client.New(apiKey)
			c := status.New(httpClient)
			dto, err := c.GetStatus(region)
			if err != nil {
				return err
			}
			bytes, err := json.Marshal(dto)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), string(bytes))
			return err
		},
	}
)
