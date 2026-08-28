# LoL Match and RSO API client implementation plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add Match v5, RSO Match v1, and the Summoner v4 RSO endpoint, covering eight regional and Bearer-authenticated operations.

**Architecture:** Match v5 owns the shared match and timeline DTOs. RSO Match v1 aliases those DTOs and adds Bearer-authenticated request methods, avoiding a second copy of the same wire schema. Summoner v4 gains its independent `/me` Bearer method.

**Tech Stack:** Go 1.26.7, standard library `encoding/json`, `net/http`, `net/url`, and existing `pkg/client` authentication boundaries.

---

**Spec:** `docs/superpowers/specs/2026-08-28-lol-api-client-design.md`

**Depends on:** Task 9 of `docs/superpowers/plans/2026-08-28-lol-standard-api-client.md` for the Summoner v4 package.

**Official contracts:** Riot Developer Portal API details for `match-v5`, `lol-rso-match-v1`, and `summoner-v4`, checked 2026-08-28.

**Required workflow:** Use @superpowers:test-driven-development for every production method. Response fixtures must include every documented field rather than only fields asserted by the immediate test. In Step 1 of every package task, also add one malformed-JSON case that expects a decode error and one non-2xx case that expects the shared `HTTPError`; these cases must call a new package method. Before every GREEN run and task commit, run `gofmt -w` on that task's Go files.

### Task 1: Match v5 ID list and replay endpoints

**Files:**
- Create: `pkg/lol/match/v5/client.go`
- Create: `pkg/lol/match/v5/response.go`
- Create: `pkg/lol/match/v5/client_test.go`

- [ ] **Step 1: Write failing request-option tests**

Define the desired API in tests:

```go
type MatchIDsOptions struct {
    StartTime *int64
    EndTime   *int64
    Queue     *int
    Type      *string
    Start     *int
    Count     *int
}

GetMatchIDsByPUUID(region, puuid string, options *MatchIDsOptions) ([]string, error)
GetReplay(region, puuid string) (*ReplayDTO, error)
```

Assert that nil options produce no query string; non-nil values encode `startTime`, `endTime`, `queue`, `type`, `start`, and `count`; zero-valued pointers remain present; PUUID is escaped; and the API-key header is sent.

Use a complete replay fixture:

```go
type ReplayDTO struct {
    Total         int      `json:"total"`
    MatchFileURLs []string `json:"matchFileURLs"`
}
```

- [ ] **Step 2: Run tests and verify RED**

Run: `GOTOOLCHAIN=go1.26.7 go test ./pkg/lol/match/v5`

Expected: build failure because the package API does not exist.

- [ ] **Step 3: Implement options, query encoding, and both methods**

Use a private `addMatchIDQueries(path string, options *MatchIDsOptions) string` or inline `url.Values`; do not add a generic options abstraction.

- [ ] **Step 4: Run package tests and verify GREEN**

Run: `GOTOOLCHAIN=go1.26.7 go test ./pkg/lol/match/v5`

Expected: PASS.

- [ ] **Step 5: Commit**

Run: `git add pkg/lol/match/v5 && git commit -m "feat(lol): add match list and replay clients"`

### Task 2: Match v5 match response

**Files:**
- Modify: `pkg/lol/match/v5/client.go`
- Modify: `pkg/lol/match/v5/response.go`
- Modify: `pkg/lol/match/v5/client_test.go`
- Create: `pkg/lol/match/v5/response_participant.go`

- [ ] **Step 1: Write a failing complete-match test**

Test:

```go
GetMatch(region, matchID string) (*MatchDto, error)
```

Assert `/lol/match/v5/matches/{matchId}`, escaped match ID, API-key header, and a complete nested response fixture.

- [ ] **Step 2: Run the test and verify RED**

Run: `GOTOOLCHAIN=go1.26.7 go test ./pkg/lol/match/v5 -run TestClientGetMatch`

Expected: build failure for `GetMatch` and its DTOs.

- [ ] **Step 3: Implement match-level DTOs**

Create these exact shapes in `response.go`:

