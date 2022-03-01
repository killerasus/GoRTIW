package RTIW

import (
	"RTIW/RTIW/Utils"
	"math"
	"math/rand"

	"github.com/engoengine/glm"
)

type Camera struct {
	Origin, LowerLeftCorner, Horizontal, Vertical, U, V, W glm.Vec3
	LensRadius                                             float32
}

func NewCamera(lookFrom, lookAt, vUp glm.Vec3, vfov, aspect, aperture, focusDist float32) *Camera {
	c := Camera{}
	c.Origin = lookFrom
	c.LensRadius = aperture / 2.0

	theta := vfov * math.Pi / 180.0
	halfHeight := math.Tan(float64(theta / 2.0))
	halfWidth := aspect * float32(halfHeight)
	localLookAt := lookAt //Can't get the address from copy parameter

	w := lookFrom.Sub(&localLookAt)
	w.Normalize()
	c.W = w

	u := vUp.Cross(&w)
	u.Normalize()
	c.U = u

	v := w.Cross(&u)
	c.V = v

	u.MulWith(halfWidth * focusDist)
	v.MulWith(float32(halfHeight * float64(focusDist)))
	w.MulWith(focusDist)

	lowerLeftCorner := lookFrom.Sub(&u)
	lowerLeftCorner.SubWith(&v)
	lowerLeftCorner.SubWith(&w)

	c.LowerLeftCorner = lowerLeftCorner
	c.Horizontal = u.Mul(2)
	c.Vertical = v.Mul(2)

	return &c
}

func (c *Camera) GetRay(s, t float32, r *rand.Rand) Ray {
	rd := Utils.RandomInUnitSphere(r)
	rd.MulWith(c.LensRadius)
	offset := c.U.Mul(rd.X())
	offset.AddScaledVec(rd.Y(), &c.V)
	origin := c.Origin.Add(&offset)

	direction := c.Horizontal.Mul(s)
	direction.AddScaledVec(t, &c.Vertical)
	direction.SubWith(&c.Origin)
	direction.SubWith(&offset)
	direction.AddWith(&c.LowerLeftCorner)

	return Ray{Origin: origin, Direction: direction}
}
