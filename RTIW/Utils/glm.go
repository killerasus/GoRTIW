package Utils

import (
	"image/color"

	"github.com/engoengine/glm"
)

func Reflect(v, n *glm.Vec3) glm.Vec3 {
	mult := 2 * v.Dot(n)
	multN := n.Mul(mult)
	return v.Sub(&multN)
}

func ColorRGBAFromVec3(vec glm.Vec3) color.RGBA {
	return color.RGBA{R: uint8(255.99 * vec.X()), G: uint8(255.99 * vec.Y()), B: uint8(255.99 * vec.Z()), A: uint8(255)}
}
