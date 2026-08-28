package v5

type TimelineDto struct {
	Metadata MetadataTimeLineDto `json:"metadata"`
	Info     InfoTimeLineDto     `json:"info"`
}

type MetadataTimeLineDto struct {
	DataVersion  string   `json:"dataVersion"`
	MatchID      string   `json:"matchId"`
	Participants []string `json:"participants"`
}

type InfoTimeLineDto struct {
	EndOfGameResult string                   `json:"endOfGameResult"`
	FrameInterval   int64                    `json:"frameInterval"`
	GameID          int64                    `json:"gameId"`
	Participants    []ParticipantTimeLineDto `json:"participants"`
	Frames          []FramesTimeLineDto      `json:"frames"`
}

type ParticipantTimeLineDto struct {
	ParticipantID int    `json:"participantId"`
	PUUID         string `json:"puuid"`
}

type FramesTimeLineDto struct {
	Events            []EventsTimeLineDto            `json:"events"`
	ParticipantFrames map[string]ParticipantFrameDto `json:"participantFrames"`
	Timestamp         int                            `json:"timestamp"`
}

type EventsTimeLineDto struct {
	Timestamp     int64  `json:"timestamp"`
	RealTimestamp int64  `json:"realTimestamp"`
	Type          string `json:"type"`
}

type ParticipantFrameDto struct {
	ChampionStats            ChampionStatsDto `json:"championStats"`
	CurrentGold              int              `json:"currentGold"`
	DamageStats              DamageStatsDto   `json:"damageStats"`
	GoldPerSecond            int              `json:"goldPerSecond"`
	JungleMinionsKilled      int              `json:"jungleMinionsKilled"`
	Level                    int              `json:"level"`
	MinionsKilled            int              `json:"minionsKilled"`
	ParticipantID            int              `json:"participantId"`
	Position                 PositionDto      `json:"position"`
	TimeEnemySpentControlled int              `json:"timeEnemySpentControlled"`
	TotalGold                int              `json:"totalGold"`
	XP                       int              `json:"xp"`
}

type ChampionStatsDto struct {
	AbilityHaste         int `json:"abilityHaste"`
	AbilityPower         int `json:"abilityPower"`
	Armor                int `json:"armor"`
	ArmorPen             int `json:"armorPen"`
	ArmorPenPercent      int `json:"armorPenPercent"`
	AttackDamage         int `json:"attackDamage"`
	AttackSpeed          int `json:"attackSpeed"`
	BonusArmorPenPercent int `json:"bonusArmorPenPercent"`
	BonusMagicPenPercent int `json:"bonusMagicPenPercent"`
	CCReduction          int `json:"ccReduction"`
	CooldownReduction    int `json:"cooldownReduction"`
	Health               int `json:"health"`
	HealthMax            int `json:"healthMax"`
	HealthRegen          int `json:"healthRegen"`
	Lifesteal            int `json:"lifesteal"`
	MagicPen             int `json:"magicPen"`
	MagicPenPercent      int `json:"magicPenPercent"`
	MagicResist          int `json:"magicResist"`
	MovementSpeed        int `json:"movementSpeed"`
	Omnivamp             int `json:"omnivamp"`
	PhysicalVamp         int `json:"physicalVamp"`
	Power                int `json:"power"`
	PowerMax             int `json:"powerMax"`
	PowerRegen           int `json:"powerRegen"`
	SpellVamp            int `json:"spellVamp"`
}

type DamageStatsDto struct {
	MagicDamageDone               int `json:"magicDamageDone"`
	MagicDamageDoneToChampions    int `json:"magicDamageDoneToChampions"`
	MagicDamageTaken              int `json:"magicDamageTaken"`
	PhysicalDamageDone            int `json:"physicalDamageDone"`
	PhysicalDamageDoneToChampions int `json:"physicalDamageDoneToChampions"`
	PhysicalDamageTaken           int `json:"physicalDamageTaken"`
	TotalDamageDone               int `json:"totalDamageDone"`
	TotalDamageDoneToChampions    int `json:"totalDamageDoneToChampions"`
	TotalDamageTaken              int `json:"totalDamageTaken"`
	TrueDamageDone                int `json:"trueDamageDone"`
	TrueDamageDoneToChampions     int `json:"trueDamageDoneToChampions"`
	TrueDamageTaken               int `json:"trueDamageTaken"`
}

type PositionDto struct {
	X int `json:"x"`
	Y int `json:"y"`
}
