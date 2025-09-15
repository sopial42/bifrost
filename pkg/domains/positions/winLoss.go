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
	// negative WinLossRatio
	// is very dangerous
	R21 WinLossRatio = 2 / 1
	R31 WinLossRatio = 3 / 1
	R41 WinLossRatio = 4 / 1
	R51 WinLossRatio = 5 / 1
)

var AvailableWinLossRatios = []WinLossRatio{
	R11,
	R12,
	R13,
	R14,
	R15,
	R21,
	R31,
	R41,
	R51,
}

func (w WinLossRatio) ComputeStoploss(price float64, tp float64) float64 {
	win := tp - price
	loss := win * float64(w)
	sl := price - loss
	return sl
}
