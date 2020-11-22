package maps

func GenMap(l int) *[][]int {
	gmap := make([][]int, l)

	// generating full filled hex map
	for i := 0; i < l; i++ {
		gmap[i] = make([]int, l*2)
		for j := 0; j < l*2; j++ {
			if i%2 == 0 && j%2 == 0 {
				gmap[i][j] = 0
			} else if i%2 != 0 && j%2 != 0 {
				gmap[i][j] = 0
			} else {
				gmap[i][j] = 1
			}
		}
	}

	return &gmap
}
