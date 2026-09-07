package Utils

import (
	"image/color"
	"math"
	"math/rand"
	"testing"

	"github.com/engoengine/glm"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func f32Near(a, b, eps float32) bool {
	d := a - b
	if d < 0 {
		d = -d
	}
	return d < eps
}

func vec3Near(a, b glm.Vec3, eps float32) bool {
	return f32Near(a[0], b[0], eps) && f32Near(a[1], b[1], eps) && f32Near(a[2], b[2], eps)
}

func vec3Len(v glm.Vec3) float32 {
	return float32(math.Sqrt(float64(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])))
}

// ---------------------------------------------------------------------------
// Reflect
// ---------------------------------------------------------------------------

// TestReflectPerpendicular verifies that reflecting a downward vector off a
// straight-up normal yields an upward vector with the same magnitude.
func TestReflectPerpendicular(t *testing.T) {
	v := glm.Vec3{0, -1, 0}
	n := glm.Vec3{0, 1, 0}
	got := Reflect(&v, &n)
	want := glm.Vec3{0, 1, 0}
	if !vec3Near(got, want, 1e-6) {
		t.Errorf("Reflect(%v, %v) = %v; want %v", v, n, got, want)
	}
}

// TestReflectDiagonal verifies a 45° reflection off a vertical normal.
// Incoming (1,-1,0) reflected off (0,1,0) should give (1,1,0).
func TestReflectDiagonal(t *testing.T) {
	v := glm.Vec3{1, -1, 0}
	n := glm.Vec3{0, 1, 0}
	got := Reflect(&v, &n)
	want := glm.Vec3{1, 1, 0}
	if !vec3Near(got, want, 1e-6) {
		t.Errorf("Reflect(%v, %v) = %v; want %v", v, n, got, want)
	}
}

// TestReflectPreservesMagnitude verifies |reflect(v,n)| == |v|.
func TestReflectPreservesMagnitude(t *testing.T) {
	v := glm.Vec3{3, -4, 0}
	n := glm.Vec3{0, 1, 0}
	got := Reflect(&v, &n)
	if !f32Near(vec3Len(got), vec3Len(v), 1e-5) {
		t.Errorf("Reflect magnitude: got %v, want %v", vec3Len(got), vec3Len(v))
	}
}

// TestReflectGrazingAngle verifies that a purely horizontal ray reflected off
// a vertical normal is unchanged (the normal component is zero).
func TestReflectGrazingAngle(t *testing.T) {
	v := glm.Vec3{1, 0, 0}
	n := glm.Vec3{0, 1, 0}
	got := Reflect(&v, &n)
	// v · n = 0, so reflect returns v unchanged
	if !vec3Near(got, v, 1e-6) {
		t.Errorf("Reflect grazing: got %v, want %v", got, v)
	}
}

// ---------------------------------------------------------------------------
// Refract
// ---------------------------------------------------------------------------

// TestRefractOkWhenDiscriminantPositive verifies ok==true for a typical
// air-to-glass scenario (niOverNt < 1, angle not at total-internal-reflection).
func TestRefractOkWhenDiscriminantPositive(t *testing.T) {
	v := glm.Vec3{0, -1, 0}   // straight down
	n := glm.Vec3{0, 1, 0}    // upward normal
	_, ok := Refract(&v, &n, 1.0/1.5) // air → glass
	if !ok {
		t.Error("Expected Refract to succeed (ok=true) for straight-on incidence")
	}
}

// TestRefractTotalInternalReflection verifies ok==false when the angle exceeds
// the critical angle (glass-to-air, steep angle).
func TestRefractTotalInternalReflection(t *testing.T) {
	// A near-parallel ray hitting a surface from glass side at a steep angle
	// triggers total internal reflection (discriminant <= 0).
	v := glm.Vec3{0.99, -0.14, 0} // nearly horizontal
	n := glm.Vec3{0, 1, 0}
	niOverNt := float32(1.5) // glass → air (> 1)
	_, ok := Refract(&v, &n, niOverNt)
	if ok {
		t.Error("Expected Refract to fail (ok=false) for total internal reflection")
	}
}

