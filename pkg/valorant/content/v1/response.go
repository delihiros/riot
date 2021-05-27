package v1

type LocalizedNamesDto struct {
	ARAE string `json:"ar-AE"`
	DEDE string `json:"de-DE"`
	ENGB string `json:"en-GB"`
	ENUS string `json:"en-US"`
	ESES string `json:"es-ES"`
	ESMX string `json:"es-MX"`
	FRFR string `json:"fr-FR"`
	IDID string `json:"id-ID"`
	ITIT string `json:"it-IT"`
	JAJP string `json:"ja-JP"`
	KOKR string `json:"ko-KR"`
	PLPL string `json:"pl-PL"`
	PTBR string `json:"pt-BR"`
	RURU string `json:"ru-RU"`
	THTH string `json:"th-TH"`
	TRTR string `json:"tr-TR"`
	VIVN string `json:"vi-VN"`
	ZHCN string `json:"zh-CN"`
	ZHTW string `json:"zh-TW"`
}

type ContentItemDto struct {
	Name           string             `json:"name"`
	LocalizedNames *LocalizedNamesDto `json:"localizedNames"`
	ID             string             `json:"id"`
	AssetName      string             `json:"assetName"`
	AssetPath      string             `json:"assetPath"`
}

type ActDto struct {
	Name           string
	LocalizedNames *LocalizedNamesDto
	ID             string
	IsActive       bool
}

type ContentDto struct {
	Version      string            `json:"version"`
	Characters   []*ContentItemDto `json:"characters"`
	Maps         []*ContentItemDto `json:"maps"`
	Chromas      []*ContentItemDto `json:"chromas"`
	Skins        []*ContentItemDto `json:"skins"`
	SkinLevels   []*ContentItemDto `json:"skinLevels"`
	Equips       []*ContentItemDto `json:"equips"`
	GameModes    []*ContentItemDto `json:"gameModes"`
	Sprays       []*ContentItemDto `json:"sprays"`
	Charms       []*ContentItemDto `json:"charms"`
	CharmLevels  []*ContentItemDto `json:"charmLevels"`
	PlayerCards  []*ContentItemDto `json:"playerCards"`
	PlayerTitles []*ContentItemDto `json:"playerTitles"`
	Acts         []*ActDto         `json:"acts"`
}
