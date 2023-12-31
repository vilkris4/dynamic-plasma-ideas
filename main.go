package main

import (
	"fmt"
	"math"
	"math/big"
)

const (
	TargetBasePlasmaElasticity = 2

	MaxBasePlasmaInMomentum          = 4200000                                              // 200 * AccountBlockBasePlasma
	TargetBasePlasmaInMomentum       = MaxBasePlasmaInMomentum / TargetBasePlasmaElasticity // 50% of MaxBasePlasmaInMomentum
	TargetPoWPlasmaPercentInMomentum = 20                                                   // Target is that 20% of plasma is generated by PoW when base plasma is over or equal to TargetBasePlasmaInMomentum.

	MaxFusionRechargeRate = 210000 // Equals 100 QSR/confirmation. This value should be dependent on maximum network throughput.

	BasePriceChangeDenominator           = 8 // Denominator to limit the maximimum target offset change per momentum.
	DifficultyPerPlasmaChangeDenominator = 8 // Denominator to limit the maximimum target offset change per momentum.

	PlasmaRechargeRateChangeMultiplier = 2 // Maximum rate at which the plasma recharge rate can change. Decreases at -50%, increases at 100%.

	AccountPlasmaRechargeRateDenominator = 210000

	AccountBlockBasePlasma      = 21000
	NumFusionUnitsForBasePlasma = 10
	PlasmaPerFusionUnit         = AccountBlockBasePlasma / NumFusionUnitsForBasePlasma

	MinDifficultyPerPlasma = 500 // Minimum PoW difficulty value

	MinBasePrice         = 1000 // Minimum base price for plasma
	BasePriceDenominator = 1000
)

type Momentum struct {
	BasePrice           uint64 // The minimum plasma price of the previous momentum
	BasePlasma          uint64 // The total base plasma used in the momentum
	UsedPlasma          uint64 // The total used plasma in the momentum => Minimum is: BasePlasma * previousMomentum.BasePrice / BasePriceDenominator
	Difficulty          uint64 // The total PoW difficulty used in the momentum
	DifficultyPerPlasma uint64 // The PoW difficulty per plasma value for the momentum
	FusionRechargeRate  uint64 // The fused plasma recharge rate for the momentum
}

