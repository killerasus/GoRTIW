package RTIW

import (
	"math"
	"math/rand"

	"github.com/engoengine/glm"
)

type Sphere struct {
	Center   glm.Vec3
	Radius   float32
	Material Material
}

func NewSphere(center glm.Vec3, radius float32, material Material) *Sphere {
	s := Sphere{Center: center, Radius: radius, Material: material}
	return &s
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

func (s *Sphere) GetMaterial() Material {
	return s.Material
}

func RandomInUnitSphere(r *rand.Rand) glm.Vec3 {
	onesVec := glm.Vec3{1, 1, 1}
	p := glm.Vec3{r.Float32(), r.Float32(), r.Float32()}
	p.SubWith(&onesVec)
	p.MulWith(2)

	for p.Len2() >= 1.0 {
		p = glm.Vec3{r.Float32(), r.Float32(), r.Float32()}
		p.SubWith(&onesVec)
		p.MulWith(2)
	}

	return p
}
