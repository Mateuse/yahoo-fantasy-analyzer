package services

import (
	"fmt"
	"strconv"

	"github.com/mateuse/yahoo-fantasy-analyzer/internal/models"
	"github.com/mateuse/yahoo-fantasy-analyzer/internal/utils"
)

func MapToLeague(response map[string]interface{}) (*models.League, error) {
	// Extract the "league" map from the response
	leagueData, ok := response["league"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response: missing league data")
	}

	// Create a new League object
	league := &models.League{}

	// Map values to the League struct
	if val, ok := leagueData["league_id"].(string); ok {
		league.LeagueID = val
	}
	if val, ok := leagueData["league_key"].(string); ok {
		league.LeagueKey = val
	}
	if val, ok := leagueData["name"].(string); ok {
		league.Name = val
	}
	if val, ok := leagueData["url"].(string); ok {
		league.URL = val
	}
	if val, ok := leagueData["logo_url"].(string); ok {
		league.LogoURL = val
	}
	if val, ok := leagueData["draft_status"].(string); ok {
		league.DraftStatus = val
	}
	if val, ok := leagueData["num_teams"].(string); ok {
		league.NumTeams, _ = strconv.Atoi(val)
	}
	if val, ok := leagueData["scoring_type"].(string); ok {
		league.ScoringType = val
	}
	if val, ok := leagueData["league_type"].(string); ok {
		league.LeagueType = val
	}
	if val, ok := leagueData["felo_tier"].(string); ok {
		league.FeloTier = val
	}
	if val, ok := leagueData["current_week"].(string); ok {
		league.CurrentWeek, _ = strconv.Atoi(val)
	}
	if val, ok := leagueData["start_week"].(string); ok {
		league.StartWeek, _ = strconv.Atoi(val)
	}
	if val, ok := leagueData["start_date"].(string); ok {
		league.StartDate = utils.ParseDate(val)
	}
	if val, ok := leagueData["end_week"].(string); ok {
		league.EndWeek, _ = strconv.Atoi(val)
	}
	if val, ok := leagueData["end_date"].(string); ok {
		league.EndDate = utils.ParseDate(val)
	}
	if val, ok := leagueData["game_code"].(string); ok {
		league.GameCode = val
	}
	if val, ok := leagueData["season"].(string); ok {
		league.Season = val
	}
	if val, ok := leagueData["weekly_deadline"].(string); ok {
		league.WeeklyDeadline = val
	}
	if val, ok := leagueData["allow_add_to_dl_extra_pos"].(string); ok {
		league.AllowAddToDLExtraPos = val == "1"
	}
	if val, ok := leagueData["is_pro_league"].(string); ok {
		league.IsProLeague = val == "1"
	}
	if val, ok := leagueData["is_cash_league"].(string); ok {
		league.IsCashLeague = val == "1"
	}
	if val, ok := leagueData["is_plus_league"].(string); ok {
		league.IsPlusLeague = val == "1"
	}
	if val, ok := leagueData["league_update_timestamp"].(string); ok {
		league.LeagueUpdateTimestamp, _ = strconv.ParseInt(val, 10, 64)
	}

	statNames := make(map[string]string)
	// Map roster positions and stat modifiers
	if settings, ok := leagueData["settings"].(map[string]interface{}); ok {
		// Map roster positions
		if rosterPositions, ok := settings["roster_positions"].(map[string]interface{}); ok {
			if positions, ok := rosterPositions["roster_position"].([]interface{}); ok {
				for _, positionData := range positions {
					if pos, ok := positionData.(map[string]interface{}); ok {
						rosterPosition := models.RosterPosition{}
						if val, ok := pos["position"].(string); ok {
							rosterPosition.Position = val
						}
						if val, ok := pos["position_type"].(string); ok {
							rosterPosition.PositionType = val
						}
						if val, ok := pos["count"].(string); ok {
							rosterPosition.Count, _ = strconv.Atoi(val)
						}
						if val, ok := pos["is_starting_position"].(string); ok {
							rosterPosition.IsStartingPosition = val == "1"
						}
						league.RosterPositions = append(league.RosterPositions, rosterPosition)
					}
				}
			}
		}

		// Extract stat categories to populate statNames
		if statCategories, ok := settings["stat_categories"].(map[string]interface{}); ok {
			if stats, ok := statCategories["stats"].(map[string]interface{}); ok {
				if statList, ok := stats["stat"].([]interface{}); ok {
					for _, statData := range statList {
						if stat, ok := statData.(map[string]interface{}); ok {
							if statID, ok := stat["stat_id"].(string); ok {
								if statName, ok := stat["name"].(string); ok {
									statNames[statID] = statName
								}
							}
						}
					}
				}
			}
		}

		// Map stat modifiers and populate StatName
		if statModifiers, ok := settings["stat_modifiers"].(map[string]interface{}); ok {
			if stats, ok := statModifiers["stats"].(map[string]interface{}); ok {
				if statList, ok := stats["stat"].([]interface{}); ok {
					for _, statData := range statList {
						if stat, ok := statData.(map[string]interface{}); ok {
							statModifier := models.StatModifier{}
							if val, ok := stat["stat_id"].(string); ok {
								statModifier.StatID = val
								// Populate StatName using the statNames map
								if name, exists := statNames[val]; exists {
									statModifier.StatName = name
								}
							}
							if val, ok := stat["value"].(string); ok {
								statModifier.Value, _ = strconv.ParseFloat(val, 64)
							}
							league.StatModifiers = append(league.StatModifiers, statModifier)
						}
					}
				}
			}
		}
	}

	return league, nil
}

