package database

// prepareUserStatements readies SQL queries for later use
import (
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx"
)

//FantasyTeam contains the account IDs and meta information for a fantasy team
type FantasyTeam struct {
	ID       int64
	Owner    int64
	Name     string
	Region   string
	Position int
	Top      int64
	Jungle   int64
	Mid      int64
	Bottom   int64
	Support  int64
}

//Queries that are prepared for easy calling
const (
	QueryGetTeams       = "getTeams"
	QueryGetRandomTeams = "getRandomTeams"
	QueryInsertTeam     = "insertTeam"
	QueryUpdateTeam     = "updateTeam"
	QueryDeleteTeam     = "deleteTeam"
)

//GetTeams for the given user id from the database if they exist
func GetTeams(userID int64) (*[]FantasyTeam, error) {
	//TODO: Update this if we expect more than 1 team per player in the future
	teams := make([]FantasyTeam, 0, 1)
	rows, _ := DBConnectionPool.Query(QueryGetTeams, userID)

	for rows.Next() {
		team := FantasyTeam{}
		rows.Scan(
			&(team.ID),
			&(team.Owner),
			&(team.Name),
			&(team.Region),
			&(team.Position),
			&(team.Top),
			&(team.Jungle),
			&(team.Mid),
			&(team.Bottom),
			&(team.Support),
		)

		teams = append(teams, team)
	}

	return &teams, rows.Err()
}

//GetRandomTeams tries to retrieve numTeams non-bench fantasy teams from the DB
func GetRandomTeams(numTeams int) (*[]FantasyTeam, error) {
	//TODO: Update this if we expect more than 1 team per player in the future
	teams := make([]FantasyTeam, numTeams, numTeams)
	rows, _ := DBConnectionPool.Query(QueryGetRandomTeams, numTeams)

	//This loop should never exceed numTeams due to the DB query's row limit
	var i int
	for i = 0; rows.Next(); i++ {
		rows.Scan(
			&(teams[i].ID),
			&(teams[i].Owner),
			&(teams[i].Name),
			&(teams[i].Region),
			&(teams[i].Position),
			&(teams[i].Top),
			&(teams[i].Jungle),
			&(teams[i].Mid),
			&(teams[i].Bottom),
			&(teams[i].Support),
		)
	}

	dbError := rows.Err()
	if dbError != nil {
		return &teams, dbError
	}

	if i < numTeams {
		return &teams, errors.New(fmt.Sprintf("Only %d teams found in database out of requested %d.", i, numTeams))
	}

	return &teams, nil
}

//AddTeam to database, errors if team already exists in that slot
func AddTeam(owner int64, teamName string, region string, position int,
	top int64, jungle int64, mid int64, bottom int64, support int64) error {
	_, err := DBConnectionPool.Exec(
		QueryInsertTeam,
		owner,
		teamName,
		strings.ToLower(region),
		position,
		top,
		jungle,
		mid,
		bottom,
		support,
	)

	return err
}

//UpdateTeam that already exists in the database
func UpdateTeam(teamID int64, owner int64, teamName string, region string, position int,
	top int64, jungle int64, mid int64, bottom int64, support int64) error {
	//Get the cache IDs
	_, err := DBConnectionPool.Exec(
		QueryUpdateTeam,
		teamID,
		owner,
		teamName,
		strings.ToLower(region),
		position,
		top,
		jungle,
		mid,
		bottom,
		support,
	)

	return err
}

//DeleteTeam that already exists in the database
func DeleteTeam(teamID int64) error {
	_, err := DBConnectionPool.Exec(
		QueryDeleteTeam,
		teamID,
	)

	return err
}

func prepareTeamStatements(conn *pgx.Conn) error {
	_, err := conn.Prepare(QueryGetTeams, `
		SELECT id, owner, team_name, region, position, top, jungle, mid, bottom, support
		FROM fantasy_friends.fantasy_team
		WHERE owner=$1 AND position != 0
		ORDER BY position
  `)
	if err != nil {
		return err
	}

	//O(nlog(n)) so could be improved, does not take teams in position 0 (reserved for bench)
	_, err = conn.Prepare(QueryGetRandomTeams, `
			SELECT id, owner, team_name, region, position, top, jungle, mid, bottom, support
			FROM fantasy_friends.fantasy_team
			WHERE position != 0
			ORDER BY RANDOM()
			LIMIT $1
	  `)
	if err != nil {
		return err
	}

	_, err = conn.Prepare(QueryInsertTeam, `
		INSERT INTO fantasy_friends.fantasy_team (owner, team_name, region, position, top, jungle, mid, bottom, support)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`)
	if err != nil {
		return err
	}

	_, err = conn.Prepare(QueryUpdateTeam, `
		UPDATE fantasy_friends.fantasy_team
		SET (owner, team_name, region, position, top, jungle, mid, bottom, support) = ($2, $3, $4, $5, $6, $7, $8, $9, $10)
		WHERE id = $1
	`)

	_, err = conn.Prepare(QueryDeleteTeam, `
		DELETE FROM fantasy_friends.fantasy_team
		WHERE id = $1
	`)

	return err
}