```text
MatchDto: metadata MetadataDto; info InfoDto
MetadataDto: dataVersion string; matchId string; participants []string
InfoDto: endOfGameResult string; gameCreation int64; gameDuration int64; gameEndTimestamp int64; gameId int64; gameMode string; gameName string; gameStartTimestamp int64; gameType string; gameVersion string; mapId int; participants []ParticipantDto; platformId string; queueId int; teams []TeamDto; tournamentCode string
TeamDto: bans []BanDto; objectives ObjectivesDto; teamId int; win bool
BanDto: championId int; pickTurn int
ObjectivesDto: baron ObjectiveDto; champion ObjectiveDto; dragon ObjectiveDto; horde ObjectiveDto; inhibitor ObjectiveDto; riftHerald ObjectiveDto; tower ObjectiveDto
ObjectiveDto: first bool; kills int
PerksDto: statPerks PerkStatsDto; styles []PerkStyleDto
PerkStatsDto: defense int; flex int; offense int
PerkStyleDto: description string; selections []PerkStyleSelectionDto; style int
PerkStyleSelectionDto: perk int; var1 int; var2 int; var3 int
MissionsDto: playerScore0 through playerScore11, all int
```

- [ ] **Step 4: Implement the complete participant DTOs**

`ParticipantDto` must include these documented fields with their exact JSON tags:

```text
allInPings int, assistMePings int, assists int, baronKills int, bountyLevel int, champExperience int, champLevel int, championId int, championName string, commandPings int, championTransform int, consumablesPurchased int, challenges ChallengesDto, damageDealtToBuildings int, damageDealtToObjectives int, damageDealtToTurrets int, damageSelfMitigated int, deaths int, detectorWardsPlaced int, doubleKills int, dragonKills int, eligibleForProgression bool, enemyMissingPings int, enemyVisionPings int, firstBloodAssist bool, firstBloodKill bool, firstTowerAssist bool, firstTowerKill bool, gameEndedInEarlySurrender bool, gameEndedInSurrender bool, holdPings int, getBackPings int, goldEarned int, goldSpent int, individualPosition string, inhibitorKills int, inhibitorTakedowns int, inhibitorsLost int, item0-item6 int, itemsPurchased int, killingSprees int, kills int, lane string, largestCriticalStrike int, largestKillingSpree int, largestMultiKill int, longestTimeSpentLiving int, magicDamageDealt int, magicDamageDealtToChampions int, magicDamageTaken int, missions MissionsDto, neutralMinionsKilled int, needVisionPings int, nexusKills int, nexusTakedowns int, nexusLost int, objectivesStolen int, objectivesStolenAssists int, onMyWayPings int, participantId int, playerScore0-playerScore11 int, pentaKills int, perks PerksDto, physicalDamageDealt int, physicalDamageDealtToChampions int, physicalDamageTaken int, placement int, playerAugment1-playerAugment4 int, playerSubteamId int, pushPings int, profileIcon int, puuid string, quadraKills int, riotIdGameName string, riotIdTagline string, role string, sightWardsBoughtInGame int, spell1Casts-spell4Casts int, subteamPlacement int, summoner1Casts int, summoner1Id int, summoner2Casts int, summoner2Id int, summonerId string, summonerLevel int, summonerName string, teamEarlySurrendered bool, teamId int, teamPosition string, timeCCingOthers int, timePlayed int, totalAllyJungleMinionsKilled int, totalDamageDealt int, totalDamageDealtToChampions int, totalDamageShieldedOnTeammates int, totalDamageTaken int, totalEnemyJungleMinionsKilled int, totalHeal int, totalHealsOnTeammates int, totalMinionsKilled int, totalTimeCCDealt int, totalTimeSpentDead int, totalUnitsHealed int, tripleKills int, trueDamageDealt int, trueDamageDealtToChampions int, trueDamageTaken int, turretKills int, turretTakedowns int, turretsLost int, unrealKills int, visionScore int, visionClearedPings int, visionWardsBoughtInGame int, wardsKilled int, wardsPlaced int, win bool
```

Do not expand ranges into arrays: `item0`, `playerScore0`, and similar names are separate JSON properties and separate Go fields.

