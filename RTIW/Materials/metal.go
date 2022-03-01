package Materials

import (
	"RTIW/RTIW"
	"RTIW/RTIW/Utils"
	"math/rand"

	"github.com/engoengine/glm"
)

type Metal struct {
	Albedo glm.Vec3
}

func NewMetal(v glm.Vec3) *Metal {
	m := Metal{Albedo: v}
	return &m
}

func (m *Metal) Scatter(ray *RTIW.Ray, hr *RTIW.HitRecord, attenuation *glm.Vec3, scatter *RTIW.Ray, rand *rand.Rand) bool {
	normDirection := ray.Direction.Normalized()
	*scatter = RTIW.Ray{Origin: hr.P, Direction: Utils.Reflect(&normDirection, &hr.Normal)}
	*attenuation = m.Albedo
	return scatter.Direction.Dot(&hr.Normal) > 0
}
