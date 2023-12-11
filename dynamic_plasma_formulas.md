# Formulas for Dynamic Plasma
## Plasma Base Price
The plasma base price of a momentum. The minimum "price" of plasma an account block has to pay in order to be eligible for inclusion in a momentum.
$$p=max\\{p_p + {p_p \cdot (B_p - B_t) \over B_{t} \cdot c},p_{min}\\}$$

* $p$: The plasma base price of the momentum

* $p_p$: The plasma base price of the previous momentum

* $B_p$: The total base plasma committed to the previous momentum

* $B_t$: The total base plasma target of a momentum. This is the target amount of base plasma to be committed to a momentum (constant: 2,100,000 plasma)

* $c$: The plasma base price change denominator (constant: 8)

* $p_{min}$: The minimum plasma base price (constant: 1)


## Rate Controlling PoW Plasma
### PoW Difficulty Per Plasma
The proof-of-work difficulty per plasma of a momentum.
$$d=max\\{d_p + {d_p \cdot (P_p - P_{t}) \over P_{t} \cdot c},d_{min}\\}$$

* $d$: The difficulty per plasma of the momentum

* $d_p$: The difficulty per plasma of the previous momentum

* $P_p$: The total PoW plasma amount of the previous momentum

* $P_t$: The momentum PoW plasma target

* $d_{min}$: The minimum difficulty per plasma (constant: 500)

### Momentum PoW Plasma Target
The target amount of PoW plasma that should be committed to a momentum.
```math
P_t=\left\{
\begin{array}{ll}
U_p \cdot \sigma, &\text{if }B_p \ge B_t \\ 
U_p, &\text{otherwise}
\end{array} 
\right.
```

* $P_t$: The momentum PoW plasma target

* $U_p$: The total used plasma of the previous momentum

* $\sigma$: The target PoW plasma percent in a momentum (constant: 20%)

* $B_p$: The total base plasma committed to the previous momentum

* $B_t$: The total base plasma target

## Rate Controlling Fused Plasma
### Account's Available Fused Plasma
The currently available fused plasma for an account.
$$F_a=max\\{F_t-F_u-F_c,0\\}$$

* $F_a$: The account's current available fused plasma

* $F_t$: The account's total fused plasma

* $F_u$: The account's uncommitted fused plasma. This is the plasma that is in the mempool waiting to be used in a confirmed account block

* $F_c$: The account's confined fused plasma. This is the account's fused plasma that is confined in previously confirmed account blocks

### Account's Confined Fused Plasma
The currently confined fused plasma of an account. This is the plasma that has to get recharged before it can be used again.
```math
F_c=\left\{
\begin{array}{ll}
max\{F_{cc}-(H_f-H_c)\cdot {F_{cc} / T_c},0\}, &\text{if }T_c>0 \\ 
0, &\text{otherwise}
\end{array} 
\right.
```

* $F_c$: The account's confined fused plasma

* $F_{cc}$: The account's confined fused plasma at confinement height

* $H_f$: The frontier momentum height

* $H_c$: The latest confinement height. This is the height that the account's latest account block was confirmed on

* $T_c$: The confirmation target of the latest confinement

### Confirmation Target
The target amount of confirmations needed to fully recharge the account's confined fused plasma.
$$T={F_b \over \gamma \cdot F_t / \omega} + max\\{T_c-(H_f - H_c),0\\}$$

* $T$: The confirmation target

* $F_b$: The fused plasma committed to the account block

* $\gamma$: The base fusion recharge rate of the confinement momentum

* $F_t$: The account's total fused plasma

* $\omega$: The fusion recharge rate denominator (constant: 210,000)

* $T_c$: The confirmation target of the latest confinement

* $H_f$: The frontier momentum height

* $H_c$: The latest confinement height. This is the height that the account's latest account block was confirmed on

### Base Fusion Recharge Rate
The base fusion recharge rate of a momentum.
```math
\gamma=\left\{
\begin{array}{ll}
 max\{\gamma_p / \alpha,\gamma_t\}, &\text{if }\gamma_p>\gamma_t \\
 min\{\gamma_p \cdot \alpha,\gamma_t\}, &\text{if }\gamma_p<\gamma_t \\
\gamma_t, &\text{otherwise}
\end{array} 
\right.
```

* $\gamma$: The base fusion recharge rate of the momentum

* $\gamma_p$: The base fusion recharge rate of the previous momentum

* $\gamma_t$: The target base fusion recharge rate of the momentum

* $\alpha$: The fusion recharge rate multiplier (constant: 2)

### Target Base Fusion Recharge Rate
The target base fusion recharge rate of a momentum.
$$\gamma_t=\gamma_{max} \cdot ({\gamma_{max} \over \beta / u})^{-2 \cdot B / B_{max}}$$

* $\gamma_t$: The target base fusion recharge rate of the momentum

* $\gamma_{max}$: The maximum base fusion recharge rate (constant: 210,000 plasma/confirmation)

* $\beta$: Account block base plasma (constant: 21,000 plasma)

* $u$: Fusion units per account block base plasma (constant: 10 QSR)

* $B$: The total base plasma committed to the momentum

* $B_{max}$: The maximum allowed total base plasma in a momentum (constant: 4,200,000 plasma)
