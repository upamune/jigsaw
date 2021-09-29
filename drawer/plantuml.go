package drawer

import "fmt"

const TypePlantUML = "plantuml"

type PlantUML struct{}

func (p *PlantUML) Draw(from, to, msg string) string {
	return fmt.Sprintf(`"%s" -> "%s": %s`, from, to, msg)
}

func (m *PlantUML) Comment(s string) string {
	return fmt.Sprintf("' %s", s)
}

func (p *PlantUML) Header() string {
	return "@startuml"
}

func (p *PlantUML) Footer() string {
	return "@enduml"
}