func MapToTeamWeek(response map[string]interface{}) (*models.Team, error) {
	teamData, ok := response["team"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response: missing team data")
	}

	team := &models.Team{}

	// Map values to the Team struct
	if val, ok := teamData["team_key"].(string); ok {
		team.TeamKey = val
	}
	if val, ok := teamData["team_id"].(string); ok {
		team.TeamID = val
	}
	if val, ok := teamData["name"].(string); ok {
		team.Name = val
	}
	if val, ok := teamData["url"].(string); ok {
		team.URL = val
	}

	// Handle nested team logos
	if logos, ok := teamData["team_logos"].([]interface{}); ok && len(logos) > 0 {
		if logoData, ok := logos[0].(map[string]interface{}); ok {
			if url, ok := logoData["url"].(string); ok {
				team.LogoURL = url
			}
		}
	}

	if val, ok := teamData["waiver_priority"].(string); ok {
		team.WaiverPriority, _ = strconv.Atoi(val)
	}
	if val, ok := teamData["number_of_moves"].(string); ok {
		team.NumberOfMoves, _ = strconv.Atoi(val)
	}
	if val, ok := teamData["number_of_trades"].(string); ok {
		team.NumberOfTrades, _ = strconv.Atoi(val)
	}
	if val, ok := teamData["league_scoring_type"].(string); ok {
		team.LeagueScoringType = val
	}
	if val, ok := teamData["draft_position"].(string); ok {
		team.DraftPosition, _ = strconv.Atoi(val)
	}

	// Map projected points
	if projectedPoints, ok := teamData["team_projected_points"].(map[string]interface{}); ok {
		if total, ok := projectedPoints["total"].(string); ok {
			team.ProjectedPoints = total
		}
	}

	// Map live projected points
	if liveProjectedPoints, ok := teamData["team_live_projected_points"].(map[string]interface{}); ok {
		if total, ok := liveProjectedPoints["total"].(string); ok {
			team.FinalPoints = total
		}
	}

	// Handle nested remaining games data
	if remainingGames, ok := teamData["team_remaining_games"].(map[string]interface{}); ok {
		if total, ok := remainingGames["total"].(map[string]interface{}); ok {
			if remaining, ok := total["remaining_games"].(string); ok {
				team.RemainingGames, _ = strconv.Atoi(remaining)
			}
			if completed, ok := total["completed_games"].(string); ok {
				team.CompletedGames, _ = strconv.Atoi(completed)
			}
		}
	}

	return team, nil
}