- [ ] **Step 5: Implement `ChallengesDto` exactly**

Use `int` for Riot `integer`, `float64` for `float`, and `[]int` for `legendaryItemUsed`. Declare every current documented field:

```text
12AssistStreakCount, baronBuffGoldAdvantageOverThreshold, controlWardTimeCoverageInRiverOrEnemyHalf, earliestBaron, earliestDragonTakedown, earliestElderDragon, earlyLaningPhaseGoldExpAdvantage, fasterSupportQuestCompletion, fastestLegendary, hadAfkTeammate, highestChampionDamage, highestCrowdControlScore, highestWardKills, junglerKillsEarlyJungle, killsOnLanersEarlyJungleAsJungler, laningPhaseGoldExpAdvantage, legendaryCount, maxCsAdvantageOnLaneOpponent, maxLevelLeadLaneOpponent, mostWardsDestroyedOneSweeper, mythicItemUsed, playedChampSelectPosition, soloTurretsLategame, takedownsFirst25Minutes, teleportTakedowns, thirdInhibitorDestroyedTime, threeWardsOneSweeperCount, visionScoreAdvantageLaneOpponent, InfernalScalePickup, fistBumpParticipation, voidMonsterKill, abilityUses, acesBefore15Minutes, alliedJungleMonsterKills, baronTakedowns, blastConeOppositeOpponentCount, bountyGold, buffsStolen, completeSupportQuestInTime, controlWardsPlaced, damagePerMinute, damageTakenOnTeamPercentage, dancedWithRiftHerald, deathsByEnemyChamps, dodgeSkillShotsSmallWindow, doubleAces, dragonTakedowns, legendaryItemUsed, effectiveHealAndShielding, elderDragonKillsWithOpposingSoul, elderDragonMultikills, enemyChampionImmobilizations, enemyJungleMonsterKills, epicMonsterKillsNearEnemyJungler, epicMonsterKillsWithin30SecondsOfSpawn, epicMonsterSteals, epicMonsterStolenWithoutSmite, firstTurretKilled, firstTurretKilledTime, flawlessAces, fullTeamTakedown, gameLength, getTakedownsInAllLanesEarlyJungleAsLaner, goldPerMinute, hadOpenNexus, immobilizeAndKillWithAlly, initialBuffCount, initialCrabCount, jungleCsBefore10Minutes, junglerTakedownsNearDamagedEpicMonster, kda, killAfterHiddenWithAlly, killedChampTookFullTeamDamageSurvived, killingSprees, killParticipation, killsNearEnemyTurret, killsOnOtherLanesEarlyJungleAsLaner, killsOnRecentlyHealedByAramPack, killsUnderOwnTurret, killsWithHelpFromEpicMonster, knockEnemyIntoTeamAndKill, kTurretsDestroyedBeforePlatesFall, landSkillShotsEarlyGame, laneMinionsFirst10Minutes, lostAnInhibitor, maxKillDeficit, mejaisFullStackInTime, moreEnemyJungleThanOpponent, multiKillOneSpell, multikills, multikillsAfterAggressiveFlash, multiTurretRiftHeraldCount, outerTurretExecutesBefore10Minutes, outnumberedKills, outnumberedNexusKill, perfectDragonSoulsTaken, perfectGame, pickKillWithAlly, poroExplosions, quickCleanse, quickFirstTurret, quickSoloKills, riftHeraldTakedowns, saveAllyFromDeath, scuttleCrabKills, shortestTimeToAceFromFirstTakedown, skillshotsDodged, skillshotsHit, snowballsHit, soloBaronKills, SWARM_DefeatAatrox, SWARM_DefeatBriar, SWARM_DefeatMiniBosses, SWARM_EvolveWeapon, SWARM_Have3Passives, SWARM_KillEnemy, SWARM_PickupGold, SWARM_ReachLevel50, SWARM_Survive15Min, SWARM_WinWith5EvolvedWeapons, soloKills, stealthWardsPlaced, survivedSingleDigitHpCount, survivedThreeImmobilizesInFight, takedownOnFirstTurret, takedowns, takedownsAfterGainingLevelAdvantage, takedownsBeforeJungleMinionSpawn, takedownsFirstXMinutes, takedownsInAlcove, takedownsInEnemyFountain, teamBaronKills, teamDamagePercentage, teamElderDragonKills, teamRiftHeraldKills, tookLargeDamageSurvived, turretPlatesTaken, turretsTakenWithRiftHerald, turretTakedowns, twentyMinionsIn3SecondsCount, twoWardsOneSweeperCount, unseenRecalls, visionScorePerMinute, wardsGuarded, wardTakedowns, wardTakedownsBefore20M
```

