package v5

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/delihiros/riot/pkg/client"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestClientGetMatchIDsByPUUID(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		assertRequest(t, r, "/lol/match/v5/matches/by-puuid/player%2F%20one/ids")
		if r.URL.RawQuery != "" {
			t.Fatalf("query = %q, want empty", r.URL.RawQuery)
		}
		return response(http.StatusOK, "200 OK", `["ASIA_1","ASIA_2"]`), nil
	})}

	got, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetMatchIDsByPUUID("asia", "player/ one", nil)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"ASIA_1", "ASIA_2"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("match IDs = %#v, want %#v", got, want)
	}
}

func TestClientGetMatchIDsByPUUIDEncodesOptions(t *testing.T) {
	startTime := int64(1710000000)
	endTime := int64(1710003600)
	queue := 420
	matchType := "ranked & special/one"
	start := 3
	count := 100
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		const want = "count=100&endTime=1710003600&queue=420&start=3&startTime=1710000000&type=ranked+%26+special%2Fone"
		if r.URL.RawQuery != want {
			t.Fatalf("query = %q, want %q", r.URL.RawQuery, want)
		}
		return response(http.StatusOK, "200 OK", `[]`), nil
	})}

	options := &MatchIDsOptions{
		StartTime: &startTime,
		EndTime:   &endTime,
		Queue:     &queue,
		Type:      &matchType,
		Start:     &start,
		Count:     &count,
	}
	if _, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetMatchIDsByPUUID("asia", "player", options); err != nil {
		t.Fatal(err)
	}
}

func TestClientGetMatchIDsByPUUIDIncludesZeroValues(t *testing.T) {
	zero64 := int64(0)
	zero := 0
	empty := ""
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		const want = "count=0&endTime=0&queue=0&start=0&startTime=0&type="
		if r.URL.RawQuery != want {
			t.Fatalf("query = %q, want %q", r.URL.RawQuery, want)
		}
		return response(http.StatusOK, "200 OK", `[]`), nil
	})}

	options := &MatchIDsOptions{
		StartTime: &zero64,
		EndTime:   &zero64,
		Queue:     &zero,
		Type:      &empty,
		Start:     &zero,
		Count:     &zero,
	}
	if _, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetMatchIDsByPUUID("asia", "player", options); err != nil {
		t.Fatal(err)
	}
}

func TestClientGetReplay(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		assertRequest(t, r, "/lol/match/v5/matches/by-puuid/player%2F%20one/replays")
		if r.URL.RawQuery != "" {
			t.Fatalf("query = %q, want empty", r.URL.RawQuery)
		}
		return response(http.StatusOK, "200 OK", `{"total":2,"matchFileURLs":["https://example.test/ASIA_1.rofl","https://example.test/ASIA_2.rofl"]}`), nil
	})}

	got, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetReplay("asia", "player/ one")
	if err != nil {
		t.Fatal(err)
	}
	want := &ReplayDTO{
		Total:         2,
		MatchFileURLs: []string{"https://example.test/ASIA_1.rofl", "https://example.test/ASIA_2.rofl"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("replay = %#v, want %#v", got, want)
	}
}

func TestClientGetMatchIDsByPUUIDReturnsNilOnPartialDecodeError(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return response(http.StatusOK, "200 OK", `["ASIA_1",2]`), nil
	})}

	got, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetMatchIDsByPUUID("asia", "player", nil)
	if err == nil {
		t.Fatal("expected partial JSON decode error")
	}
	if got != nil {
		t.Fatalf("match IDs = %#v, want nil on decode error", got)
	}
}

func TestClientGetReplayReturnsNilOnMalformedJSON(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return response(http.StatusOK, "200 OK", `{`), nil
	})}

	got, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetReplay("asia", "player")
	if err == nil {
		t.Fatal("expected malformed JSON error")
	}
	if got != nil {
		t.Fatalf("replay = %#v, want nil on decode error", got)
	}
}

func TestClientGetReplayReturnsSharedHTTPError(t *testing.T) {
	const body = `{"status":"forbidden"}`
	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return response(http.StatusForbidden, "403 Forbidden", body), nil
	})}

	got, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetReplay("asia", "player")
	if got != nil {
		t.Fatalf("replay = %#v, want nil on HTTP error", got)
	}
	var responseErr *client.HTTPError
	if !errors.As(err, &responseErr) {
		t.Fatalf("error = %v, want *client.HTTPError", err)
	}
	if responseErr.StatusCode != http.StatusForbidden {
		t.Fatalf("status code = %d, want %d", responseErr.StatusCode, http.StatusForbidden)
	}
	if string(responseErr.Body) != body {
		t.Fatalf("body = %q, want %q", responseErr.Body, body)
	}
}