// TestRefractStraightOn verifies that a straight-on ray (v parallel to -n)
// refracted with niOverNt==1 passes through unchanged.
func TestRefractStraightOn(t *testing.T) {
	v := glm.Vec3{0, -1, 0}
	n := glm.Vec3{0, 1, 0}
	refracted, ok := Refract(&v, &n, 1.0) // same medium
	if !ok {
		t.Fatal("Expected Refract to succeed with niOverNt=1")
	}
	// With niOverNt=1 and straight-on incidence the refracted ray == v.Normalised()
	want := glm.Vec3{0, -1, 0}
	if !vec3Near(refracted, want, 1e-5) {
		t.Errorf("Refract straight-on: got %v, want %v", refracted, want)
	}
}

// ---------------------------------------------------------------------------
// RandomInUnitSphere
// ---------------------------------------------------------------------------

// TestRandomInUnitSphereInsideSphere verifies every returned point is strictly
// inside the unit sphere (length < 1) over many samples.
func TestRandomInUnitSphereInsideSphere(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	for i := 0; i < 1000; i++ {
		p := RandomInUnitSphere(r)
		if p.Len2() >= 1.0 {
			t.Errorf("Sample %d: point %v is outside unit sphere (len2=%v)", i, p, p.Len2())
		}
	}
}

// TestRandomInUnitSphereNotConstant verifies the function produces different
// values (i.e., it is not returning the same point every time).
func TestRandomInUnitSphereNotConstant(t *testing.T) {
	r := rand.New(rand.NewSource(0))
	first := RandomInUnitSphere(r)
	allSame := true
	for i := 0; i < 20; i++ {
		p := RandomInUnitSphere(r)
		if p != first {
			allSame = false
			break
		}
	}
	if allSame {
		t.Error("RandomInUnitSphere returned the same point every time")
	}
}

// TestRandomInUnitSphereDifferentSeeds verifies results differ across seeds.
func TestRandomInUnitSphereDifferentSeeds(t *testing.T) {
	r1 := rand.New(rand.NewSource(1))
	r2 := rand.New(rand.NewSource(2))
	p1 := RandomInUnitSphere(r1)
	p2 := RandomInUnitSphere(r2)
	if p1 == p2 {
		t.Error("Expected different random points for different seeds")
	}
}

// ---------------------------------------------------------------------------
// Schlick
// ---------------------------------------------------------------------------

// TestSchlickCosine1 verifies Schlick(1, any) == r0 (no boost at normal incidence).
func TestSchlickCosine1(t *testing.T) {
	refidx := float32(1.5)
	r0 := (1.0 - refidx) / (1.0 + refidx)
	r0 = r0 * r0
	got := Schlick(1.0, refidx)
	if !f32Near(got, r0, 1e-6) {
		t.Errorf("Schlick(1, %v) = %v; want %v", refidx, got, r0)
	}
}

// TestSchlickCosine0 verifies Schlick(0, any) == 1.0 (full reflectance at
// grazing incidence).
func TestSchlickCosine0(t *testing.T) {
	got := Schlick(0.0, 1.5)
	if !f32Near(got, 1.0, 1e-6) {
		t.Errorf("Schlick(0, 1.5) = %v; want 1.0", got)
	}
}

// TestSchlickRange verifies the result is always in [0, 1] for a range of
// cosine values.
func TestSchlickRange(t *testing.T) {
	refidx := float32(1.5)
	for _, cosine := range []float32{0, 0.1, 0.25, 0.5, 0.75, 0.9, 1.0} {
		got := Schlick(cosine, refidx)
		if got < 0 || got > 1 {
			t.Errorf("Schlick(%v, %v) = %v; want value in [0,1]", cosine, refidx, got)
		}
	}
}

// TestSchlickMonotone verifies that Schlick is monotonically decreasing in
// cosine (reflectance decreases as the angle becomes more head-on).
func TestSchlickMonotone(t *testing.T) {
	refidx := float32(1.5)
	prev := Schlick(0.0, refidx)
	for _, cosine := range []float32{0.2, 0.4, 0.6, 0.8, 1.0} {
		cur := Schlick(cosine, refidx)
		if cur > prev+1e-6 {
			t.Errorf("Schlick not monotone: Schlick(%v) = %v > Schlick(prev) = %v", cosine, cur, prev)
		}
		prev = cur
	}
}