Go identifiers that begin with a digit or contain uppercase portal prefixes must be made legal without changing their JSON tags, for example `TwelveAssistStreakCount json:"12AssistStreakCount"` and `SwarmDefeatAatrox json:"SWARM_DefeatAatrox"`.

- [ ] **Step 6: Implement `GetMatch` and verify GREEN**

Run: `gofmt -w pkg/lol/match/v5 && GOTOOLCHAIN=go1.26.7 go test ./pkg/lol/match/v5 -run TestClientGetMatch`

Expected: PASS with the complete fixture.

- [ ] **Step 7: Commit**

Run: `git add pkg/lol/match/v5 && git commit -m "feat(lol): add match v5 result client"`

### Task 3: Match v5 timeline

**Files:**
- Modify: `pkg/lol/match/v5/client.go`
- Modify: `pkg/lol/match/v5/client_test.go`
- Create: `pkg/lol/match/v5/response_timeline.go`

- [ ] **Step 1: Write a failing complete-timeline test**

Test `GetTimeline(region, matchID string) (*TimelineDto, error)` and the escaped `/lol/match/v5/matches/{matchId}/timeline` path.

- [ ] **Step 2: Run the test and verify RED**

Run: `GOTOOLCHAIN=go1.26.7 go test ./pkg/lol/match/v5 -run TestClientGetTimeline`

Expected: build failure for the missing timeline method and DTOs.

- [ ] **Step 3: Implement the complete timeline model**

Create:

```text
TimelineDto: metadata MetadataTimeLineDto; info InfoTimeLineDto
MetadataTimeLineDto: dataVersion string; matchId string; participants []string
InfoTimeLineDto: endOfGameResult string; frameInterval int64; gameId int64; participants []ParticipantTimeLineDto; frames []FramesTimeLineDto
ParticipantTimeLineDto: participantId int; puuid string
FramesTimeLineDto: events []EventsTimeLineDto; participantFrames map[string]ParticipantFrameDto; timestamp int
EventsTimeLineDto: timestamp int64; realTimestamp int64; type string
ParticipantFrameDto: championStats ChampionStatsDto; currentGold int; damageStats DamageStatsDto; goldPerSecond int; jungleMinionsKilled int; level int; minionsKilled int; participantId int; position PositionDto; timeEnemySpentControlled int; totalGold int; xp int
ChampionStatsDto: abilityHaste int; abilityPower int; armor int; armorPen int; armorPenPercent int; attackDamage int; attackSpeed int; bonusArmorPenPercent int; bonusMagicPenPercent int; ccReduction int; cooldownReduction int; health int; healthMax int; healthRegen int; lifesteal int; magicPen int; magicPenPercent int; magicResist int; movementSpeed int; omnivamp int; physicalVamp int; power int; powerMax int; powerRegen int; spellVamp int
DamageStatsDto: magicDamageDone int; magicDamageDoneToChampions int; magicDamageTaken int; physicalDamageDone int; physicalDamageDoneToChampions int; physicalDamageTaken int; totalDamageDone int; totalDamageDoneToChampions int; totalDamageTaken int; trueDamageDone int; trueDamageDoneToChampions int; trueDamageTaken int
PositionDto: x int; y int
```

Use a map for `participantFrames` because its JSON object keys are participant IDs.

- [ ] **Step 4: Implement the method and verify GREEN**

Run: `gofmt -w pkg/lol/match/v5 && GOTOOLCHAIN=go1.26.7 go test ./pkg/lol/match/v5`

