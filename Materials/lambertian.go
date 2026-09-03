package Materials

import (
	"GoRTIW"
	"GoRTIW/Utils"
	"math/rand"

	"github.com/engoengine/glm"
)

type Lambertian struct {
	Albedo glm.Vec3
}

func NewLambertian(v glm.Vec3) *Lambertian {
	return &Lambertian{Albedo: v}
}

func (l *Lambertian) Scatter(ray *GoRTIW.Ray, hr *GoRTIW.HitRecord, attenuation *glm.Vec3, scatter *GoRTIW.Ray, rand *rand.Rand) bool {
	target := hr.P
	target.AddWith(&hr.Normal)
	point := Utils.RandomInUnitSphere(rand)
	target.AddWith(&point)
	target.SubWith(&hr.P)
	*scatter = GoRTIW.Ray{Origin: hr.P, Direction: target}
	*attenuation = l.Albedo
	return true
}
