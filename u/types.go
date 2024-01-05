package u

type Region [2]float64

type Subtitle struct {
	Region     Region
	Transcript string
}

type RegionWithIndex struct {
	Region Region
	Index  int
}
