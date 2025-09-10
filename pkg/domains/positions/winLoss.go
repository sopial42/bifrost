package positions

type WinLossRatio float64

const (
	// positive WinLossRatio
	// should be used
	R11 WinLossRatio = 1.0 / 1
	R12 WinLossRatio = 1.0 / 2
	R13 WinLossRatio = 1.0 / 3
	R14 WinLossRatio = 1.0 / 4
	R15 WinLossRatio = 1.0 / 5
	R16 WinLossRatio = 1.0 / 6
	R17 WinLossRatio = 1.0 / 7
	R18 WinLossRatio = 1.0 / 8
	R19 WinLossRatio = 1.0 / 9
	// negative WinLossRatio
	// is very dangerous
	R21 WinLossRatio = 2 / 1
	R31 WinLossRatio = 3 / 1
	R41 WinLossRatio = 4 / 1
	R51 WinLossRatio = 5 / 1
	R61 WinLossRatio = 6 / 1
	R71 WinLossRatio = 7 / 1
	R81 WinLossRatio = 8 / 1
	R91 WinLossRatio = 9 / 1
)



func (w WinLossRatio) ComputeStoploss(price float64, tp float64) float64 {
	win := tp - price
	loss := win * float64(w)
	sl := price - loss
	return sl
}