func main() {
	txCount := uint64(100)                                              // The transaction count of the frontier momentum
	basePlasma := txCount * AccountBlockBasePlasma                      // Consumed base plasma by momentum, assuming all TXs require AccountBlockBasePlasma
	basePrice := uint64(1000)                                           // The minimum plasma price of the frontier momentum. Used to calculate the minimum plasma price for the next momentum.
	usedPlasma := basePlasma * uint64(basePrice) / BasePriceDenominator // basePlasma is the minimum the usedPlasma will be for the momentum
	poWPlasmaPercentInMomentum := uint64(20)                            // The percentage of total plasma generated by PoW in the frontier momentum
	difficultyPerPlasma := uint64(1000)                                 // The PoW difficulty per plasma of the frontier momentum, adjusts dynamically
	fusionRechargeRate := uint64(2100)                                  // The fused plasma recharge rate of the frontier momentum, adjusts dynamically

	if basePrice < MinBasePrice {
		fmt.Println("Base price is too small")
		return
	}

	if basePlasma > MaxBasePlasmaInMomentum {
		fmt.Println("Used base plasma in momentum is too high")
		return
	}

	if difficultyPerPlasma < MinDifficultyPerPlasma {
		fmt.Println("PoW difficulty per plasma is too small")
		return
	}

	if fusionRechargeRate > MaxFusionRechargeRate {
		fmt.Println("Fused plasma recharge rage is too high")
		return
	}

	// Construct frontier momentum
	frontier := Momentum{
		BasePrice:           basePrice,
		BasePlasma:          basePlasma,
		UsedPlasma:          usedPlasma,
		Difficulty:          usedPlasma * difficultyPerPlasma * poWPlasmaPercentInMomentum / 100,
		DifficultyPerPlasma: difficultyPerPlasma,
		FusionRechargeRate:  fusionRechargeRate,
	}

	// Attack simulation
	attackerQsr := uint64(1000000)
	attackerQsrPerAddress := uint64(50)                       // The amount of QSR an attacker has in one address
	attackTargetQsrCutoff := uint64(50)                       // The target QSR amount the attacker wants to cutoff from using the network
	attackerTxs := int64(attackerQsr / attackTargetQsrCutoff) // QSR is split evenly into separate transactions
	attackerFilledMomentums := uint64(0)
	averageRechargeRate := uint64(0)

	rounds := 5000
	for i := 0; i < rounds; i++ {
		fmt.Println("------------------------------------------------------------------------------------------------")
		fmt.Println(fmt.Sprintf("Round #%d", i+1))
		fmt.Println("------")

		newDifficultyPerPlasma := GetPoWDifficultyPerPlasma(frontier)
		fmt.Println("New DifficultyPerPlasma is", newDifficultyPerPlasma)
		fmt.Println(fmt.Sprintf("DifficultyPerPlasma change is %.1f%%", (float64(newDifficultyPerPlasma)/float64(frontier.DifficultyPerPlasma)-1)*100))
		fmt.Println("------")

		newBasePrice := GetBasePrice(frontier)
		fmt.Println("New BasePrice is", newBasePrice)
		fmt.Println(fmt.Sprintf("BasePrice change is %.1f%%", (float64(newBasePrice)/float64(frontier.BasePrice)-1)*100))
		fmt.Println(fmt.Sprintf("Fused QSR needed to meet account block base plasma is %.5f QSR", (float64(AccountBlockBasePlasma)*float64(newBasePrice)/BasePriceDenominator)/float64(PlasmaPerFusionUnit)))
		fmt.Println("------")

		targetRechargeRate := GetFusionRechargeRateTarget(frontier)
		fmt.Println(fmt.Sprintf("New FusionRechargeRate target is %d (%.2f QSR/confirmation)", targetRechargeRate, float64(targetRechargeRate)/float64(PlasmaPerFusionUnit)))

		newFusionRechargeRate := GetFusionRechargeRate(frontier, targetRechargeRate)
		fmt.Println(fmt.Sprintf("New FusionRechargeRate is %d (%.2f QSR/confirmation)", newFusionRechargeRate, float64(newFusionRechargeRate)/float64(PlasmaPerFusionUnit)))
		fmt.Println(fmt.Sprintf("FusionRechargeRate change is %.1f%%", (float64(newFusionRechargeRate)/float64(frontier.FusionRechargeRate)-1)*100))
		fmt.Println("------------------------------------------------------------------------------------------------")

		// Simulate tx count change
		txCount := uint64(200)

		// Simulate attack
		minQsrRequired := (float64(AccountBlockBasePlasma) * float64(newBasePrice) / BasePriceDenominator) / float64(PlasmaPerFusionUnit)

		if minQsrRequired < float64(attackTargetQsrCutoff) {
			attackerTxs -= int64(txCount)
			averageRechargeRate = averageRechargeRate + uint64(float64(newFusionRechargeRate)*GetAddressRechargeRateMultiplier(attackerQsrPerAddress*PlasmaPerFusionUnit))
			fmt.Println("attacker recharge multiplier", GetAddressRechargeRateMultiplier(attackerQsrPerAddress*PlasmaPerFusionUnit))
			fmt.Println(fmt.Sprintf("attacker recharge rate %d (%.2f QSR/confirmation)", uint64(float64(newFusionRechargeRate)*GetAddressRechargeRateMultiplier(attackerQsrPerAddress*PlasmaPerFusionUnit)), (float64(newFusionRechargeRate) * GetAddressRechargeRateMultiplier(attackerQsrPerAddress*PlasmaPerFusionUnit) / float64(PlasmaPerFusionUnit))))

			attackerFilledMomentums++
		} else {
			txCount = uint64(0) // Attacker has cutoff everyone else by raising the base price.
		}
		fmt.Println("Attacker TXs left in mempool", attackerTxs)
		fmt.Println("Attacker overall recharge rate", attackerTxs)

		if attackerTxs <= 0 && minQsrRequired < float64(attackTargetQsrCutoff) {
			fmt.Println("------")
			fmt.Println("Attack result")
			fmt.Println("------")
			rechargeRateInQsr := float64((averageRechargeRate / uint64(attackerFilledMomentums))) / float64(PlasmaPerFusionUnit)
			momentumsUntilRecharged := float64(attackerQsrPerAddress) / rechargeRateInQsr
			blockingTime := float64(i*10) / 60
			attackCooldownTime := momentumsUntilRecharged * 10 / 60
			fmt.Println(fmt.Sprintf("Average FusionRechargeRate is %d (%.2f QSR/confirmation)", (averageRechargeRate / uint64(attackerFilledMomentums)), rechargeRateInQsr))
			fmt.Println(fmt.Sprintf("Able to block users with less than %d fused QSR for %.2f minutes", attackTargetQsrCutoff, blockingTime))
			fmt.Println(fmt.Sprintf("Attack can be repeated in %.2f minutes", attackCooldownTime))
			fmt.Println(fmt.Sprintf("Bandwidth ratio %.2f", blockingTime/attackCooldownTime))
			fmt.Println("------------------------------------------------------------------------------------------------")
			break
		}

		// Simulate changes
		basePlasma = txCount * AccountBlockBasePlasma
		usedPlasma = basePlasma
		poWPlasmaPercentInMomentum = uint64(0)

		frontier = Momentum{
			BasePrice:           newBasePrice,
			BasePlasma:          basePlasma,
			UsedPlasma:          usedPlasma,
			Difficulty:          usedPlasma * newDifficultyPerPlasma * poWPlasmaPercentInMomentum / 100,
			DifficultyPerPlasma: newDifficultyPerPlasma,
			FusionRechargeRate:  newFusionRechargeRate,
		}
	}
}

