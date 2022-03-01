package RTIW

import (
	"math/rand"

	"github.com/engoengine/glm"
)

type Material interface {
	Scatter(ray *Ray, h *HitRecord, attenuation *glm.Vec3, scatter *Ray, rand *rand.Rand) bool
}
