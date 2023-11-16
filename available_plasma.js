// Pseudo code example for calculating available fused plasma
// for account chain based on dynamic plasma recharge rate.

function main() {
    const fusedPlasma = fusedPlasma(account);
    const rechargedPlasma = rechargedPlasma(account, fusedPlasma);
    const confinedPlasma = confinedPlasma(account, rechargedPlasma);
    const availablePlasma = availablePlasma(fusedPlasma, confinedPlasma);
}

function fusedPlasma() {
    const fused = getStakeBeneficialAmount(account.address);
    return fusedAmountToPlasma(fused);
}

// Get the amount of fused plasma that has been recharged since the previous
// account block was sent by the account chain.
function rechargedPlasma(account, fusedPlasma) {
    const frontierBlock = account.storage.frontierBlock();
    const confirmationMomentum = frontierBlock.getConfirmationMomentum();
    const frontierMomentum = getFrontierMomentum();
    const confirmations = frontierMomentum.height - confirmationMomentum.height;
    const rechargeRate = confirmationMomentum.FusionRechargeRate * rechargeRateMultiplier(fusedPlasma)
    const rechargedPlasma = confirmations * rechargeRate;
    return Math.min(rechargedPlasma, fusedPlasma);
}

// Get the amount of plasma that is currently "confined" in account blocks.
// Confined plasma has to get recharged before it can be used again.
// The confined plasma amount is updated to the storage whenever an account
// block is added to the mempool.
function confinedPlasma(account, rechargedPlasma) {
   const confined = account.storage.getConfinedPlasma() - rechargedPlasma;
   return Math.max(confined, 0);
}

// Get the amount of available fused plasma for the account.
function availablePlasma(fusedPlasma, confinedPlasma) {
    const available = fusedPlasma - confinedPlasma;
    return Math.max(available, 0);
}

// The recharge rate multiplier increases an account's recharge rate when
// the account has fused at least the threshold plasma amount. 
// This should probably be accompanied with an account level limit of how many
// transactions an account can send within n momentums when the account's
// fused plasma exceeds or equals thresholdPlasma.
function rechargeRateMultiplier(fusedPlasma) {
    const thresholdPlasma = 2100000 // Over or equal to 1000 QSR
    if (fusedPlasma >= thresholdPlasma) { 
        return fusedPlasma / thresholdPlasma
    } else {
        return 1;
    }
}