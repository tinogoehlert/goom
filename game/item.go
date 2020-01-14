package game

// Item - an item in the game
type Item struct {
	*DoomThing
	category string
	ref      string
	consumed bool
}

// ItemFromDef creates item from definition
func ItemFromDef(x, y, angle float32, def *ItemDef) *Item {
	item := &Item{DoomThing: NewDoomThing(x, y, angle, def.Sprite, false)}
	item.animations["idle"] = []byte(def.Animation)
	item.currentAnimation = item.animations["idle"]
	item.id = def.ID
	item.category = def.Category
	item.ref = def.Reference
	return item
}

// IsShown determines if consumable was consumed
func (item *Item) IsShown() bool {
	return !item.consumed
}
