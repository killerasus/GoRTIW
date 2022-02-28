package RTIW

import "github.com/engoengine/glm"

type Camera struct {
	Origin, LowerLeftCorner, Horizontal, Vertical glm.Vec3
}

func (c *Camera) GetRay(u, v float32) Ray {
	direction := c.Horizontal.Mul(u)
	direction.AddScaledVec(v, &c.Vertical)
	direction.Sub(&c.Origin)
	direction.AddWith(&c.LowerLeftCorner)
	return Ray{Origin: c.Origin, Direction: direction}
}
