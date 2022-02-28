package RTIW

import (
	"math"

	"github.com/engoengine/glm"
)

type Sphere struct {
	Center glm.Vec3
	Radius float32
}

func (s *Sphere) Hit(ray *Ray, tMin, tMax float32, hit *HitRecord) bool {
	oc := ray.Origin.Sub(&s.Center)
	a := ray.Direction.Dot(&ray.Direction)
	b := oc.Dot(&ray.Direction)
	c := oc.Dot(&oc) - s.Radius*s.Radius
	discriminant := b*b - a*c
	if discriminant > 0 {
		temp := (-b - float32(math.Sqrt(float64(b*b-a*c)))) / a
		if temp < tMax && temp > tMin {
			hit.T = temp
			hit.P = ray.PointAtParameter(hit.T)
			hit.Normal = hit.P.Sub(&s.Center)
			hit.Normal.Normalize()
			return true
		}
		temp = (-b + float32(math.Sqrt(float64(b*b-a*c)))) / a
		if temp < tMax && temp > tMin {
			hit.T = temp
			hit.P = ray.PointAtParameter(hit.T)
			hit.Normal = hit.P.Sub(&s.Center)
			hit.Normal.Normalize()
			return true
		}
	}
	return false
}
