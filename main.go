package main

import (
	"RTIW/RTIW"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"

	"github.com/engoengine/glm"
)

func computeColor(ray *RTIW.Ray, surfaces *RTIW.Surfaces) color.RGBA {
	hitRecord := RTIW.HitRecord{}
	if surfaces.Hit(ray, 0, math.MaxFloat32, &hitRecord) {
		n := glm.Vec3{hitRecord.Normal.X() + 1, hitRecord.Normal.Y() + 1, hitRecord.Normal.Z() + 1}
		n = n.Mul(0.5)
		return color.RGBA{R: uint8(255.99 * n.X()), G: uint8(255.99 * n.Y()), B: uint8(255.99 * n.Z()), A: uint8(255)}
	}
	unitDirection := ray.Direction.Normalized()
	t := 0.5 * (unitDirection.Y() + 1.0)
	interpA := glm.Vec3{1.0, 1.0, 1.0}
	interpB := glm.Vec3{0.5, 0.7, 1.0}
	paramA := interpA.Mul(1.0 - t)
	paramB := interpB.Mul(t)
	computed := paramA.Add(&paramB)
	return color.RGBA{R: uint8(255.99 * computed.X()), G: uint8(255.99 * computed.Y()), B: uint8(255.99 * computed.Z()), A: 1}
}

func main() {

	file, err := os.Create("output.jpg")
	if err != nil {
		log.Fatal("error creating ouput file: ", err)
	}

	defer file.Close()

	lowerLeftCorner := glm.Vec3{-2.0, -1.0, -1.0}
	horizontal := glm.Vec3{4.0, 0.0, 0.0}
	vertical := glm.Vec3{0.0, 2.0, 0.0}
	origin := glm.Vec3{0.0, 0.0, 0.0}

	nx := 200
	ny := 100

	output := image.NewRGBA(image.Rect(0, 0, nx, ny))

	var surfaces RTIW.Surfaces
	sphere1 := RTIW.Sphere{Center: glm.Vec3{0.0, 0.0, -1}, Radius: 0.5}
	sphere2 := RTIW.Sphere{Center: glm.Vec3{0.0, -100.5, -1}, Radius: 100}
	surfaces.List = append(surfaces.List, &sphere1)
	surfaces.List = append(surfaces.List, &sphere2)

	for j := 0; j < ny; j++ {
		for i := 0; i < nx; i++ {
			u := float32(float32(i) / float32(nx))
			v := float32(float32(j) / float32(ny))
			paramVert := vertical.Mul(v)
			paramHorz := horizontal.Mul(u)
			horzVert := paramHorz.Add(&paramVert)
			direction := lowerLeftCorner.Add(&horzVert)
			ray := RTIW.Ray{Origin: origin, Direction: direction}
			c := computeColor(&ray, &surfaces)
			output.SetRGBA(i, ny-j, c)
		}
	}

	err = jpeg.Encode(file, output, nil)
	if err != nil {
		log.Fatal("error enconding jpg: ", err)
	}
}