Expected: PASS.

- [ ] **Step 5: Commit**

Run: `git add pkg/lol/match/v5 && git commit -m "feat(lol): add match v5 timeline client"`

### Task 4: RSO Match v1

**Files:**
- Create: `pkg/lol/rso/match/v1/client.go`
- Create: `pkg/lol/rso/match/v1/response.go`
- Create: `pkg/lol/rso/match/v1/client_test.go`

- [ ] **Step 1: Write failing Bearer-authentication tests**

Test:

```go
GetMatchIDs(region, accessToken string, options *MatchIDsOptions) ([]string, error)
GetMatch(region, accessToken, matchID string) (*MatchDto, error)
GetTimeline(region, accessToken, matchID string) (*TimelineDto, error)
```

Assert the three `/lol/rso-match/v1` paths, the same six optional ID filters, `Authorization: Bearer access-token`, and an absent `X-Riot-Token` header.

- [ ] **Step 2: Run tests and verify RED**

Run: `GOTOOLCHAIN=go1.26.7 go test ./pkg/lol/rso/match/v1`

Expected: build failure.

- [ ] **Step 3: Implement type aliases and methods**

Alias the stable wire types instead of duplicating them:

```go
type MatchIDsOptions = matchv5.MatchIDsOptions
type MatchDto = matchv5.MatchDto
type TimelineDto = matchv5.TimelineDto
```

Use `GetWithRegionAndHeaders` with only the Bearer header. Reuse equivalent query encoding locally or through an exported pure encoder only if the Match v5 implementation already needs one; do not expose transport internals.

- [ ] **Step 4: Run package tests and verify GREEN**

Run: `GOTOOLCHAIN=go1.26.7 go test ./pkg/lol/rso/match/v1`

Expected: PASS.

- [ ] **Step 5: Commit**

Run: `git add pkg/lol/rso/match/v1 && git commit -m "feat(lol): add RSO match v1 client"`

### Task 5: Summoner v4 RSO endpoint

**Files:**
- Modify: `pkg/lol/summoner/v4/client.go`
- Modify: `pkg/lol/summoner/v4/client_test.go`

- [ ] **Step 1: Write a failing `/me` test**

Test:

```go
GetByAccessToken(region, accessToken string) (*SummonerDTO, error)
```

Assert `/lol/summoner/v4/summoners/me`, Bearer authentication, no Riot API-key header, and complete `SummonerDTO` decoding.

- [ ] **Step 2: Run the test and verify RED**

Run: `GOTOOLCHAIN=go1.26.7 go test ./pkg/lol/summoner/v4 -run TestClientGetByAccessToken`

Expected: build failure for the missing method.

- [ ] **Step 3: Implement the Bearer method**

Use `GetWithRegionAndHeaders` and prepend `Bearer ` inside the method, matching `pkg/account/v1`.

- [ ] **Step 4: Run package tests and verify GREEN**

Run: `GOTOOLCHAIN=go1.26.7 go test ./pkg/lol/summoner/v4`

Expected: PASS.

- [ ] **Step 5: Commit**

Run: `git add pkg/lol/summoner/v4 && git commit -m "feat(lol): add summoner RSO endpoint"`

### Task 6: Match and RSO integration verification

**Files:**
- Verify only: Match v5, RSO Match v1, and Summoner v4 packages

- [ ] **Step 1: Verify all eight operations are present**

Count Match v5 = 4, RSO Match v1 = 3, and Summoner RSO = 1.

- [ ] **Step 2: Run full tests and race detector**

Run: `GOTOOLCHAIN=go1.26.7 go test ./...`

Run: `GOTOOLCHAIN=go1.26.7 go test -race ./...`

Expected: both PASS.

- [ ] **Step 3: Run static analysis and formatting checks**

Run: `GOTOOLCHAIN=go1.26.7 go vet ./...`

Run: `gofmt -l pkg/lol/match pkg/lol/rso pkg/lol/summoner`

Expected: no output; every task formatted its files before committing.

Run: `git diff --check`

Expected: no errors.
