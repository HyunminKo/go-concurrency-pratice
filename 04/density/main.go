package main

import "fmt"

type Metal struct {
	mass   float64
	volume float64
}

func (m *Metal) Density() float64 {
	return m.mass / m.volume
}

type Gas struct {
	pressure      float64
	temperature   float64
	molecularMass float64
}

func (g *Gas) Density() float64 {
	var density float64
	density = (g.molecularMass * g.pressure) / (0.0821 * (g.temperature + 273))
	return density
}

type Dense interface {
	Density() float64
}

func IsDenser(a, b Dense) bool {
	return a.Density() > b.Density()
}

func main() {
	gold := Metal{478, 24}
	silver := Metal{100, 10}

	result := IsDenser(&gold, &silver)
	if result {
		fmt.Println("Gold has higher density than silver")
	} else {
		fmt.Println("Silver has higher density than gold")
	}

	oxygen := Gas{
		pressure:      5,
		temperature:   27,
		molecularMass: 32,
	}

	hydrogen := Gas{
		pressure:      1,
		temperature:   0,
		molecularMass: 2,
	}

	result = IsDenser(&oxygen, &hydrogen)
	if result {
		fmt.Println("Oxygen has higher density than hydrogen")
	} else {
		fmt.Println("Hydrogen has higher density than Oxygen")
	}

}