func TestClientGetMatch(t *testing.T) {
	fixture := completeMatchFixture(t)
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		assertRequest(t, r, "/lol/match/v5/matches/ASIA_1%2F%20special")
		if r.URL.RawQuery != "" {
			t.Fatalf("query = %q, want empty", r.URL.RawQuery)
		}
		return response(http.StatusOK, "200 OK", fixture), nil
	})}

	got, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetMatch("asia", "ASIA_1/ special")
	if err != nil {
		t.Fatal(err)
	}
	assertSameJSON(t, got, fixture)
}

func TestMatchResponseWireSchemas(t *testing.T) {
	assertWireSchema(t, reflect.TypeFor[MatchDto](), matchSchema)
	assertWireSchema(t, reflect.TypeFor[MetadataDto](), metadataSchema)
	assertWireSchema(t, reflect.TypeFor[InfoDto](), infoSchema)
	assertWireSchema(t, reflect.TypeFor[TeamDto](), teamSchema)
	assertWireSchema(t, reflect.TypeFor[BanDto](), banSchema)
	assertWireSchema(t, reflect.TypeFor[ObjectivesDto](), objectivesSchema)
	assertWireSchema(t, reflect.TypeFor[ObjectiveDto](), objectiveSchema)
	assertWireSchema(t, reflect.TypeFor[PerksDto](), perksSchema)
	assertWireSchema(t, reflect.TypeFor[PerkStatsDto](), perkStatsSchema)
	assertWireSchema(t, reflect.TypeFor[PerkStyleDto](), perkStyleSchema)
	assertWireSchema(t, reflect.TypeFor[PerkStyleSelectionDto](), perkStyleSelectionSchema)
	assertWireSchema(t, reflect.TypeFor[MissionsDto](), missionsSchema)
	assertWireSchema(t, reflect.TypeFor[ParticipantDto](), participantSchema)
	assertWireSchema(t, reflect.TypeFor[ChallengesDto](), challengesSchema)
}

func TestClientGetMatchReturnsNilOnPartialDecodeError(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return response(http.StatusOK, "200 OK", `{"metadata":{"dataVersion":"2"},"info":{"gameId":"invalid"}}`), nil
	})}

	got, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetMatch("asia", "ASIA_1")
	if err == nil {
		t.Fatal("expected partial JSON decode error")
	}
	if got != nil {
		t.Fatalf("match = %#v, want nil on decode error", got)
	}
}

func TestClientGetMatchReturnsSharedHTTPError(t *testing.T) {
	const body = `{"status":"not found"}`
	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return response(http.StatusNotFound, "404 Not Found", body), nil
	})}

	got, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetMatch("asia", "ASIA_1")
	if got != nil {
		t.Fatalf("match = %#v, want nil on HTTP error", got)
	}
	var responseErr *client.HTTPError
	if !errors.As(err, &responseErr) {
		t.Fatalf("error = %v, want *client.HTTPError", err)
	}
	if responseErr.StatusCode != http.StatusNotFound {
		t.Fatalf("status code = %d, want %d", responseErr.StatusCode, http.StatusNotFound)
	}
	if string(responseErr.Body) != body {
		t.Fatalf("body = %q, want %q", responseErr.Body, body)
	}
}

func TestClientGetMatchReturnsNilOnRequestError(t *testing.T) {
	got, err := New(client.NewWithHTTPClient("test-key", http.DefaultClient)).GetMatch("not/a/region", "ASIA_1")
	if err == nil {
		t.Fatal("expected request error")
	}
	if got != nil {
		t.Fatalf("match = %#v, want nil on request error", got)
	}
}

func TestClientGetTimeline(t *testing.T) {
	fixture := completeTimelineFixture(t)
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		assertRequest(t, r, "/lol/match/v5/matches/ASIA_1%2F%20special/timeline")
		if r.URL.RawQuery != "" {
			t.Fatalf("query = %q, want empty", r.URL.RawQuery)
		}
		return response(http.StatusOK, "200 OK", fixture), nil
	})}

	got, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetTimeline("asia", "ASIA_1/ special")
	if err != nil {
		t.Fatal(err)
	}
	assertSameJSON(t, got, fixture)
}

func TestTimelineResponseWireSchemas(t *testing.T) {
	assertWireSchema(t, reflect.TypeFor[TimelineDto](), timelineSchema)
	assertWireSchema(t, reflect.TypeFor[MetadataTimeLineDto](), metadataTimelineSchema)
	assertWireSchema(t, reflect.TypeFor[InfoTimeLineDto](), infoTimelineSchema)
	assertWireSchema(t, reflect.TypeFor[ParticipantTimeLineDto](), participantTimelineSchema)
	assertWireSchema(t, reflect.TypeFor[FramesTimeLineDto](), framesTimelineSchema)
	assertWireSchema(t, reflect.TypeFor[EventsTimeLineDto](), eventsTimelineSchema)
	assertWireSchema(t, reflect.TypeFor[ParticipantFrameDto](), participantFrameSchema)
	assertWireSchema(t, reflect.TypeFor[ChampionStatsDto](), championStatsSchema)
	assertWireSchema(t, reflect.TypeFor[DamageStatsDto](), damageStatsSchema)
	assertWireSchema(t, reflect.TypeFor[PositionDto](), positionSchema)
}

