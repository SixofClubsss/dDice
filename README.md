# dDice

dDice, decentralized dice games built on [Dero's](https://dero.io) decentralized application platform. 

[dDice origin](https://github.com/newvcas8372/dDice)

[dDice7 expansion](#ddice7)

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

### dDice7
dDice7 was built off the dDice foundation. dDice7 is a craps style expansion. dDice7 UI is built using [Fyne toolkit](https://fyne.io/) and powered by [Gnomon](https://github.com/civilware/Gnomon) and [dReams](https://dreamdapps.io). 

![goMod](https://img.shields.io/github/go-mod/go-version/SixofClubsss/dDice.svg)![goReport](https://goreportcard.com/badge/github.com/SixofClubsss/dDice)[![goDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/SixofClubsss/dDice)

dDice7 can be played on [dReams](https://dreamdapps.io/download) app.

### SCID (Contract ID)
[fed996730a15744c941d4722db0b1a36dc650939dbf66c246aa7e74f38e409cd](https://explorer.dero.io/tx/fed996730a15744c941d4722db0b1a36dc650939dbf66c246aa7e74f38e409cd)

### Added features
- Craps style game with proposition (one time) and place (recurring) bets
- Multiplayer table where each roll effects all players at the table
- Multi currency betting

### Bets
- Prop (single roll)
    - (0) Under 7 pays 1:1
    - (1) Over 7 pays 1:1
    - (2) Any 7 pays 4:1
    - (3) Any crap (2 or 3 or 12)  pays 7:1 
    - (4) Ace deuce (1 and 2) pays 15:1
    - (5) YO (5 & 6) pays 15:1
    - (6) Aces (1 & 1) pays 30:1
    - (7) Midnight (6 & 6) pays 30:1
- Place (multi roll)
    - (0) Place 4 pays 9:5
    - (1) Place 5 pays 7:5
    - (2) Place 6 pays 7:6
    - (3) Place 8 pays 7:6
    - (4) Place 9 pays 7:5
    - (5) Place 10 pays 9:5
    - (6) Field (3, 4, 9, 10 11) pays 1:1 and (2, 12) pays 2:1
    - Inside (Place on  5, 6, 8 and 9)
    - Outside (Place on  4 and 10)

### e.x.1 (Roll dice with Over 7 bet - Wagering 0.05 DERO):
TX Fee: ~0.00380
```
curl -u user:pass http://127.0.0.1:10103/json_rpc -d '{"jsonrpc":"2.0","id":"0","method":"transfer","params":{"transfers":
[{"destination":"dero1qy3ef7nudlawzmwk90n4dtqpyadgl4cw04th3wnrnk2g5nz457adqqg7j9kw3", "burn":5000}], "fees":380, "scid":"fed996730a15744c941d4722db0b1a36dc650939dbf66c246aa7e74f38e409cd","ringsize":2, "sc_rpc":[{"name":"entrypoint","datatype":"S","value":"Roll"}, {"name":"bet","datatype":"U","value":1}] }}' -H 'Content-Type: application/json';
```

### e.x.2 (Place bet on 8 - Wagering 10 dReams):
TX Fee: ~0.00280
```
curl -u user:pass http://127.0.0.1:10103/json_rpc -d '{"jsonrpc":"2.0","id":"0","method":"transfer","params":{"transfers":
[{"scid":"ad2e7b37c380cc1aed3a6b27224ddfc92a2d15962ca1f4d35e530dba0f9575a9", "burn":1000000}], "fees":280, "scid":"fed996730a15744c941d4722db0b1a36dc650939dbf66c246aa7e74f38e409cd","ringsize":2, "sc_rpc":[{"name":"entrypoint","datatype":"S","value":"Place"}, {"name":"p","datatype":"U","value":3}] }}' -H 'Content-Type: application/json';
```

### e.x.3 (Inside bet - Wagering 0.2 DERO total, 0.05 placed on 5, 6, 8 and 9 respectively):
TX Fee: ~0.00450
```
curl -u user:pass http://127.0.0.1:10103/json_rpc -d '{"jsonrpc":"2.0","id":"0","method":"transfer","params":{"transfers":
[{"destination":"dero1qy3ef7nudlawzmwk90n4dtqpyadgl4cw04th3wnrnk2g5nz457adqqg7j9kw3", "burn":20000}], "fees":450, "scid":"fed996730a15744c941d4722db0b1a36dc650939dbf66c246aa7e74f38e409cd","ringsize":2, "sc_rpc":[{"name":"entrypoint","datatype":"S","value":"Inside"}] }}' -H 'Content-Type: application/json';
```

### Licensing
dDice is free and open source.
dDice smart contract forked from [newvcas8372/dDice](https://github.com/newvcas8372/dDice) under [BSD 3-Clause License](https://github.com/newvcas8372/dDice/blob/main/LICENSE).    
dDice expansion source code is published under the [MIT License](https://github.com/SixofClubsss/dDice/blob/main/dice/LICENSE).   
Copyright © 2024 SixofClubs        
