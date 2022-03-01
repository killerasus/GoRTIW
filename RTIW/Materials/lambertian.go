package Materials

import (
	"RTIW/RTIW"
	"RTIW/RTIW/Utils"
	"math/rand"

	"github.com/engoengine/glm"
)

type Lambertian struct {
	Albedo glm.Vec3
}

func NewLambertian(v glm.Vec3) *Lambertian {
	l := Lambertian{Albedo: v}
	return &l
}

func (l *Lambertian) Scatter(ray *RTIW.Ray, hr *RTIW.HitRecord, attenuation *glm.Vec3, scatter *RTIW.Ray, rand *rand.Rand) bool {
	target := hr.P
	target.AddWith(&hr.Normal)
	point := Utils.RandomInUnitSphere(rand)
	target.AddWith(&point)
	target.SubWith(&hr.P)
	*scatter = RTIW.Ray{Origin: hr.P, Direction: target}
	*attenuation = l.Albedo
	return true
}