func TestClientGetTimelineReturnsNilOnPartialDecodeError(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return response(http.StatusOK, "200 OK", `{"metadata":{"dataVersion":"2"},"info":{"gameId":"invalid"}}`), nil
	})}

	got, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetTimeline("asia", "ASIA_1")
	if err == nil {
		t.Fatal("expected partial JSON decode error")
	}
	if got != nil {
		t.Fatalf("timeline = %#v, want nil on decode error", got)
	}
}

func TestClientGetTimelineReturnsSharedHTTPError(t *testing.T) {
	const body = `{"status":"not found"}`
	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return response(http.StatusNotFound, "404 Not Found", body), nil
	})}

	got, err := New(client.NewWithHTTPClient("test-key", httpClient)).GetTimeline("asia", "ASIA_1")
	if got != nil {
		t.Fatalf("timeline = %#v, want nil on HTTP error", got)
	}
	var responseErr *client.HTTPError
	if !errors.As(err, &responseErr) {
		t.Fatalf("error = %v, want *client.HTTPError", err)
	}
	if responseErr.StatusCode != http.StatusNotFound {
		t.Fatalf("status code = %d, want %d", responseErr.StatusCode, http.StatusNotFound)
	}
	if string(responseErr.Body) != body {
		t.Fatalf("body = %q, want %q", responseErr.Body, body)
	}
}

func TestClientGetTimelineReturnsNilOnRequestError(t *testing.T) {
	got, err := New(client.NewWithHTTPClient("test-key", http.DefaultClient)).GetTimeline("not/a/region", "ASIA_1")
	if err == nil {
		t.Fatal("expected request error")
	}
	if got != nil {
		t.Fatalf("timeline = %#v, want nil on request error", got)
	}
}

type wireField struct {
	goName string
	json   string
	typeOf string
}

func schema(spec string) []wireField {
	lines := strings.Split(strings.TrimSpace(spec), "\n")
	fields := make([]wireField, 0, len(lines))
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) != 3 {
			panic("invalid test schema line: " + line)
		}
		fields = append(fields, wireField{goName: parts[0], json: parts[1], typeOf: parts[2]})
	}
	return fields
}

func assertWireSchema(t *testing.T, typ reflect.Type, want []wireField) {
	t.Helper()
	if typ.NumField() != len(want) {
		t.Fatalf("%s has %d fields, want %d", typ, typ.NumField(), len(want))
	}
	for i, field := range want {
		got := typ.Field(i)
		if got.Name != field.goName || got.Tag.Get("json") != field.json || got.Type.String() != field.typeOf {
			t.Errorf("%s field %d = %s %s %q, want %s %s %q", typ, i, got.Name, got.Type, got.Tag.Get("json"), field.goName, field.typeOf, field.json)
		}
	}
}

func completeMatchFixture(t *testing.T) string {
	t.Helper()
	challenges := fixtureObject(challengesSchema, nil)
	missions := fixtureObject(missionsSchema, nil)
	participant := fixtureObject(participantSchema, map[string]any{
		"challenges": challenges,
		"missions":   missions,
		"perks": map[string]any{
			"statPerks": map[string]any{"defense": 601, "flex": 602, "offense": 603},
			"styles": []any{map[string]any{
				"description": "primary-style",
				"selections":  []any{map[string]any{"perk": 604, "var1": 605, "var2": 606, "var3": 607}},
				"style":       608,
			}},
		},
	})
	fixture := map[string]any{
		"metadata": map[string]any{
			"dataVersion":  "2",
			"matchId":      "ASIA_1/ special",
			"participants": []any{"participant-puuid"},
		},
		"info": map[string]any{
			"endOfGameResult":    "GameComplete",
			"gameCreation":       int64(1700000000001),
			"gameDuration":       int64(1801),
			"gameEndTimestamp":   int64(1700001801002),
			"gameId":             int64(12345678901),
			"gameMode":           "CLASSIC",
			"gameName":           "complete fixture",
			"gameStartTimestamp": int64(1700000001003),
			"gameType":           "MATCHED_GAME",
			"gameVersion":        "26.17.1",
			"mapId":              11,
			"participants":       []any{participant},
			"platformId":         "JP1",
			"queueId":            420,
			"teams": []any{map[string]any{
				"bans": []any{map[string]any{"championId": 101, "pickTurn": 1}},
				"objectives": map[string]any{
					"baron":      map[string]any{"first": true, "kills": 1},
					"champion":   map[string]any{"first": true, "kills": 2},
					"dragon":     map[string]any{"first": true, "kills": 3},
					"horde":      map[string]any{"first": true, "kills": 4},
					"inhibitor":  map[string]any{"first": true, "kills": 5},
					"riftHerald": map[string]any{"first": true, "kills": 6},
					"tower":      map[string]any{"first": true, "kills": 7},
				},
				"teamId": 100,
				"win":    true,
			}},
			"tournamentCode": "TOURNAMENT-CODE",
		},
	}
	body, err := json.Marshal(fixture)
	if err != nil {
		t.Fatal(err)
	}
	return string(body)
}

