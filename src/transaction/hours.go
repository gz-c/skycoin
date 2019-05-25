package transaction

import (
	"errors"
	"math/big"
	"sort"

	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/util/fee"
	"github.com/skycoin/skycoin/src/util/mathutil"
)

// DistributeSpendHours calculates how many coin hours to transfer to the change address and how
// many to transfer to each of the other destination addresses.
// Input hours are split by BurnFactor (rounded down) to meet the fee requirement.
// The remaining hours are split in half, one half goes to the change address
// and the other half goes to the destination addresses.
// If the remaining hours are an odd number, the change address gets the extra hour.
// If the amount assigned to the destination addresses is not perfectly divisible by the
// number of destination addresses, the extra hours are distributed to some of these addresses.
// Returns the number of hours to send to the change address,
// an array of length nAddrs with the hours to give to each destination address,
// and a sum of these values.
func DistributeSpendHours(inputHours, nAddrs uint64, haveChange bool) (uint64, []uint64, uint64) {
	feeHours := fee.RequiredFee(inputHours, params.UserVerifyTxn.BurnFactor)
	remainingHours := inputHours - feeHours

	var changeHours uint64
	if haveChange {
		// Split the remaining hours between the change output and the other outputs
		changeHours = remainingHours / 2

		// If remainingHours is an odd number, give the extra hour to the change output
		if remainingHours%2 == 1 {
			changeHours++
		}
	}

	// Distribute the remaining hours equally amongst the destination outputs
	remainingAddrHours := remainingHours - changeHours
	addrHoursShare := remainingAddrHours / nAddrs

	// Due to integer division, extra coin hours might remain after dividing by len(toAddrs)
	// Allocate these extra hours to the toAddrs
	addrHours := make([]uint64, nAddrs)
	for i := range addrHours {
		addrHours[i] = addrHoursShare
	}

	extraHours := remainingAddrHours - (addrHoursShare * nAddrs)
	i := 0
	for extraHours > 0 {
		addrHours[i] = addrHours[i] + 1
		i++
		extraHours--
	}

	// Assert that the hour calculation is correct
	var spendHours uint64
	for _, h := range addrHours {
		spendHours += h
	}
	spendHours += changeHours
	if spendHours != remainingHours {
		logger.Panicf("spendHours != remainingHours (%d != %d), calculation error", spendHours, remainingHours)
	}

	return changeHours, addrHours, spendHours
}

// unsortableHighToLowUint64 allows us to sort and unsort a []uint64.
type unsortableHighToLowUint64 struct {
	items   []uint64
	indices []int
}

func newUnsortableHighToLowUint64(items []uint64) *unsortableHighToLowUint64 {
	x := &unsortableHighToLowUint64{
		items:   make([]uint64, len(items)),
		indices: make([]int, len(items)),
	}

	copy(x.items, items)
	for i := 0; i < len(x.indices); i++ {
		x.indices[i] = i
	}

	return x
}

func (u *unsortableHighToLowUint64) Swap(i, j int) {
	u.items[i], u.items[j] = u.items[j], u.items[i]
	u.indices[i], u.indices[j] = u.indices[j], u.indices[i]
}

func (u *unsortableHighToLowUint64) Less(i, j int) bool {
	return u.items[i] > u.items[j]
}

func (u *unsortableHighToLowUint64) Len() int {
	return len(u.items)
}

// DistributeCoinHoursProportional distributes hours amongst coins proportional to the coins amount
func DistributeCoinHoursProportional(coins []uint64, hours uint64) ([]uint64, error) {
	if len(coins) == 0 {
		return nil, errors.New("DistributeCoinHoursProportional coins array must not be empty")
	}

	coinsInt := make([]*big.Int, len(coins))

	var total uint64
	for i, c := range coins {
		if c == 0 {
			return nil, errors.New("DistributeCoinHoursProportional coins array has a zero value")
		}

		var err error
		total, err = mathutil.AddUint64(total, c)
		if err != nil {
			return nil, err
		}

		cInt64, err := mathutil.Uint64ToInt64(c)
		if err != nil {
			return nil, err
		}

		coinsInt[i] = big.NewInt(cInt64)
	}

	addrHours := make([]uint64, len(coins))
	if hours >= uint64(len(coins)) {
		// When total hours > the number of outputs, assign each output 1 hour, then
		// assign the remaining hours proportionally.
		for i := 0; i < len(coins); i++ {
			addrHours[i] = 1
		}

		hours -= uint64(len(coins))
	} else if hours > 0 {
		// If there aren't enough hours to give every output 1 hour,
		// make sure the outputs with more coins get the 1 hour
		sortedCoins := newUnsortableHighToLowUint64(coins)
		sort.Stable(sortedCoins)

		for i := uint64(0); i < hours; i++ {
			addrHours[sortedCoins.indices[i]] = 1
		}

		hours = 0
	}

	if hours == 0 {
		return addrHours, nil
	}

	totalInt64, err := mathutil.Uint64ToInt64(total)
	if err != nil {
		return nil, err
	}
	totalInt := big.NewInt(totalInt64)

	hoursInt64, err := mathutil.Uint64ToInt64(hours)
	if err != nil {
		return nil, err
	}
	hoursInt := big.NewInt(hoursInt64)

	var assignedHours uint64
	for i, c := range coinsInt {
		// Scale the ratio of coins to total coins proportionally by calculating
		// (coins * totalHours) / totalCoins
		// The remainder is truncated, remaining hours are appended after this
		num := &big.Int{}
		num.Mul(c, hoursInt)

		fracInt := big.Int{}
		fracInt.Div(num, totalInt)

		if !fracInt.IsUint64() {
			return nil, errors.New("DistributeCoinHoursProportional calculated fractional hours is not representable as a uint64")
		}

		fracHours := fracInt.Uint64()
		addrHours[i], err = mathutil.AddUint64(addrHours[i], fracHours)
		if err != nil {
			return nil, err
		}

		assignedHours, err = mathutil.AddUint64(assignedHours, fracHours)
		if err != nil {
			return nil, err
		}
	}

	if hours < assignedHours {
		return nil, errors.New("DistributeCoinHoursProportional assigned hours exceeding input hours, this is a bug")
	}

	remainingHours := hours - assignedHours

	if remainingHours > uint64(len(coins)) {
		return nil, errors.New("DistributeCoinHoursProportional remaining hours exceed len(coins), this is a bug")
	}

	// For remaining hours lost due to fractional cutoff when scaling,
	// first provide at least 1 coin hour to coins that were assigned 0.
	i := 0
	for remainingHours > 0 && i < len(coins) {
		if addrHours[i] == 0 {
			addrHours[i] = 1
			remainingHours--
		}
		i++
	}

	// Then, assign the extra coin hours
	i = 0
	for remainingHours > 0 {
		addrHours[i] = addrHours[i] + 1
		remainingHours--
		i++
	}

	return addrHours, nil
}
