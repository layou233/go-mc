package registry

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/Tnze/go-mc/chat"
	"github.com/Tnze/go-mc/nbt"
)

type NetworkCodec struct {
	ChatType        Registry[ChatType]       `registry:"minecraft:chat_type"`
	DamageType      Registry[DamageType]     `registry:"minecraft:damage_type"`
	DimensionType   Registry[Dimension]      `registry:"minecraft:dimension_type"`
	TrimMaterial    Registry[nbt.RawMessage] `registry:"minecraft:trim_material"`
	TrimPattern     Registry[nbt.RawMessage] `registry:"minecraft:trim_pattern"`
	WorldGenBiome   Registry[nbt.RawMessage] `registry:"minecraft:worldgen/biome"`
	Wolfvariant     Registry[nbt.RawMessage] `registry:"minecraft:wolf_variant"`
	PaintingVariant Registry[nbt.RawMessage] `registry:"minecraft:painting_variant"`
	BannerPattern   Registry[nbt.RawMessage] `registry:"minecraft:banner_pattern"`
	Enchantment     Registry[nbt.RawMessage] `registry:"minecraft:enchantment"`
	JukeboxSong     Registry[nbt.RawMessage] `registry:"minecraft:jukebox_song"`
}

type ChatType struct {
	Chat      chat.Decoration `nbt:"chat"`
	Narration chat.Decoration `nbt:"narration"`
}

type DamageType struct {
	MessageID        string  `nbt:"message_id"`
	Scaling          string  `nbt:"scaling"`
	Exhaustion       float32 `nbt:"exhaustion"`
	Effects          string  `nbt:"effects,omitempty"`
	DeathMessageType string  `nbt:"death_message_type,omitempty"`
}

type Dimension struct {
	FixedTime          int64   `nbt:"fixed_time,omitempty"`
	HasSkylight        bool    `nbt:"has_skylight"`
	HasCeiling         bool    `nbt:"has_ceiling"`
	Ultrawarm          bool    `nbt:"ultrawarm"`
	Natural            bool    `nbt:"natural"`
	CoordinateScale    float64 `nbt:"coordinate_scale"`
	BedWorks           bool    `nbt:"bed_works"`
	RespawnAnchorWorks byte    `nbt:"respawn_anchor_works"`
	MinY               int32   `nbt:"min_y"`
	Height             int32   `nbt:"height"`
	LogicalHeight      int32   `nbt:"logical_height"`
	InfiniteBurn       string  `nbt:"infiniburn"`
	Effects            string  `nbt:"effects"`
	AmbientLight       float64 `nbt:"ambient_light"`

	PiglinSafe                  byte           `nbt:"piglin_safe"`
	HasRaids                    byte           `nbt:"has_raids"`
	MonsterSpawnLightLevel      nbt.RawMessage `nbt:"monster_spawn_light_level"` // Tag_Int or {type:"minecraft:uniform", value:{min_inclusive: Tag_Int, max_inclusive: Tag_Int}}
	MonsterSpawnBlockLightLimit int32          `nbt:"monster_spawn_block_light_limit"`
}

// InsertNBTDataIntoRegistry insert data (entry, data) into the registry.
// The codec should be a pointer of a struct. And the registry should be a field of the codec struct.
//
// This function is a temporary solution while the registry system isn't implemented well.
func InsertNBTDataIntoRegistry(codec any, registry, entry string, data nbt.RawMessage) error {
	codecVal := reflect.ValueOf(codec)
	if codecVal.Kind() != reflect.Pointer {
		return errors.New("codec is not a pointer")
	}

	codecVal = codecVal.Elem()
	if codecVal.Kind() != reflect.Struct {
		return errors.New("codec is not a pointer of struct")
	}

	codecTyp := codecVal.Type()

	numField := codecVal.NumField()
	for i := 0; i < numField; i++ {
		registryID, ok := codecTyp.Field(i).Tag.Lookup("registry")
		if !ok {
			continue
		}
		if registryID == registry {
			fieldVal := codecVal.Field(i).Addr()
			args := []reflect.Value{reflect.ValueOf(entry), reflect.ValueOf(data)}
			err := fieldVal.MethodByName("InsertNBT").Call(args)[0]
			if !err.IsNil() {
				return err.Interface().(error)
			}
			if registry == "minecraft:chat_type" {
				fmt.Println(fieldVal.Interface())
			}
			return nil
		}
	}
	return errors.New("registry " + registry + " not found in the codec")
}
