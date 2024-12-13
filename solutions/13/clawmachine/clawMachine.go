package clawmachine

import (
	"errors"
	"log/slog"
	"math"

	"gonum.org/v1/gonum/mat"
)

const (
	// Yes this is small... intentionally! For part 2, in which values climb to 10^12
	floatIntConversionEpsilon float64 = 0.001
)

var (
	unitConversionFix          *mat.VecDense = mat.NewVecDense(2, []float64{10000000000000, 10000000000000})
	buttonTokenCosts           *mat.VecDense = mat.NewVecDense(2, []float64{3, 1})
	ErrorFloatNotCastableToInt error         = errors.New("cannot cast float to int with significant loss of accuracy")
	ErrorIncomputable          error         = errors.New("cannot determine a button sequence that can reach the prize")
)

func floatIntConversion(x float64) (int, error) {
	intCast := int(math.Round(x))
	if math.Abs(x-float64(intCast)) > floatIntConversionEpsilon {
		return 0, ErrorFloatNotCastableToInt
	}
	return intCast, nil
}

type ClawMachine struct {
	coordinateIncrementMatrix *mat.Dense
	targetCoordinateVector    *mat.VecDense
}

func NewClawMachine(buttonAX, buttonAY, buttonBX, buttonBY, prizeX, prizeY float64) *ClawMachine {
	coordinateIncrementMatrix := mat.NewDense(2, 2, []float64{
		buttonAX, buttonBX, buttonAY, buttonBY,
	})
	targetCoordinateVector := mat.NewVecDense(2, []float64{
		prizeX, prizeY,
	})
	return &ClawMachine{
		coordinateIncrementMatrix,
		targetCoordinateVector,
	}
}

func (machine *ClawMachine) FixUnitConversion() {
	machine.targetCoordinateVector.AddVec(machine.targetCoordinateVector, unitConversionFix)
}

func (machine *ClawMachine) ComputeLowestTokenCost() (int, error) {
	var buttonPushes mat.VecDense
	err := buttonPushes.SolveVec(machine.coordinateIncrementMatrix, machine.targetCoordinateVector)
	if err != nil {
		return 0, ErrorIncomputable
	}

	buttonPushes.MulElemVec(&buttonPushes, buttonTokenCosts)
	buttonCostsRawData := buttonPushes.RawVector().Data
	slog.Debug("cost computed", "total cost", buttonCostsRawData)

	totalCostInt := 0
	for _, cost := range buttonCostsRawData {
		intCost, err := floatIntConversion(cost)
		if err != nil {
			return 0, ErrorFloatNotCastableToInt
		}
		totalCostInt += intCost
	}
	return totalCostInt, nil
}
