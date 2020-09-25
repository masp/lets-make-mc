package biome

type DimensionDef struct {
	PiglinSafe      byte    `nbt:"piglin_safe"`
	Natural         byte    `nbt:"natural"`
	AmbientLight    float32 `nbt:"ambient_light"`
	FixedTime       float32 `nbt:"fixed_time"`
	Infiniburn      string  `nbt:"infiniburn"`
	RespawnAnchor   byte    `nbt:"respawn_anchor_works"`
	HasSkylight     byte    `nbt:"has_skylight"`
	BedWorks        byte    `nbt:"bed_works"`
	Effects         string  `nbt:"effects"`
	HasRaids        byte    `nbt:"has_raids"`
	LogicalHeight   int32   `nbt:"logical_height"`
	CoordinateScale float32 `nbt:"coordinate_scale"`
	Ultrawarm       float32 `nbt:"ultrawarm"`
	HasCeiling      float32 `nbt:"has_ceiling"`
}

type DimensionRoot struct {
	TypeName   string           `nbt:"type"`
	Dimensions []DimensionEntry `nbt:"value"`
}

type DimensionEntry struct {
	Name    string       `nbt:"name"`
	ID      int32        `nbt:"id"`
	Element DimensionDef `nbt:"element"`
}

type BiomeDef struct {
	Precipitation string `nbt:"precipitation"`
	Effects       struct {
		SkyColor      int32 `nbt:"sky_color"`
		WaterFogColor int32 `nbt:"water_fog_color"`
		FogColor      int32 `nbt:"fog_color"`
		WaterColor    int32 `nbt:"water_color"`
		MoodSound     struct {
			TickDelay         int32   `nbt:"tick_delay"`
			Offset            float64 `nbt:"offset"`
			Sound             string  `nbt:"sound"`
			BlockSearchExtent int32   `nbt:"block_search_extent"`
		} `nbt:"mood_sound"`
	} `nbt:"effects"`
	Depth    float32 `nbt:"depth"`
	Temp     float32 `nbt:"temperature"`
	Scale    float32 `nbt:"scale"`
	Downfall float32 `nbt:"downfall"`
	Category string  `nbt:"category"`
}

type BiomeEntry struct {
	Name    string   `nbt:"name"`
	ID      int32    `nbt:"id"`
	Element BiomeDef `nbt:"element"`
}

type BiomeRoot struct {
	TypeName string       `nbt:"type"`
	Biomes   []BiomeEntry `nbt:"value"`
}

type DimensionBiomeRegistry struct {
	Dimensions DimensionRoot `nbt:"minecraft:dimension_type"`
	Biomes     BiomeRoot     `nbt:"minecraft:worldgen/biome"`
}

func OverworldDimension() DimensionDef {
	return DimensionDef{
		PiglinSafe:      0,
		Natural:         1,
		AmbientLight:    0,
		Infiniburn:      "minecraft:infiniburn_overworld",
		RespawnAnchor:   0,
		HasSkylight:     1,
		BedWorks:        1,
		Effects:         "minecraft:overworld",
		HasRaids:        1,
		LogicalHeight:   256,
		CoordinateScale: 1.0,
		Ultrawarm:       0,
		HasCeiling:      0,
	}
}

func plainBiome() BiomeDef {
	b := BiomeDef{}
	b.Precipitation = "rain"
	b.Effects.SkyColor = 7907327
	b.Effects.WaterFogColor = 329011
	b.Effects.FogColor = 12638463
	b.Effects.WaterColor = 4159204
	b.Effects.MoodSound.TickDelay = 6000
	b.Effects.MoodSound.Offset = 2.0
	b.Effects.MoodSound.Sound = "minecraft:ambient.cave"
	b.Effects.MoodSound.BlockSearchExtent = 8
	b.Depth = 0.125
	b.Temp = 0.8
	b.Scale = 0.05
	b.Downfall = 0.4
	b.Category = "plains"
	return b
}

const (
	OverworldID = 0
	PlainsID    = 0
)

func BuildRegistry() DimensionBiomeRegistry {
	return DimensionBiomeRegistry{
		Dimensions: DimensionRoot{
			TypeName: "minecraft:dimension_type",
			Dimensions: []DimensionEntry{
				{
					Name:    "minecraft:overworld",
					ID:      OverworldID,
					Element: OverworldDimension(),
				},
			},
		},
		Biomes: BiomeRoot{
			TypeName: "minecraft:worldgen/biome",
			Biomes: []BiomeEntry{
				{
					Name:    "minecraft:plains",
					ID:      PlainsID,
					Element: plainBiome(),
				},
			},
		},
	}
}
