package lnwallet

import (
	"github.com/decred/dcrd/dcrutil"
	"github.com/decred/dcrwallet/wallet/v2/txrules"
)

// lnTxVersion is the version that transactions need to be defined to use so
// that they are usable as ln transactions.
const lnTxVersion uint16 = 2

// DefaultDustLimit is used to calculate the default dust threshold limit,
// assuming the network uses the defaultRelayFeePerKb of the wallet.
func DefaultDustLimit() dcrutil.Amount {
	return DustThresholdForRelayFee(AtomPerKByte(txrules.DefaultRelayFeePerKb))
}

// DustThresholdForRelayFee returns the minimum amount an output of the given
// size should have in order to not be considered a dust output for networks
// using the provided relay fee.
//
// It is assumed the output is paying to a P2PKH script.
func DustThresholdForRelayFee(relayFeeRate AtomPerKByte) dcrutil.Amount {
	// Size to redeem a p2pkh script is the size of an input + size of a
	// serialized p2pkh signature script (varint length + OP_DATA_73 + 73 +
	// OP_DATA_33 + 33)
	inputRedeemSize := inputSize + P2PKHSigScriptSize

	// Calculate the total (estimated) cost to the network.  This is
	// calculated using the serialize size of the output plus the serial
	// size of a transaction input which redeems it.  The output is assumed
	// to be compressed P2PKH as this is the most common script type. The serialized
	// varint size of a P2PKH script is 1.
	scriptSize := P2PKHPkScriptSize
	totalSize := outputSize + 1 + scriptSize + inputRedeemSize

	// Calculate the relay fee for this test tx in atoms, given its
	// estimated totalSize and the provided relayFeeRate in atoms/kB.
	relayFee := totalSize * int64(relayFeeRate) / 1000

	// Threshold for dustiness is determined as 3 times the relay fee.
	return dcrutil.Amount(3 * relayFee)
}
