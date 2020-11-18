package main

import (
	"fmt"
	"math"
)

type circle struct {
	radius float64
}
type triangle struct {
	a, b, c float64
}

type rectangle struct {
	h, w float64
}

func (c circle) String() string {
	return fmt.Sprint("Circle (Radius: ", c.radius, ")")
}

func (t triangle) String() string {
	return fmt.Sprint("Triangle (Sides: ", t.a, ", ", t.b, ", ", t.c, ")")
}

func (r rectangle) String() string {
	return fmt.Sprint("rectangle (Sides: ", r.h, ", ", r.w, ")")
}

func (c circle) area() float64 {
	return math.Pi * c.radius * c.radius
}

func (t triangle) area() float64 {
	p := (t.a + t.b + t.c) / 2.0
	return math.Sqrt(p * (p - t.a) * (p - t.b) * (p - t.c))
}

func (r rectangle) area() float64 {
	return r.h * r.w
}
func (t triangle) angles() []float64 {
	return []float64{angle(t.b, t.c, t.a), angle(t.a, t.c, t.b), angle(t.a, t.b, t.c)}
}

func angle(a, b, c float64) float64 {
	return math.Acos((a*a+b*b-c*c)/(2*a*b)) * 180.0 / math.Pi
}

type shape interface {
	area() float64
}

func main() {

	shapes := []shape{
		circle{1.0},
		triangle{10, 4, 7},
		rectangle{5, 10},
	}

	for _, s := range shapes {
		fmt.Println(s, "\tArea:", s.area())
		t, ok := s.(triangle)
		if ok {
			fmt.Println(t.angles())
		}
	}
}
