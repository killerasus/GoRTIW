package Materials

import (
	"GoRTIW"
	"GoRTIW/Utils"
	"math/rand"

	"github.com/engoengine/glm"
)

type Dielectric struct {
	ReflectionIdx float32
}

func NewDielectric(ri float32) *Dielectric {
	return &Dielectric{ReflectionIdx: ri}
}

func (d *Dielectric) Scatter(ray *GoRTIW.Ray, hr *GoRTIW.HitRecord, attenuation *glm.Vec3, scatter *GoRTIW.Ray, rand *rand.Rand) bool {
	outwardNormal := glm.Vec3{}
	reflected := Utils.Reflect(&ray.Direction, &hr.Normal)
	niOverNt := float32(0.0)
	*attenuation = glm.Vec3{1.0, 1.0, 1.0}
	refracted := glm.Vec3{}
	cosine := float32(0.0)
	reflectProb := float32(0.0)

	if ray.Direction.Dot(&hr.Normal) > 0 {
		outwardNormal = hr.Normal.Inverse()
		niOverNt = d.ReflectionIdx
		cosine = d.ReflectionIdx * ray.Direction.Dot(&hr.Normal) / ray.Direction.Len()
	} else {
		outwardNormal = hr.Normal
		niOverNt = 1.0 / d.ReflectionIdx
		cosine = -ray.Direction.Dot(&hr.Normal) / ray.Direction.Len()
	}

	refracted, ok := Utils.Refract(&ray.Direction, &outwardNormal, niOverNt)
	if ok {
		reflectProb = Utils.Schlick(cosine, d.ReflectionIdx)
	} else {
		reflectProb = 1.0
	}

	if rand.Float32() < reflectProb {
		*scatter = GoRTIW.Ray{Origin: hr.P, Direction: reflected}
	} else {
		*scatter = GoRTIW.Ray{Origin: hr.P, Direction: refracted}
	}

	return true
}
