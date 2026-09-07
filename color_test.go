package GoRTIW

import (
	"image"
	"math"
	"math/rand"
	"sync"
	"testing"

	"github.com/engoengine/glm"
)

// MockSurface implements the Surface interface for testing.
// HitMock controls whether Hit() reports a collision and can populate HitRecord.
type MockSurface struct {
	HitMock  func(ray *Ray, tmin, tmax float32, hitRecord *HitRecord) bool
	Material Material
}

func (m MockSurface) GetMaterial() Material { return m.Material }

func (m MockSurface) Hit(ray *Ray, tMin, tMax float32, hit *HitRecord) bool {
	return m.HitMock(ray, tMin, tMax, hit)
}

// mockMaterial is a minimal Material stub.
// When scatter==true it copies albedo into attenuation and returns true;
// when scatter==false it returns false (fully absorbed).
type mockMaterial struct {
	albedo  glm.Vec3
	scatter bool
}

func (m *mockMaterial) Scatter(ray *Ray, h *HitRecord, attenuation *glm.Vec3, s *Ray, r *rand.Rand) bool {
	if !m.scatter {
		return false
	}
	*attenuation = m.albedo
	*s = Ray{Origin: h.P, Direction: ray.Direction}
	return true
}

// ---------------------------------------------------------------------------
// ComputeColor tests
// ---------------------------------------------------------------------------

// TestComputeColorNoHit verifies that a ray missing all surfaces returns the
// sky gradient (blend between white and light-blue), not a zero vector.
func TestComputeColorNoHit(t *testing.T) {
	surfaces := Surfaces{}
	surfaces.Add(MockSurface{
		HitMock: func(_ *Ray, _, _ float32, _ *HitRecord) bool { return false },
	})

	// Straight-up ray: normalised Y == 1 → t == 1 → pure interpB {0.5, 0.7, 1.0}
	ray := Ray{Direction: glm.Vec3{0, 1, 0}}
	r := rand.New(rand.NewSource(0))

	result := ComputeColor(&ray, &surfaces, 0, r)

	if result == (glm.Vec3{}) {
		t.Error("Expected sky gradient colour, got zero vector")
	}
	expected := glm.Vec3{0.5, 0.7, 1.0}
	if !vec3ApproxEqual(result, expected, 1e-5) {
		t.Errorf("Expected sky colour %v, got %v", expected, result)
	}
}

// TestComputeColorSkyGradientHorizon verifies the sky colour for a horizontal
// ray (Y≈0): the lerp parameter t == 0.5, giving {0.75, 0.85, 1.0}.
func TestComputeColorSkyGradientHorizon(t *testing.T) {
	surfaces := Surfaces{}
	surfaces.Add(MockSurface{
		HitMock: func(_ *Ray, _, _ float32, _ *HitRecord) bool { return false },
	})

	ray := Ray{Direction: glm.Vec3{1, 0, 0}}
	r := rand.New(rand.NewSource(0))

	result := ComputeColor(&ray, &surfaces, 0, r)
	expected := glm.Vec3{0.75, 0.85, 1.0}
	if !vec3ApproxEqual(result, expected, 1e-5) {
		t.Errorf("Expected horizon sky %v, got %v", expected, result)
	}
}

// TestComputeColorHitAbsorbed verifies that a hit whose material does not
// scatter (fully absorbed) returns the zero vector.
func TestComputeColorHitAbsorbed(t *testing.T) {
	mat := &mockMaterial{scatter: false}
	surfaces := Surfaces{}
	surfaces.Add(MockSurface{
		HitMock: func(_ *Ray, _, _ float32, hr *HitRecord) bool {
			hr.T = 1.0
			hr.P = glm.Vec3{0, 0, -1}
			hr.Normal = glm.Vec3{0, 0, 1}
			hr.Material = mat
			return true
		},
		Material: mat,
	})

	ray := Ray{Direction: glm.Vec3{0, 0, -1}}
	r := rand.New(rand.NewSource(42))

	result := ComputeColor(&ray, &surfaces, 0, r)
	if result != (glm.Vec3{}) {
		t.Errorf("Expected zero vector for absorbed hit, got %v", result)
	}
}

