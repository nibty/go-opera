package emitter

import (
	"time"
	
	"github.com/Fantom-foundation/lachesis-base/emitter/ancestor"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/Fantom-foundation/lachesis-base/inter/pos"
	"github.com/Fantom-foundation/lachesis-base/utils/piecefunc"

	"github.com/Fantom-foundation/go-opera/inter"
)

func scalarUpdMetric(diff idx.Event, weight pos.Weight, totalWeight pos.Weight) ancestor.Metric {
	return ancestor.Metric(scalarUpdMetricF(uint64(diff)*piecefunc.DecimalUnit)) * ancestor.Metric(weight) / ancestor.Metric(totalWeight)
}

func updMetric(median, cur, upd idx.Event, validatorIdx idx.Validator, validators *pos.Validators) ancestor.Metric {
	if upd <= median || upd <= cur {
		return 0
	}
	weight := validators.GetWeightByIdx(validatorIdx)
	if median < cur {
		return scalarUpdMetric(upd-median, weight, validators.TotalWeight()) - scalarUpdMetric(cur-median, weight, validators.TotalWeight())
	}
	return scalarUpdMetric(upd-median, weight, validators.TotalWeight())
}

func kickStartMetric(metric ancestor.Metric, seq idx.Event) ancestor.Metric {
	// kickstart metric in a beginning of epoch, when there's nothing to observe yet
	if seq <= 2 && metric < 0.9*piecefunc.DecimalUnit {
		metric += 0.1 * piecefunc.DecimalUnit
	}
	if seq <= 1 && metric <= 0.8*piecefunc.DecimalUnit {
		metric += 0.2 * piecefunc.DecimalUnit
	}
	return metric
}

func eventMetric(orig ancestor.Metric, seq idx.Event) ancestor.Metric {
	return kickStartMetric(ancestor.Metric(eventMetricF(uint64(orig))), seq)
}

// Function to get the top 50 elements from a slice
 func top50(slice []idx.ValidatorID) []idx.ValidatorID {
     if len(slice) > 50 {
         return slice[:50] // Return the first 100 elements
     }
     return slice // Return the slice as is if it has less than or equal to 100 elements
 }

 // Function to check if a number is in the top 50 elements of a slice
 func isInTop50(number idx.ValidatorID, slice []idx.ValidatorID) bool {
     var v idx.ValidatorID
     top50Slice := top50(slice)
     for _, v = range top50Slice {
         if v == number {
             return true
         }
     }
     return false
 }

func (em *Emitter) isAllowedToEmit(e inter.EventI, eTxs bool, metric ancestor.Metric, selfParent *inter.Event) bool {
	// for now allow only vals up to ID 30 to emit:
	passedTime := e.CreationTime().Time().Sub(em.prevEmittedAtTime)
	if passedTime < 0 {
		passedTime = 0
	}
	passedTimeIdle := e.CreationTime().Time().Sub(em.prevIdleTime)
	if passedTimeIdle < 0 {
		passedTimeIdle = 0
	}
	// metric is a decimal (0.0, 1.0], being an estimation of how much the event will advance the consensus
	adjustedPassedTime := time.Duration(ancestor.Metric(passedTime/piecefunc.DecimalUnit) * metric)
	adjustedPassedIdleTime := time.Duration(ancestor.Metric(passedTimeIdle/piecefunc.DecimalUnit) * metric)
	passedBlocks := em.world.GetLatestBlockIndex() - em.prevEmittedAtBlock

	supermajority := true
	// Filter this node's events if not in top50 supermajority of stakers
        if e.Creator() == em.config.Validator.ID && !isInTop50(e.Creator(), em.validators.SortedIDs()) {
                //fmt.Println("This node is not in supermajority")
                //supermajority = false
				// disable check
				supermajority = true
        }

    if (supermajority) {
	if em.stakeRatio[e.Creator()] < 0.35*piecefunc.DecimalUnit {
		// top validators emit event right after transaction is originated
		passedTimeIdle = passedTime
	} else if em.stakeRatio[e.Creator()] < 0.7*piecefunc.DecimalUnit {
		// top validators emit event right after transaction is originated
		passedTimeIdle = (passedTimeIdle + passedTime) / 2
	}
	if passedTimeIdle > passedTime {
		passedTimeIdle = passedTime
	}
	// Forbid emitting if not enough power and power is decreasing
	{
		threshold := em.config.EmergencyThreshold
		if e.GasPowerLeft().Min() <= threshold {
			if selfParent != nil && e.GasPowerLeft().Min() < selfParent.GasPowerLeft().Min() {
				em.Periodic.Warn(10*time.Second, "Not enough power to emit event, waiting",
					"power", e.GasPowerLeft().String(),
					"selfParentPower", selfParent.GasPowerLeft().String(),
					"stake%", 100*float64(em.validators.Get(e.Creator()))/float64(em.validators.TotalWeight()))
				return false
			}
		}
	}
	// Enforce emitting if passed too many time/blocks since previous event
	{
		rules := em.world.GetRules()
		maxBlocks := rules.Economy.BlockMissedSlack/2 + 1
		if rules.Economy.BlockMissedSlack > maxBlocks && maxBlocks < rules.Economy.BlockMissedSlack-5 {
			maxBlocks = rules.Economy.BlockMissedSlack - 5
		}
		if passedTime >= em.intervals.Max ||
			passedBlocks >= maxBlocks*4/5 && metric >= piecefunc.DecimalUnit/2 ||
			passedBlocks >= maxBlocks {
			return true
		}
	}
	// Slow down emitting if power is low
	{
		threshold := (em.config.NoTxsThreshold + em.config.EmergencyThreshold) / 2
		if e.GasPowerLeft().Min() <= threshold {
			// it's emitter, so no need in determinism => fine to use float
			minT := float64(em.intervals.Min)
			maxT := float64(em.intervals.Max)
			factor := float64(e.GasPowerLeft().Min()) / float64(threshold)
			adjustedEmitInterval := time.Duration(maxT - (maxT-minT)*factor)
			if passedTime < adjustedEmitInterval {
				return true
			}
		}
	}
	// Slow down emitting if no txs to confirm/originate
	{
		if passedTime < em.intervals.Max &&
			em.idle() &&
			!eTxs {
			return true
		}
	}
	// Emitting is controlled by the efficiency metric
	{
		if passedTime < em.intervals.Min {
			return true
		}
		if adjustedPassedTime < em.intervals.Min &&
			!em.idle() {
			return true
		}
		if adjustedPassedIdleTime < em.intervals.Confirming &&
			!em.idle() &&
			!eTxs {
			return true
		}
	}
    // only allow top validators
	return true
	} else { 
	        // Enforce emitting if passed Max time (10 mins)
                if passedTime >= em.intervals.Max {
                        return true
                }
	}
	return false
}

func (em *Emitter) recheckIdleTime() {
	em.world.Lock()
	defer em.world.Unlock()
	if em.idle() {
		em.prevIdleTime = time.Now()
	}
}
