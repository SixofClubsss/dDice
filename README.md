## contract/dDice.bas
Attempt at similar product as Ether Dice etc. Dice rolling game in which you can choose between a 2x and a 10x multiplier (increment by 1s [e.g. 2x, 3x, 4x, ... 10x]) and roll high or low.
The high and low numbers are defined as such:
```
    2x --> 50 or over --> 49 or under
    3x --> 67 or over --> 33 or under
    4x --> 75 or over --> 25 or under
    5x --> 80 or over --> 20 or under
    6x --> 84 or over --> 16 or under
    7x --> 86 or over --> 14 or under
    8x --> 88 or over --> 12 or under
    9x --> 89 or over --> 11 or under
    10x --> 90 or over --> 10 or under
```

### Disclaimer
We are not responsible for any lost funds through the usage of this contract. Please deploy and utilize at your own risk. ALWAYS USE RINGSIZE=2 when interacting with this contract to prevent loss of funds, see [line 16](https://github.com/newvcas8372/dDice/blob/main/contract/dDice.bas#L77) in each roll function. You CAN donate anonymously with ringsize > 2.

### SCID (Contract ID)
[fe61b1ac6edbe18180d2863f05d1dfb26a767abdfc0488cbe4970d950ef45de8](https://explorer.dero.io/tx/fe61b1ac6edbe18180d2863f05d1dfb26a767abdfc0488cbe4970d950ef45de8)

### e.x.1 (Roll High with 2x Multiplier - Wagering 0.05 DERO):
TX Fee: ~0.00258
```
curl http://127.0.0.1:10103/json_rpc -d '{"jsonrpc":"2.0","id":"0","method":"scinvoke","params":{"sc_dero_deposit":5000,"ringsize":2,"scid":"fe61b1ac6edbe18180d2863f05d1dfb26a767abdfc0488cbe4970d950ef45de8","sc_rpc":[{"name":"entrypoint","datatype":"S","value":"RollDiceHigh"},{"name":"multiplier","datatype":"U","value":2}] }}' -H 'Content-Type: application/json'
```

### e.x.2 (Roll Low with 5x Multiplier - Wagering 0.1 DERO):
TX Fee: ~0.00258
```
curl http://127.0.0.1:10103/json_rpc -d '{"jsonrpc":"2.0","id":"0","method":"scinvoke","params":{"sc_dero_deposit":10000,"ringsize":2,"scid":"fe61b1ac6edbe18180d2863f05d1dfb26a767abdfc0488cbe4970d950ef45de8","sc_rpc":[{"name":"entrypoint","datatype":"S","value":"RollDiceLow"},{"name":"multiplier","datatype":"U","value":5}] }}' -H 'Content-Type: application/json'
```

### DONATE (Donates DERO to dDice Liquidity Anonymously - Donating 1 DERO):
TX Fee: ~0.00289
```
curl http://127.0.0.1:10103/json_rpc -d '{"jsonrpc":"2.0","id":"0","method":"scinvoke","params":{"sc_dero_deposit":100000,"ringsize":16,"scid":"fe61b1ac6edbe18180d2863f05d1dfb26a767abdfc0488cbe4970d950ef45de8","sc_rpc":[{"name":"entrypoint","datatype":"S","value":"Donate"}] }}' -H 'Content-Type: application/json'
```

### DERO Dice Template (Install your own!)

Deploy the `contract/dDice.bas` contents and list the deployed SCID into your dApp.

Install dDice
```
curl --request POST --data-binary @dDice.bas http://127.0.0.1:10103/install_sc
```
Cost to deploy: ~0.08745 (possibly optimized over time/updates)

Cost to play: ~0.00258 (possibly optimized over time/updates)

Comment-heavy codebase:
```go
/*  dDice.bas
    Original Version: https://github.com/Nelbert442/dero-smartcontracts/tree/main/DERO-Dice
    Updated Version: https://github.com/newvcas8372/dDice
    Updated Author: newvcas8372
*/

Function InitializePrivate() Uint64
    10  IF EXISTS("owner") == 0 THEN GOTO 15 ELSE GOTO 999
    15  STORE("owner", SIGNER())
    20  STORE("minWager", 5000)  // Sets minimum wager (DERO is 5 atomic units)
    30  STORE("maxWager", 500000)  // Sets maximum wager (DERO is 5 atomic units)
    40  STORE("sc_giveback", 9800)  // Sets the SC giveback on reward payout, 2% to pool, 98% to winner (9800) for example
    50  STORE("balance", 0) // Tracks balance

    // Defines the over/under amounts to hit via RANDOM() in order to win for each func
    60  STORE("Over-x2", 50)
    61  STORE("Under-x2", 49)
    65  STORE("Over-x3", 67)
    66  STORE("Under-x3", 33)
    70  STORE("Over-x4", 75)
    71  STORE("Under-x4", 25)
    75  STORE("Over-x5", 80)
    76  STORE("Under-x5", 20)
    80  STORE("Over-x6", 84)
    81  STORE("Under-x6", 16)
    85  STORE("Over-x7", 86)
    86  STORE("Under-x7", 14)
    90  STORE("Over-x8", 88)
    91  STORE("Under-x8", 12)
    95  STORE("Over-x9", 89)
    96  STORE("Under-x9", 11)
    100 STORE("Over-x10", 90)
    101 STORE("Under-x10", 10)

    // In-contract stats tracking for total plays (per multiplier) and wins to calculate historical odds
    120 STORE("2xPlays", 0)
    121 STORE("2xWins", 0)
    125 STORE("3xPlays", 0)
    126 STORE("3xWins", 0)
    130 STORE("4xPlays", 0)
    131 STORE("4xWins", 0)
    135 STORE("5xPlays", 0)
    136 STORE("5xWins", 0)
    140 STORE("6xPlays", 0)
    141 STORE("6xWins", 0)
    145 STORE("7xPlays", 0)
    146 STORE("7xWins", 0)
    150 STORE("8xPlays", 0)
    151 STORE("8xWins", 0)
    155 STORE("9xPlays", 0)
    156 STORE("9xWins", 0)
    160 STORE("10xPlays", 0)
    161 STORE("10xWins", 0)

    190 STORE("minMultiplier", 2) // Sets the minimum multiplier. If this is modified, be sure to add over/under references above
    191 STORE("maxMultiplier", 10)  // Sets the maximum multiplier. If this is modified, be sure to add over/under references above

    210 RETURN 0
    999 RETURN 1
End Function

// Donates balance to the SC. This can be done anonymously as no SIGNER() method is used
Function Donate() Uint64
    10  DIM balance, dvalue as Uint64
    11  LET dvalue = DEROVALUE()
    15  IF dvalue == 0 THEN GOTO 85 // If value is 0, simply return

	50  LET balance = LOAD("balance") + dvalue
	60  STORE("balance", balance)

	85 RETURN 0
End Function

// Call to roll dice against over-x values in order to win
Function RollDiceHigh(multiplier Uint64) Uint64
    10  DIM rolledNum, targetNumber, payoutAmount, minWager, maxWager, minMultiplier, maxMultiplier, currentHeight, betAmount as Uint64
    11  DIM sendToAddr as String
    13  LET currentHeight = BLOCK_HEIGHT()
    14  LET betAmount = DEROVALUE()
    15  LET sendToAddr = SIGNER()
    16  IF ADDRESS_STRING(sendToAddr) == "" THEN GOTO 500   // If ringsize is != 2, we just return 0, append balance and close out. We cannot send funds back or anything, so it is added to SC balance. This should be WARNING on all dApp frontends

    40  LET minWager = LOAD("minWager")
    41  LET maxWager = LOAD("maxWager")
    42  LET minMultiplier = LOAD("minMultiplier")
    43  LET maxMultiplier = LOAD("maxMultiplier")
    45  IF betAmount < minWager THEN GOTO 900 // If value is less than stored minimum wager, send bet DERO back.
    50  IF betAmount > maxWager THEN GOTO 900 // If value is more than stored maximum wager, send bet DERO back
    55  LET payoutAmount = LOAD("sc_giveback") * betAmount * multiplier / 10000
    
    60  IF EXISTS("Over-x" + ITOA(multiplier)) == 1 THEN GOTO 70 ELSE GOTO 900

    70  LET rolledNum = RANDOM(99)  // Randomly choose a number between 0 and 99
    80  LET targetNumber = LOAD("Over-x" + ITOA(multiplier))
    85  STORE(ITOA(multiplier) + "xPlays", LOAD(ITOA(multiplier) + "xPlays") + 1)   // Append 1 play to the multiplier plays for stats/odds
    90  IF rolledNum >= targetNumber THEN GOTO 100 ELSE GOTO 500

    100 IF LOAD("balance") < payoutAmount THEN GOTO 700 // If balance cannot cover the potential winnings, error out and send DERO back to SIGNER()
    120 SEND_DERO_TO_ADDRESS(sendToAddr, payoutAmount)
    125 STORE("balance", LOAD("balance") + (betAmount - payoutAmount))
    126 STORE(ITOA(multiplier) + "xWins", LOAD(ITOA(multiplier) + "xWins") + 1) // Append 1 win to the multiplier wins for stats/odds
    130 RETURN 0

    500 STORE("balance", LOAD("balance") + betAmount)
    505 RETURN 0

    700 STORE(ITOA(multiplier) + "xWins", LOAD(ITOA(multiplier) + "xWins") + 1)
    710 SEND_DERO_TO_ADDRESS(sendToAddr, betAmount)
    720 RETURN 0

    900 SEND_DERO_TO_ADDRESS(sendToAddr, betAmount)
    910 RETURN 0
End Function

// Call to roll dice against under-x values in order to win
Function RollDiceLow(multiplier Uint64) Uint64
    10  DIM rolledNum, targetNumber, payoutAmount, minWager, maxWager, minMultiplier, maxMultiplier, currentHeight, betAmount as Uint64
    11  DIM sendToAddr as String
    13  LET currentHeight = BLOCK_HEIGHT()
    14  LET betAmount = DEROVALUE()
    15  LET sendToAddr = SIGNER()
    16  IF ADDRESS_STRING(sendToAddr) == "" THEN GOTO 500   // If ringsize is != 2, we just return 0, append balance and close out. We cannot send funds back or anything, so it is added to SC balance. This should be WARNING on all dApp frontends

    40  LET minWager = LOAD("minWager")
    41  LET maxWager = LOAD("maxWager")
    42  LET minMultiplier = LOAD("minMultiplier")
    43  LET maxMultiplier = LOAD("maxMultiplier")
    45  IF betAmount < minWager THEN GOTO 900 // If value is less than stored minimum wager, send bet DERO back.
    50  IF betAmount > maxWager THEN GOTO 900 // If value is more than stored maximum wager, send bet DERO back
    55  LET payoutAmount = LOAD("sc_giveback") * betAmount * multiplier / 10000
    
    60  IF EXISTS("Under-x" + ITOA(multiplier)) == 1 THEN GOTO 70 ELSE GOTO 900

    70  LET rolledNum = RANDOM(99)  // Randomly choose a number between 0 and 99
    80  LET targetNumber = LOAD("Under-x" + ITOA(multiplier))
    85  STORE(ITOA(multiplier) + "xPlays", LOAD(ITOA(multiplier) + "xPlays") + 1)   // Append 1 play to the multiplier plays for stats/odds
    90  IF rolledNum <= targetNumber THEN GOTO 100 ELSE GOTO 500

    100 IF LOAD("balance") < payoutAmount THEN GOTO 700 // If balance cannot cover the potential winnings, error out and send DERO back to SIGNER()
    120 SEND_DERO_TO_ADDRESS(sendToAddr, payoutAmount)
    125 STORE("balance", LOAD("balance") + (betAmount - payoutAmount))
    126 STORE(ITOA(multiplier) + "xWins", LOAD(ITOA(multiplier) + "xWins") + 1) // Append 1 win to the multiplier wins for stats/odds
    130 RETURN 0

    500 STORE("balance", LOAD("balance") + betAmount)
    505 RETURN 0

    700 STORE(ITOA(multiplier) + "xWins", LOAD(ITOA(multiplier) + "xWins") + 1)
    710 SEND_DERO_TO_ADDRESS(sendToAddr, betAmount)
    720 RETURN 0

    900 SEND_DERO_TO_ADDRESS(sendToAddr, betAmount)
    910 RETURN 0
End Function

// Transfer ownership to another address
Function TransferOwnership(newowner String) Uint64 
    10  IF LOAD("owner") == SIGNER() THEN GOTO 30 
    20  RETURN 1
    30  STORE("tmpowner",ADDRESS_RAW(newowner))
    40  RETURN 0
End Function

// Claim ownership
Function ClaimOwnership() Uint64 
    10  IF LOAD("tmpowner") == SIGNER() THEN GOTO 30 
    20  RETURN 1
    30  STORE("owner",SIGNER())
    40  RETURN 0
End Function

// Withdraw a given amount of DERO from the contract
Function Withdraw(amount Uint64) Uint64
    10  IF LOAD("owner") == SIGNER() THEN GOTO 20 ELSE GOTO 50
    20  IF LOAD("balance") < amount THEN GOTO 50
    30  SEND_DERO_TO_ADDRESS(SIGNER(), amount)
    40  STORE("balance", LOAD("balance") - amount)
    50  RETURN 0
End Function
```

### Future
Future item will include a template of a dApp/local html/js page which will query the stats and provide necessary details for playing etc.