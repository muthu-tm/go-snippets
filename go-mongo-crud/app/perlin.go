// Package perlin provides coherent noise function over 1, 2 or 3 dimensions
// This code is go adaptagion based on C implementation that can be found here:
// http://git.gnome.org/browse/gegl/tree/operations/common/perlin/perlin.c
// (original copyright Ken Perlin)
package app

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

var wg sync.WaitGroup

// General constants
const (
	B  = 0x100
	N  = 0x1000
	BM = 0xff
)

// Perlin is the noise generator
type Perlin struct {
	alpha float64
	beta  float64
	n     int

	p  [B + B + 2]int
	g3 [B + B + 2][3]float64
	g2 [B + B + 2][2]float64
	g1 [B + B + 2]float64
}

// NewPerlin creates new Perlin noise generator
// In what follows "alpha" is the weight when the sum is formed.
// Typically it is 2, As this approaches 1 the function is noisier.
// "beta" is the harmonic scaling/spacing, typically 2, n is the
// number of iterations and seed is the math.rand seed value to use
func NewPerlin(alpha, beta float64, n int, seed int64) *Perlin {
	return NewPerlinRandSource(alpha, beta, n, rand.NewSource(seed))
}

// NewPerlinRandSource creates new Perlin noise generator
// In what follows "alpha" is the weight when the sum is formed.
// Typically it is 2, As this approaches 1 the function is noisier.
// "beta" is the harmonic scaling/spacing, typically 2, n is the
// number of iterations and source is source of pseudo-random int64 values
func NewPerlinRandSource(alpha, beta float64, n int, source rand.Source) *Perlin {
	var p Perlin
	var i int

	p.alpha = alpha
	p.beta = beta
	p.n = n

	r := rand.New(source)

	for i = 0; i < B; i++ {
		p.p[i] = i
		p.g1[i] = float64((r.Int()%(B+B))-B) / B

		for j := 0; j < 2; j++ {
			p.g2[i][j] = float64((r.Int()%(B+B))-B) / B
		}

		normalize2(&p.g2[i])
	}

	for ; i > 0; i-- {
		k := p.p[i]
		j := r.Int() % B
		p.p[i] = p.p[j]
		p.p[j] = k
	}

	for i := 0; i < B+2; i++ {
		p.p[B+i] = p.p[i]
		p.g1[B+i] = p.g1[i]
		for j := 0; j < 2; j++ {
			p.g2[B+i][j] = p.g2[i][j]
		}
		for j := 0; j < 3; j++ {
			p.g3[B+i][j] = p.g3[i][j]
		}
	}

	return &p
}

func normalize2(v *[2]float64) {
	s := math.Sqrt(v[0]*v[0] + v[1]*v[1])
	v[0] = v[0] / s
	v[1] = v[1] / s
}

func at2(rx, ry float64, q [2]float64) float64 {
	return rx*q[0] + ry*q[1]
}

func at3(rx, ry, rz float64, q [3]float64) float64 {
	return rx*q[0] + ry*q[1] + rz*q[2]
}

func sCurve(t float64) float64 {
	return t * t * (3. - 2.*t)
}

func lerp(t, a, b float64) float64 {
	return a + t*(b-a)
}

func (p *Perlin) noise2(vec [2]float64) float64 {

	t := vec[0] + N
	bx0 := int(t) & BM
	bx1 := (bx0 + 1) & BM
	rx0 := t - float64(int(t))
	rx1 := rx0 - 1.

	t = vec[1] + N
	by0 := int(t) & BM
	by1 := (by0 + 1) & BM
	ry0 := t - float64(int(t))
	ry1 := ry0 - 1.

	i := p.p[bx0]
	j := p.p[bx1]

	b00 := p.p[i+by0]
	b10 := p.p[j+by0]
	b01 := p.p[i+by1]
	b11 := p.p[j+by1]

	sx := sCurve(rx0)
	sy := sCurve(ry0)

	q := p.g2[b00]
	u := at2(rx0, ry0, q)
	q = p.g2[b10]
	v := at2(rx1, ry0, q)
	a := lerp(sx, u, v)

	q = p.g2[b01]
	u = at2(rx0, ry1, q)
	q = p.g2[b11]
	v = at2(rx1, ry1, q)
	b := lerp(sx, u, v)

	return lerp(sy, a, b)
}

// Noise2D Generates 2-dimensional Perlin Noise value
func (p *Perlin) Noise2D(x, y float64, c chan float64) {
	var scale float64 = 1
	var sum float64
	var px [2]float64

	px[0] = x
	px[1] = y

	for i := 0; i < p.n; i++ {
		val := p.noise2(px)
		sum += val / scale
		scale *= p.alpha
		px[0] *= p.beta
		px[1] *= p.beta
	}

	fmt.Printf("%0.0f\t%0.0f\t%0.4f\n", x*10, y*10, sum)
	c <- sum
}

func (p *Perlin) GetNoise2D(x int, y int, c chan float64) (noises []float64) {
	noises = make([]float64, 0)
	for i := 0.; i <= float64(x); i++ {
		for j := 0.; j <= float64(y); j++ {
			go p.Noise2D(i/10, j/10, c)
			wg.Add(1)
			go func() {
				defer wg.Done()
				noises = append(noises, <-c)
			}()
		}
	}

	wg.Wait()

	return
}

func GetPerlinNoise2D(c *gin.Context) {
	noiseChan := make(chan float64)
	x, _ := strconv.Atoi(c.PostForm("x"))
	y, _ := strconv.Atoi(c.PostForm("y"))
	p := NewPerlin(2, 2, 2, 12345)

	noises := p.GetNoise2D(x, y, noiseChan)
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": noises})
}