func completeTimelineFixture(t *testing.T) string {
	t.Helper()
	championStats := fixtureObjectAt(championStatsSchema, nil, 100)
	damageStats := fixtureObjectAt(damageStatsSchema, nil, 200)
	position := fixtureObjectAt(positionSchema, nil, 300)
	participantFrame := fixtureObjectAt(participantFrameSchema, map[string]any{
		"championStats": championStats,
		"damageStats":   damageStats,
		"position":      position,
	}, 400)
	fixture := map[string]any{
		"metadata": map[string]any{
			"dataVersion":  "2",
			"matchId":      "ASIA_1/ special",
			"participants": []any{"participant-puuid"},
		},
		"info": map[string]any{
			"endOfGameResult": "GameComplete",
			"frameInterval":   int64(60001),
			"gameId":          int64(12345678901),
			"participants": []any{map[string]any{
				"participantId": 1,
				"puuid":         "participant-puuid",
			}},
			"frames": []any{map[string]any{
				"events": []any{map[string]any{
					"timestamp":     int64(61001),
					"realTimestamp": int64(1700000000002),
					"type":          "GAME_END",
				}},
				"participantFrames": map[string]any{"1": participantFrame},
				"timestamp":         62003,
			}},
		},
	}
	body, err := json.Marshal(fixture)
	if err != nil {
		t.Fatal(err)
	}
	return string(body)
}

func fixtureObject(fields []wireField, replacements map[string]any) map[string]any {
	return fixtureObjectAt(fields, replacements, 1000)
}

func fixtureObjectAt(fields []wireField, replacements map[string]any, base int) map[string]any {
	object := make(map[string]any, len(fields))
	for i, field := range fields {
		if replacement, ok := replacements[field.json]; ok {
			object[field.json] = replacement
			continue
		}
		switch field.typeOf {
		case "int", "int64":
			object[field.json] = base + i
		case "float64":
			object[field.json] = float64(base+i) + 0.25
		case "bool":
			object[field.json] = true
		case "string":
			object[field.json] = field.json + "-value"
		case "[]int":
			object[field.json] = []any{base + i, base + 1000 + i}
		default:
			panic("missing fixture replacement for " + field.goName + " " + field.typeOf)
		}
	}
	return object
}

func assertSameJSON(t *testing.T, got any, want string) {
	t.Helper()
	gotJSON, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	var gotValue, wantValue any
	if err := json.Unmarshal(gotJSON, &gotValue); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(want), &wantValue); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gotValue, wantValue) {
		t.Fatalf("decoded match does not preserve complete fixture\ngot:  %s\nwant: %s", gotJSON, want)
	}
}

var matchSchema = schema(`
Metadata metadata v5.MetadataDto
Info info v5.InfoDto
`)

var timelineSchema = schema(`
Metadata metadata v5.MetadataTimeLineDto
Info info v5.InfoTimeLineDto
`)

var metadataTimelineSchema = schema(`
DataVersion dataVersion string
MatchID matchId string
Participants participants []string
`)

var infoTimelineSchema = schema(`
EndOfGameResult endOfGameResult string
FrameInterval frameInterval int64
GameID gameId int64
Participants participants []v5.ParticipantTimeLineDto
Frames frames []v5.FramesTimeLineDto
`)

var participantTimelineSchema = schema(`
ParticipantID participantId int
PUUID puuid string
`)

var framesTimelineSchema = schema(`
Events events []v5.EventsTimeLineDto
ParticipantFrames participantFrames map[string]v5.ParticipantFrameDto
Timestamp timestamp int
`)

var eventsTimelineSchema = schema(`
Timestamp timestamp int64
RealTimestamp realTimestamp int64
Type type string
`)

var participantFrameSchema = schema(`
ChampionStats championStats v5.ChampionStatsDto
CurrentGold currentGold int
DamageStats damageStats v5.DamageStatsDto
GoldPerSecond goldPerSecond int
JungleMinionsKilled jungleMinionsKilled int
Level level int
MinionsKilled minionsKilled int
ParticipantID participantId int
Position position v5.PositionDto
TimeEnemySpentControlled timeEnemySpentControlled int
TotalGold totalGold int
XP xp int
`)

