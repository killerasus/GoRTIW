package RTIW

import "github.com/engoengine/glm"

type HitRecord struct {
	T         float32
	P, Normal glm.Vec3
	Material  Material
}
