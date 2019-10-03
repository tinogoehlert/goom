package game

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// ThingDef a DOOM thing
type ThingDef struct {
	ID        int    `yaml:"id"`
	Sprite    string `yaml:"sprite"`
	Animation string `yaml:"anim"`
}

// MonsterDef monster definitions
type MonsterDef struct {
	ID         int               `yaml:"id"`
	Health     int               `yaml:"health"`
	Sprite     string            `yaml:"sprite"`
	Animations map[string]string `yaml:"anim"`
}

// DefStore holds DOOM definitions e.g. monsters, weapons and obstacles
type DefStore struct {
	Monsters  []MonsterDef `yaml:"monsters"`
	Obstacles []ThingDef   `yaml:"obstacles"`
}

func NewDefStore(defConfig string) *DefStore {
	yamlFile, err := ioutil.ReadFile(defConfig)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	var ds = &DefStore{}
	err = yaml.Unmarshal(yamlFile, ds)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return ds
}

func (ds *DefStore) GetMonsterDef(id int) *MonsterDef {
	for _, md := range ds.Monsters {
		if md.ID == id {
			return &md
		}
	}
	return nil
}

func (ds *DefStore) GetObstacleDef(id int) *ThingDef {
	for _, obs := range ds.Obstacles {
		if obs.ID == id {
			return &obs
		}
	}
	return nil
}
