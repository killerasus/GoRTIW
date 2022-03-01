package Shapes

import (
	"RTIW/RTIW"
	"math"

	"github.com/engoengine/glm"
)

type Sphere struct {
	Center   glm.Vec3
	Radius   float32
	Material RTIW.Material
}

func NewSphere(center glm.Vec3, radius float32, material RTIW.Material) *Sphere {
	s := Sphere{Center: center, Radius: radius, Material: material}
	return &s
}

func (s *Sphere) Hit(ray *RTIW.Ray, tMin, tMax float32, hit *RTIW.HitRecord) bool {
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

func (s *Sphere) GetMaterial() RTIW.Material {
	return s.Material
}
