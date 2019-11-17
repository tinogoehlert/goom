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
	Sounds     map[string]string `yaml:"sounds"`
	Animations map[string]string `yaml:"anim"`
}

// ItemDef monster definitions
type ItemDef struct {
	ID        int    `yaml:"id"`
	Sprite    string `yaml:"sprite"`
	Animation string `yaml:"anim"`
	Category  string `yaml:"category"`
	Reference string `yaml:"ref"`
}

// DefStore holds DOOM definitions e.g. monsters, weapons and obstacles
type DefStore struct {
	Monsters  []MonsterDef `yaml:"monsters"`
	Obstacles []ThingDef   `yaml:"obstacles"`
	Items     []ItemDef    `yaml:"items"`
	Weapons   []Weapon     `yaml:"weapons"`
}

// NewDefStore creates a new definition store from yaml file
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

// GetMonsterDef gets monster definition by ID
func (ds *DefStore) GetMonsterDef(id int) *MonsterDef {
	for _, md := range ds.Monsters {
		if md.ID == id {
			return &md
		}
	}
	return nil
}

// GetObstacleDef gets obstacle definition by ID
func (ds *DefStore) GetObstacleDef(id int) *ThingDef {
	for _, obs := range ds.Obstacles {
		if obs.ID == id {
			return &obs
		}
	}
	return nil
}

// GetItemDef gets obstacle definition by ID
func (ds *DefStore) GetItemDef(id int) *ItemDef {
	for _, item := range ds.Items {
		if item.ID == id {
			return &item
		}
	}
	return nil
}

// GetWeapon gets weapon definition by Name
func (ds *DefStore) GetWeapon(name string) *Weapon {
	for _, w := range ds.Weapons {
		if w.Name == name {
			return &w
		}
	}
	return nil
}