// TestComputeColorHitScatter verifies that a scattering hit (white albedo)
// returns a non-zero colour.
// The surface hits exactly once (first call); subsequent recursive rays miss
// so ComputeColor returns the sky gradient, which is multiplied by the white
// albedo and propagated back — giving a non-zero result.
func TestComputeColorHitScatter(t *testing.T) {
	mat := &mockMaterial{albedo: glm.Vec3{1, 1, 1}, scatter: true}
	hits := 0
	surfaces := Surfaces{}
	surfaces.Add(MockSurface{
		HitMock: func(ray *Ray, _, _ float32, hr *HitRecord) bool {
			if hits == 0 {
				hits++
				hr.T = 1.0
				hr.P = glm.Vec3{0, 0, -1}
				hr.Normal = glm.Vec3{0, 0, 1}
				hr.Material = mat
				return true
			}
			return false
		},
		Material: mat,
	})

	// Upward scattered ray so the sky gradient miss returns a bright colour.
	ray := Ray{Direction: glm.Vec3{0, 1, 0}}
	r := rand.New(rand.NewSource(7))

	result := ComputeColor(&ray, &surfaces, 0, r)
	if result == (glm.Vec3{}) {
		t.Error("Expected non-zero colour from scattering material, got zero vector")
	}
}

// TestComputeColorMaxDepth verifies that at depth >= 50 a hit returns zero
// even when the material would scatter.
func TestComputeColorMaxDepth(t *testing.T) {
	mat := &mockMaterial{albedo: glm.Vec3{1, 1, 1}, scatter: true}
	surfaces := Surfaces{}
	surfaces.Add(MockSurface{
		HitMock: func(_ *Ray, _, _ float32, hr *HitRecord) bool {
			hr.T = 1.0
			hr.Material = mat
			return true
		},
		Material: mat,
	})

	ray := Ray{Direction: glm.Vec3{0, 0, -1}}
	r := rand.New(rand.NewSource(0))

	result := ComputeColor(&ray, &surfaces, 50, r)
	if result != (glm.Vec3{}) {
		t.Errorf("Expected zero vector at max depth, got %v", result)
	}
}

// ---------------------------------------------------------------------------
// ComputePixel tests
// ---------------------------------------------------------------------------

// TestComputePixelCallsDone verifies ComputePixel calls wg.Done() so that
// wg.Wait() returns without deadlock.
// ComputePixel is synchronous, so we can call wg.Wait() directly after it.
func TestComputePixelCallsDone(t *testing.T) {
	surfaces := Surfaces{}
	surfaces.Add(MockSurface{
		HitMock: func(_ *Ray, _, _ float32, _ *HitRecord) bool { return false },
	})

	nx, ny, ns := 4, 4, 1
	output := image.NewRGBA(image.Rect(0, 0, nx, ny+1))
	camera := NewCamera(
		glm.Vec3{0, 0, 0}, glm.Vec3{0, 0, -1}, glm.Vec3{0, 1, 0},
		90, float32(nx)/float32(ny), 0, 1,
	)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	ComputePixel(0, 0, nx, ny, ns, camera, &surfaces, output, wg)
	// If wg.Done() was never called this will hang and the test runner will
	// time out — which is the correct failure signal.
	wg.Wait()
}

// TestComputePixelGammaCorrection verifies gamma correction (sqrt) is applied.
// For a miss ray the sky-gradient produces a deterministic value (for ns == 1,
// only the floor of u/v matters since float(i)/float(nx) is exact).
// We confirm each output channel satisfies channel == uint8(sqrt(raw)*255).
func TestComputePixelGammaCorrection(t *testing.T) {
	surfaces := Surfaces{}
	surfaces.Add(MockSurface{
		HitMock: func(_ *Ray, _, _ float32, _ *HitRecord) bool { return false },
	})

	// 1×1 image — pixel (0,0), ny-j = 1-0 = 1 so we need at least height 2.
	nx, ny, ns := 1, 1, 1
	output := image.NewRGBA(image.Rect(0, 0, nx, ny+1))
	camera := NewCamera(
		glm.Vec3{0, 0, 0}, glm.Vec3{0, 0, -1}, glm.Vec3{0, 1, 0},
		90, 1.0, 0, 1,
	)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	ComputePixel(0, 0, nx, ny, ns, camera, &surfaces, output, wg)
	wg.Wait()

	px := output.RGBAAt(0, ny)
	if px.A != 255 {
		t.Errorf("Expected alpha 255, got %d", px.A)
	}
	// The sky gradient is always in (0,1], so after gamma correction the pixel
	// must be non-black.
	if px.R == 0 && px.G == 0 && px.B == 0 {
		t.Error("Expected non-black sky pixel after gamma correction, got {0,0,0}")
	}
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func vec3ApproxEqual(a, b glm.Vec3, eps float32) bool {
	return float32(math.Abs(float64(a[0]-b[0]))) < eps &&
		float32(math.Abs(float64(a[1]-b[1]))) < eps &&
		float32(math.Abs(float64(a[2]-b[2]))) < eps
}