var championStatsSchema = schema(`
AbilityHaste abilityHaste int
AbilityPower abilityPower int
Armor armor int
ArmorPen armorPen int
ArmorPenPercent armorPenPercent int
AttackDamage attackDamage int
AttackSpeed attackSpeed int
BonusArmorPenPercent bonusArmorPenPercent int
BonusMagicPenPercent bonusMagicPenPercent int
CCReduction ccReduction int
CooldownReduction cooldownReduction int
Health health int
HealthMax healthMax int
HealthRegen healthRegen int
Lifesteal lifesteal int
MagicPen magicPen int
MagicPenPercent magicPenPercent int
MagicResist magicResist int
MovementSpeed movementSpeed int
Omnivamp omnivamp int
PhysicalVamp physicalVamp int
Power power int
PowerMax powerMax int
PowerRegen powerRegen int
SpellVamp spellVamp int
`)

var damageStatsSchema = schema(`
MagicDamageDone magicDamageDone int
MagicDamageDoneToChampions magicDamageDoneToChampions int
MagicDamageTaken magicDamageTaken int
PhysicalDamageDone physicalDamageDone int
PhysicalDamageDoneToChampions physicalDamageDoneToChampions int
PhysicalDamageTaken physicalDamageTaken int
TotalDamageDone totalDamageDone int
TotalDamageDoneToChampions totalDamageDoneToChampions int
TotalDamageTaken totalDamageTaken int
TrueDamageDone trueDamageDone int
TrueDamageDoneToChampions trueDamageDoneToChampions int
TrueDamageTaken trueDamageTaken int
`)

var positionSchema = schema(`
X x int
Y y int
`)

var metadataSchema = schema(`
DataVersion dataVersion string
MatchID matchId string
Participants participants []string
`)

var infoSchema = schema(`
EndOfGameResult endOfGameResult string
GameCreation gameCreation int64
GameDuration gameDuration int64
GameEndTimestamp gameEndTimestamp int64
GameID gameId int64
GameMode gameMode string
GameName gameName string
GameStartTimestamp gameStartTimestamp int64
GameType gameType string
GameVersion gameVersion string
MapID mapId int
Participants participants []v5.ParticipantDto
PlatformID platformId string
QueueID queueId int
Teams teams []v5.TeamDto
TournamentCode tournamentCode string
`)

var teamSchema = schema(`
Bans bans []v5.BanDto
Objectives objectives v5.ObjectivesDto
TeamID teamId int
Win win bool
`)

var banSchema = schema(`
ChampionID championId int
PickTurn pickTurn int
`)

var objectivesSchema = schema(`
Baron baron v5.ObjectiveDto
Champion champion v5.ObjectiveDto
Dragon dragon v5.ObjectiveDto
Horde horde v5.ObjectiveDto
Inhibitor inhibitor v5.ObjectiveDto
RiftHerald riftHerald v5.ObjectiveDto
Tower tower v5.ObjectiveDto
`)

var objectiveSchema = schema(`
First first bool
Kills kills int
`)

var perksSchema = schema(`
StatPerks statPerks v5.PerkStatsDto
Styles styles []v5.PerkStyleDto
`)

var perkStatsSchema = schema(`
Defense defense int
Flex flex int
Offense offense int
`)

var perkStyleSchema = schema(`
Description description string
Selections selections []v5.PerkStyleSelectionDto
Style style int
`)

var perkStyleSelectionSchema = schema(`
Perk perk int
Var1 var1 int
Var2 var2 int
Var3 var3 int
`)

var missionsSchema = schema(`
PlayerScore0 playerScore0 int
PlayerScore1 playerScore1 int
PlayerScore2 playerScore2 int
PlayerScore3 playerScore3 int
PlayerScore4 playerScore4 int
PlayerScore5 playerScore5 int
PlayerScore6 playerScore6 int
PlayerScore7 playerScore7 int
PlayerScore8 playerScore8 int
PlayerScore9 playerScore9 int
PlayerScore10 playerScore10 int
PlayerScore11 playerScore11 int
`)

