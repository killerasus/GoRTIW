package Utils

import (
	"image/color"
	"math/rand"

	"github.com/engoengine/glm"
)

func Reflect(v, n *glm.Vec3) glm.Vec3 {
	mult := 2 * v.Dot(n)
	multN := n.Mul(mult)
	return v.Sub(&multN)
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

func ColorRGBAFromVec3(vec glm.Vec3) color.RGBA {
	return color.RGBA{R: uint8(255.99 * vec.X()), G: uint8(255.99 * vec.Y()), B: uint8(255.99 * vec.Z()), A: uint8(255)}
}
