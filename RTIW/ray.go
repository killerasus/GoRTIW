package RTIW

import (
	"github.com/engoengine/glm"
)

type Ray struct {
	Origin, Direction glm.Vec3
}

func (r *Ray) PointAtParameter(t float32) glm.Vec3 {
	param := r.Direction.Mul(t)
	return r.Origin.Add(&param)
}