var participantSchema = schema(`
AllInPings allInPings int
AssistMePings assistMePings int
Assists assists int
BaronKills baronKills int
BountyLevel bountyLevel int
ChampExperience champExperience int
ChampLevel champLevel int
ChampionID championId int
ChampionName championName string
CommandPings commandPings int
ChampionTransform championTransform int
ConsumablesPurchased consumablesPurchased int
Challenges challenges v5.ChallengesDto
DamageDealtToBuildings damageDealtToBuildings int
DamageDealtToObjectives damageDealtToObjectives int
DamageDealtToTurrets damageDealtToTurrets int
DamageSelfMitigated damageSelfMitigated int
Deaths deaths int
DetectorWardsPlaced detectorWardsPlaced int
DoubleKills doubleKills int
DragonKills dragonKills int
EligibleForProgression eligibleForProgression bool
EnemyMissingPings enemyMissingPings int
EnemyVisionPings enemyVisionPings int
FirstBloodAssist firstBloodAssist bool
FirstBloodKill firstBloodKill bool
FirstTowerAssist firstTowerAssist bool
FirstTowerKill firstTowerKill bool
GameEndedInEarlySurrender gameEndedInEarlySurrender bool
GameEndedInSurrender gameEndedInSurrender bool
HoldPings holdPings int
GetBackPings getBackPings int
GoldEarned goldEarned int
GoldSpent goldSpent int
IndividualPosition individualPosition string
InhibitorKills inhibitorKills int
InhibitorTakedowns inhibitorTakedowns int
InhibitorsLost inhibitorsLost int
Item0 item0 int
Item1 item1 int
Item2 item2 int
Item3 item3 int
Item4 item4 int
Item5 item5 int
Item6 item6 int
ItemsPurchased itemsPurchased int
KillingSprees killingSprees int
Kills kills int
Lane lane string
LargestCriticalStrike largestCriticalStrike int
LargestKillingSpree largestKillingSpree int
LargestMultiKill largestMultiKill int
LongestTimeSpentLiving longestTimeSpentLiving int
MagicDamageDealt magicDamageDealt int
MagicDamageDealtToChampions magicDamageDealtToChampions int
MagicDamageTaken magicDamageTaken int
Missions missions v5.MissionsDto
NeutralMinionsKilled neutralMinionsKilled int
NeedVisionPings needVisionPings int
NexusKills nexusKills int
NexusTakedowns nexusTakedowns int
NexusLost nexusLost int
ObjectivesStolen objectivesStolen int
ObjectivesStolenAssists objectivesStolenAssists int
OnMyWayPings onMyWayPings int
ParticipantID participantId int
PlayerScore0 playerScore0 int
PlayerScore1 playerScore1 int
PlayerScore2 playerScore2 int
PlayerScore3 playerScore3 int
PlayerScore4 playerScore4 int
PlayerScore5 playerScore5 int
PlayerScore6 playerScore6 int
PlayerScore7 playerScore7 int
PlayerScore8 playerScore8 int
PlayerScore9 playerScore9 int
PlayerScore10 playerScore10 int
PlayerScore11 playerScore11 int
PentaKills pentaKills int
Perks perks v5.PerksDto
PhysicalDamageDealt physicalDamageDealt int
PhysicalDamageDealtToChampions physicalDamageDealtToChampions int
PhysicalDamageTaken physicalDamageTaken int
Placement placement int
PlayerAugment1 playerAugment1 int
PlayerAugment2 playerAugment2 int
PlayerAugment3 playerAugment3 int
PlayerAugment4 playerAugment4 int
PlayerSubteamID playerSubteamId int
PushPings pushPings int
ProfileIcon profileIcon int
PUUID puuid string
QuadraKills quadraKills int
RiotIDGameName riotIdGameName string
RiotIDTagline riotIdTagline string
Role role string
SightWardsBoughtInGame sightWardsBoughtInGame int
Spell1Casts spell1Casts int
Spell2Casts spell2Casts int
Spell3Casts spell3Casts int
Spell4Casts spell4Casts int
SubteamPlacement subteamPlacement int
Summoner1Casts summoner1Casts int
Summoner1ID summoner1Id int
Summoner2Casts summoner2Casts int
Summoner2ID summoner2Id int
SummonerID summonerId string
SummonerLevel summonerLevel int
SummonerName summonerName string
TeamEarlySurrendered teamEarlySurrendered bool
TeamID teamId int
TeamPosition teamPosition string
TimeCCingOthers timeCCingOthers int
TimePlayed timePlayed int
TotalAllyJungleMinionsKilled totalAllyJungleMinionsKilled int
TotalDamageDealt totalDamageDealt int
TotalDamageDealtToChampions totalDamageDealtToChampions int
TotalDamageShieldedOnTeammates totalDamageShieldedOnTeammates int
TotalDamageTaken totalDamageTaken int
TotalEnemyJungleMinionsKilled totalEnemyJungleMinionsKilled int
TotalHeal totalHeal int
TotalHealsOnTeammates totalHealsOnTeammates int
TotalMinionsKilled totalMinionsKilled int
TotalTimeCCDealt totalTimeCCDealt int
TotalTimeSpentDead totalTimeSpentDead int
TotalUnitsHealed totalUnitsHealed int
TripleKills tripleKills int
TrueDamageDealt trueDamageDealt int
TrueDamageDealtToChampions trueDamageDealtToChampions int
TrueDamageTaken trueDamageTaken int
TurretKills turretKills int
TurretTakedowns turretTakedowns int
TurretsLost turretsLost int
UnrealKills unrealKills int
VisionScore visionScore int
VisionClearedPings visionClearedPings int
VisionWardsBoughtInGame visionWardsBoughtInGame int
WardsKilled wardsKilled int
WardsPlaced wardsPlaced int
Win win bool
`)

