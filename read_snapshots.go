package main

import (
	"log"

	"gonum.org/v1/hdf5"
)

type Vec2 struct {
	x float64
	y float64
}

func read_snapshot(filename string) (uint, float64, []Vec2, []float64) {
	// HDF5 파일 열기
	file, err := hdf5.OpenFile(filename, hdf5.F_ACC_RDONLY)
	if err != nil {
		log.Fatalf("파일을 열 수 없습니다: %v", err)
	}
	defer file.Close()

	world, _ := file.OpenGroup("world")
	defer world.Close()

	tAtt, _ := world.OpenAttribute("t")
	defer tAtt.Close()
	var t float64
	tAtt.Read(&t, hdf5.T_NATIVE_DOUBLE)

	atoms, _ := world.OpenGroup("atoms")
	defer atoms.Close()

	pos_data, _ := atoms.OpenDataset("pos")
	defer pos_data.Close()
	dataspace := pos_data.Space()
	dims, _, _ := dataspace.SimpleExtentDims()
	N := dims[0]
	pos := make([]Vec2, N)
	pos_data.Read(&pos)

	pol_data, _ := atoms.OpenDataset("polarity")
	defer pol_data.Close()
	pol := make([]float64, N)
	pol_data.Read(&pol)

	return N, t, pos, pol
}