// TestSchlickSameIndexFormula verifies that with identical refractive indices
// (refidx == 1.0), r0 == 0, so Schlick(cosine, 1.0) == (1-cosine)^5.
func TestSchlickSameIndexFormula(t *testing.T) {
	refidx := float32(1.0)
	for _, cosine := range []float32{0, 0.5, 1.0} {
		got := Schlick(cosine, refidx)
		// r0 = 0 → Schlick = 0 + 1*pow(1-cosine,5) = (1-cosine)^5
		want := float32(math.Pow(float64(1.0-cosine), 5))
		if !f32Near(got, want, 1e-6) {
			t.Errorf("Schlick(%v, 1.0) = %v; want %v", cosine, got, want)
		}
	}
}

// ---------------------------------------------------------------------------
// ColorRGBAFromVec3
// ---------------------------------------------------------------------------

// TestColorRGBAFromVec3Black verifies that the zero vector maps to black.
func TestColorRGBAFromVec3Black(t *testing.T) {
	got := ColorRGBAFromVec3(glm.Vec3{0, 0, 0})
	want := color.RGBA{R: 0, G: 0, B: 0, A: 255}
	if got != want {
		t.Errorf("ColorRGBAFromVec3({0,0,0}) = %v; want %v", got, want)
	}
}

// TestColorRGBAFromVec3White verifies that {1,1,1} maps to full white (255,255,255).
func TestColorRGBAFromVec3White(t *testing.T) {
	got := ColorRGBAFromVec3(glm.Vec3{1, 1, 1})
	// 255.99 * 1.0 = 255 (integer truncation)
	want := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	if got != want {
		t.Errorf("ColorRGBAFromVec3({1,1,1}) = %v; want %v", got, want)
	}
}

// TestColorRGBAFromVec3AlwaysFullAlpha verifies alpha is always 255.
func TestColorRGBAFromVec3AlwaysFullAlpha(t *testing.T) {
	cases := []glm.Vec3{{0, 0, 0}, {0.5, 0.5, 0.5}, {1, 1, 1}}
	for _, v := range cases {
		got := ColorRGBAFromVec3(v)
		if got.A != 255 {
			t.Errorf("ColorRGBAFromVec3(%v).A = %d; want 255", v, got.A)
		}
	}
}

// TestColorRGBAFromVec3Channels verifies per-channel mapping for a known value.
// For channel value 0.5: uint8(255.99 * 0.5) == uint8(127.995) == 127.
func TestColorRGBAFromVec3Channels(t *testing.T) {
	v := glm.Vec3{1.0, 0.5, 0.0}
	got := ColorRGBAFromVec3(v)
	// Mirror the production cast: uint8(255.99 * channel)
	wantR := uint8(float32(255.99) * v[0])
	wantG := uint8(float32(255.99) * v[1])
	wantB := uint8(float32(255.99) * v[2])
	if got.R != wantR || got.G != wantG || got.B != wantB {
		t.Errorf("ColorRGBAFromVec3(%v) = {R:%d G:%d B:%d}; want {R:%d G:%d B:%d}",
			v, got.R, got.G, got.B, wantR, wantG, wantB)
	}
}

// TestColorRGBAFromVec3IndependentChannels verifies each channel is mapped
// independently (cross-channel bleed would show up here).
func TestColorRGBAFromVec3IndependentChannels(t *testing.T) {
	red   := ColorRGBAFromVec3(glm.Vec3{1, 0, 0})
	green := ColorRGBAFromVec3(glm.Vec3{0, 1, 0})
	blue  := ColorRGBAFromVec3(glm.Vec3{0, 0, 1})

	if red.R != 255 || red.G != 0 || red.B != 0 {
		t.Errorf("Pure red: got %v", red)
	}
	if green.R != 0 || green.G != 255 || green.B != 0 {
		t.Errorf("Pure green: got %v", green)
	}
	if blue.R != 0 || blue.G != 0 || blue.B != 255 {
		t.Errorf("Pure blue: got %v", blue)
	}
}