var challengesSchema = schema(`
TwelveAssistStreakCount 12AssistStreakCount int
BaronBuffGoldAdvantageOverThreshold baronBuffGoldAdvantageOverThreshold int
ControlWardTimeCoverageInRiverOrEnemyHalf controlWardTimeCoverageInRiverOrEnemyHalf float64
EarliestBaron earliestBaron int
EarliestDragonTakedown earliestDragonTakedown int
EarliestElderDragon earliestElderDragon int
EarlyLaningPhaseGoldExpAdvantage earlyLaningPhaseGoldExpAdvantage int
FasterSupportQuestCompletion fasterSupportQuestCompletion int
FastestLegendary fastestLegendary int
HadAfkTeammate hadAfkTeammate int
HighestChampionDamage highestChampionDamage int
HighestCrowdControlScore highestCrowdControlScore int
HighestWardKills highestWardKills int
JunglerKillsEarlyJungle junglerKillsEarlyJungle int
KillsOnLanersEarlyJungleAsJungler killsOnLanersEarlyJungleAsJungler int
LaningPhaseGoldExpAdvantage laningPhaseGoldExpAdvantage int
LegendaryCount legendaryCount int
MaxCSAdvantageOnLaneOpponent maxCsAdvantageOnLaneOpponent float64
MaxLevelLeadLaneOpponent maxLevelLeadLaneOpponent int
MostWardsDestroyedOneSweeper mostWardsDestroyedOneSweeper int
MythicItemUsed mythicItemUsed int
PlayedChampSelectPosition playedChampSelectPosition int
SoloTurretsLategame soloTurretsLategame int
TakedownsFirst25Minutes takedownsFirst25Minutes int
TeleportTakedowns teleportTakedowns int
ThirdInhibitorDestroyedTime thirdInhibitorDestroyedTime int
ThreeWardsOneSweeperCount threeWardsOneSweeperCount int
VisionScoreAdvantageLaneOpponent visionScoreAdvantageLaneOpponent float64
InfernalScalePickup InfernalScalePickup int
FistBumpParticipation fistBumpParticipation int
VoidMonsterKill voidMonsterKill int
AbilityUses abilityUses int
AcesBefore15Minutes acesBefore15Minutes int
AlliedJungleMonsterKills alliedJungleMonsterKills float64
BaronTakedowns baronTakedowns int
BlastConeOppositeOpponentCount blastConeOppositeOpponentCount int
BountyGold bountyGold int
BuffsStolen buffsStolen int
CompleteSupportQuestInTime completeSupportQuestInTime int
ControlWardsPlaced controlWardsPlaced int
DamagePerMinute damagePerMinute float64
DamageTakenOnTeamPercentage damageTakenOnTeamPercentage float64
DancedWithRiftHerald dancedWithRiftHerald int
DeathsByEnemyChamps deathsByEnemyChamps int
DodgeSkillShotsSmallWindow dodgeSkillShotsSmallWindow int
DoubleAces doubleAces int
DragonTakedowns dragonTakedowns int
LegendaryItemUsed legendaryItemUsed []int
EffectiveHealAndShielding effectiveHealAndShielding float64
ElderDragonKillsWithOpposingSoul elderDragonKillsWithOpposingSoul int
ElderDragonMultikills elderDragonMultikills int
EnemyChampionImmobilizations enemyChampionImmobilizations int
EnemyJungleMonsterKills enemyJungleMonsterKills float64
EpicMonsterKillsNearEnemyJungler epicMonsterKillsNearEnemyJungler int
EpicMonsterKillsWithin30SecondsOfSpawn epicMonsterKillsWithin30SecondsOfSpawn int
EpicMonsterSteals epicMonsterSteals int
EpicMonsterStolenWithoutSmite epicMonsterStolenWithoutSmite int
FirstTurretKilled firstTurretKilled int
FirstTurretKilledTime firstTurretKilledTime float64
FlawlessAces flawlessAces int
FullTeamTakedown fullTeamTakedown int
GameLength gameLength float64
GetTakedownsInAllLanesEarlyJungleAsLaner getTakedownsInAllLanesEarlyJungleAsLaner int
GoldPerMinute goldPerMinute float64
HadOpenNexus hadOpenNexus int
ImmobilizeAndKillWithAlly immobilizeAndKillWithAlly int
InitialBuffCount initialBuffCount int
InitialCrabCount initialCrabCount int
JungleCSBefore10Minutes jungleCsBefore10Minutes float64
JunglerTakedownsNearDamagedEpicMonster junglerTakedownsNearDamagedEpicMonster int
KDA kda float64
KillAfterHiddenWithAlly killAfterHiddenWithAlly int
KilledChampTookFullTeamDamageSurvived killedChampTookFullTeamDamageSurvived int
KillingSprees killingSprees int
KillParticipation killParticipation float64
KillsNearEnemyTurret killsNearEnemyTurret int
KillsOnOtherLanesEarlyJungleAsLaner killsOnOtherLanesEarlyJungleAsLaner int
KillsOnRecentlyHealedByAramPack killsOnRecentlyHealedByAramPack int
KillsUnderOwnTurret killsUnderOwnTurret int
KillsWithHelpFromEpicMonster killsWithHelpFromEpicMonster int
KnockEnemyIntoTeamAndKill knockEnemyIntoTeamAndKill int
KTurretsDestroyedBeforePlatesFall kTurretsDestroyedBeforePlatesFall int
LandSkillShotsEarlyGame landSkillShotsEarlyGame int
LaneMinionsFirst10Minutes laneMinionsFirst10Minutes int
LostAnInhibitor lostAnInhibitor int
MaxKillDeficit maxKillDeficit int
MejaisFullStackInTime mejaisFullStackInTime int
MoreEnemyJungleThanOpponent moreEnemyJungleThanOpponent float64
MultiKillOneSpell multiKillOneSpell int
Multikills multikills int
MultikillsAfterAggressiveFlash multikillsAfterAggressiveFlash int
MultiTurretRiftHeraldCount multiTurretRiftHeraldCount int
OuterTurretExecutesBefore10Minutes outerTurretExecutesBefore10Minutes int
OutnumberedKills outnumberedKills int
OutnumberedNexusKill outnumberedNexusKill int
PerfectDragonSoulsTaken perfectDragonSoulsTaken int
PerfectGame perfectGame int
PickKillWithAlly pickKillWithAlly int
PoroExplosions poroExplosions int
QuickCleanse quickCleanse int
QuickFirstTurret quickFirstTurret int
QuickSoloKills quickSoloKills int
RiftHeraldTakedowns riftHeraldTakedowns int
SaveAllyFromDeath saveAllyFromDeath int
ScuttleCrabKills scuttleCrabKills int
ShortestTimeToAceFromFirstTakedown shortestTimeToAceFromFirstTakedown float64
SkillshotsDodged skillshotsDodged int
SkillshotsHit skillshotsHit int
SnowballsHit snowballsHit int
SoloBaronKills soloBaronKills int
SwarmDefeatAatrox SWARM_DefeatAatrox int
SwarmDefeatBriar SWARM_DefeatBriar int
SwarmDefeatMiniBosses SWARM_DefeatMiniBosses int
SwarmEvolveWeapon SWARM_EvolveWeapon int
SwarmHave3Passives SWARM_Have3Passives int
SwarmKillEnemy SWARM_KillEnemy int
SwarmPickupGold SWARM_PickupGold float64
SwarmReachLevel50 SWARM_ReachLevel50 int
SwarmSurvive15Min SWARM_Survive15Min int
SwarmWinWith5EvolvedWeapons SWARM_WinWith5EvolvedWeapons int
SoloKills soloKills int
StealthWardsPlaced stealthWardsPlaced int
SurvivedSingleDigitHPCount survivedSingleDigitHpCount int
SurvivedThreeImmobilizesInFight survivedThreeImmobilizesInFight int
TakedownOnFirstTurret takedownOnFirstTurret int
Takedowns takedowns int
TakedownsAfterGainingLevelAdvantage takedownsAfterGainingLevelAdvantage int
TakedownsBeforeJungleMinionSpawn takedownsBeforeJungleMinionSpawn int
TakedownsFirstXMinutes takedownsFirstXMinutes int
TakedownsInAlcove takedownsInAlcove int
TakedownsInEnemyFountain takedownsInEnemyFountain int
TeamBaronKills teamBaronKills int
TeamDamagePercentage teamDamagePercentage float64
TeamElderDragonKills teamElderDragonKills int
TeamRiftHeraldKills teamRiftHeraldKills int
TookLargeDamageSurvived tookLargeDamageSurvived int
TurretPlatesTaken turretPlatesTaken int
TurretsTakenWithRiftHerald turretsTakenWithRiftHerald int
TurretTakedowns turretTakedowns int
TwentyMinionsIn3SecondsCount twentyMinionsIn3SecondsCount int
TwoWardsOneSweeperCount twoWardsOneSweeperCount int
UnseenRecalls unseenRecalls int
VisionScorePerMinute visionScorePerMinute float64
WardsGuarded wardsGuarded int
WardTakedowns wardTakedowns int
WardTakedownsBefore20M wardTakedownsBefore20M int
`)

func assertRequest(t *testing.T, r *http.Request, path string) {
	t.Helper()
	if r.Method != http.MethodGet {
		t.Fatalf("method = %q, want GET", r.Method)
	}
	if r.URL.Host != "asia.api.riotgames.com" {
		t.Fatalf("host = %q, want asia.api.riotgames.com", r.URL.Host)
	}
	if r.URL.EscapedPath() != path {
		t.Fatalf("escaped path = %q, want %q", r.URL.EscapedPath(), path)
	}
	if r.Header.Get("X-Riot-Token") != "test-key" {
		t.Fatalf("X-Riot-Token = %q, want test-key", r.Header.Get("X-Riot-Token"))
	}
}

func response(statusCode int, status, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Status:     status,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}
