package Materials

import (
	"RTIW/RTIW"
	"RTIW/RTIW/Utils"
	"math/rand"

	"github.com/engoengine/glm"
)

type Metal struct {
	Albedo glm.Vec3
	Fuzz   float32
}

func NewMetal(v glm.Vec3, f float32) *Metal {
	m := Metal{Albedo: v}
	if f < 1.0 {
		m.Fuzz = f
	} else {
		m.Fuzz = 1.0
	}
	return &m
}

func (m *Metal) Scatter(ray *RTIW.Ray, hr *RTIW.HitRecord, attenuation *glm.Vec3, scatter *RTIW.Ray, rand *rand.Rand) bool {
	normDirection := ray.Direction.Normalized()
	reflected := Utils.Reflect(&normDirection, &hr.Normal)
	inSphere := RTIW.RandomInUnitSphere(rand)
	reflected.AddScaledVec(m.Fuzz, &inSphere)
	*scatter = RTIW.Ray{Origin: hr.P, Direction: reflected}
	*attenuation = m.Albedo
	return scatter.Direction.Dot(&hr.Normal) > 0
}