func GetAddressRechargeRateMultiplier(addressFusedPlasma uint64) float64 {
	return float64(addressFusedPlasma) / float64(AccountPlasmaRechargeRateDenominator)
}

// Get the fused plasma recharge rate target for a given momentum.
// Recharge rate is dependent on momentum fullness.
// Get rate multiplier using exponential function => f(x) = MaxFusionRechargeRate * 10000^-x
// Assuming MaxRechargeRatePerConfirmation == 210,000 plasma/confirmation:
// At 0% full => 210,000 plasma/confirmation => 100 QSR/confirmation
// At 50% full => 2,100 plasma/confirmation => 1 QSR/confirmation
// At 100% full => 21 plasma/confirmation => 0.01 QSR/confirmation
func GetFusionRechargeRateTarget(frontier Momentum) uint64 {
	fullness := float64(frontier.BasePlasma) / float64(MaxBasePlasmaInMomentum)
	multiplier := math.Pow(10000, -fullness)
	return uint64(float64(MaxFusionRechargeRate) * multiplier)
}

// Get the fusion recharge rate based on current target.
func GetFusionRechargeRate(frontier Momentum, target uint64) uint64 {
	if frontier.FusionRechargeRate == target {
		return frontier.FusionRechargeRate
	}
	rate := uint64(0)
	if frontier.FusionRechargeRate > target {
		rate = frontier.FusionRechargeRate / PlasmaRechargeRateChangeMultiplier
		if rate < target {
			rate = target
		}
	} else {
		rate = frontier.FusionRechargeRate * PlasmaRechargeRateChangeMultiplier
		if rate > target {
			rate = target
		}
	}
	return rate
}

// Get new base price for plasma based on the previous momentum
func GetBasePrice(frontier Momentum) uint64 {
	basePrice := GetTargetOffsetMultiplier(TargetBasePlasmaInMomentum, frontier.BasePlasma, frontier.BasePrice, BasePriceChangeDenominator)
	if basePrice < MinBasePrice {
		basePrice = MinBasePrice
	}
	return basePrice
}

// Get new PoW difficulty per plasma based on the previous momentum
func GetPoWDifficultyPerPlasma(frontier Momentum) uint64 {
	targetPlasma := frontier.UsedPlasma
	if frontier.BasePlasma >= TargetBasePlasmaInMomentum {
		targetPlasma = frontier.UsedPlasma * TargetPoWPlasmaPercentInMomentum / 100
	}
	difficultyPerPlasma := GetTargetOffsetMultiplier(targetPlasma, frontier.Difficulty/frontier.DifficultyPerPlasma, frontier.DifficultyPerPlasma, DifficultyPerPlasmaChangeDenominator)
	if difficultyPerPlasma < MinDifficultyPerPlasma {
		difficultyPerPlasma = MinDifficultyPerPlasma
	}
	return difficultyPerPlasma
}

// Get the new value of a multiplier depending on how far away the actual value is from the target value.
// m = m_c + (m_c * (A - T)) / (T * c)
func GetTargetOffsetMultiplier(target uint64, actual uint64, currentMultiplier uint64, changeDenominator uint64) uint64 {
	if actual == target {
		return currentMultiplier
	}
	bigCurrent := new(big.Int).SetUint64(currentMultiplier)
	bigActual := new(big.Int).SetUint64(actual)
	bigTarget := new(big.Int).SetUint64(target)
	distance := bigActual.Sub(bigActual, bigTarget)
	numerator := new(big.Int).Mul(bigCurrent, distance)
	denominator := new(big.Int).Mul(bigTarget, new(big.Int).SetUint64(changeDenominator))
	delta := new(big.Int).Div(numerator, denominator)
	multiplier := bigCurrent.Add(bigCurrent, delta)
	if multiplier.Cmp(big.NewInt(0)) == +1 {
		return multiplier.Uint64()
	} else {
		return uint64(0)
	}
}
