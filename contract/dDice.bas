/*  dDice.bas
    Original Version: https://github.com/Nelbert442/dero-smartcontracts/tree/main/DERO-Dice
    Updated Version: https://github.com/newvcas8372/dDice
    Updated Author: newvcas8372
*/

Function InitializePrivate() Uint64
    10  IF EXISTS("owner") == 0 THEN GOTO 15 ELSE GOTO 999
    15  STORE("owner", SIGNER())
    20  STORE("minWager", 5000)
    30  STORE("maxWager", 500000)
    40  STORE("sc_giveback", 9800)
    50  STORE("balance", 0)

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

    190 STORE("minMultiplier", 2)
    191 STORE("maxMultiplier", 10)

    210 RETURN 0
    999 RETURN 1
End Function

Function Donate() Uint64
    10  DIM balance, dvalue as Uint64
    11  LET dvalue = DEROVALUE()
    15  IF dvalue == 0 THEN GOTO 85

	50  LET balance = LOAD("balance") + dvalue
	60  STORE("balance", balance)

	85 RETURN 0
End Function

Function RollDiceHigh(multiplier Uint64) Uint64
    10  DIM rolledNum, targetNumber, payoutAmount, minWager, maxWager, minMultiplier, maxMultiplier, currentHeight, betAmount as Uint64
    11  DIM sendToAddr as String
    13  LET currentHeight = BLOCK_HEIGHT()
    14  LET betAmount = DEROVALUE()
    15  LET sendToAddr = SIGNER()
    16  IF ADDRESS_STRING(sendToAddr) == "" THEN GOTO 500

    40  LET minWager = LOAD("minWager")
    41  LET maxWager = LOAD("maxWager")
    42  LET minMultiplier = LOAD("minMultiplier")
    43  LET maxMultiplier = LOAD("maxMultiplier")
    45  IF betAmount < minWager THEN GOTO 900
    50  IF betAmount > maxWager THEN GOTO 900
    55  LET payoutAmount = LOAD("sc_giveback") * betAmount * multiplier / 10000
    
    60  IF EXISTS("Over-x" + ITOA(multiplier)) == 1 THEN GOTO 70 ELSE GOTO 900

    70  LET rolledNum = RANDOM(99)
    80  LET targetNumber = LOAD("Over-x" + ITOA(multiplier))
    85  STORE(ITOA(multiplier) + "xPlays", LOAD(ITOA(multiplier) + "xPlays") + 1)
    90  IF rolledNum >= targetNumber THEN GOTO 100 ELSE GOTO 500

    100 IF LOAD("balance") < payoutAmount THEN GOTO 700
    120 SEND_DERO_TO_ADDRESS(sendToAddr, payoutAmount)
    125 STORE("balance", LOAD("balance") + (betAmount - payoutAmount))
    126 STORE(ITOA(multiplier) + "xWins", LOAD(ITOA(multiplier) + "xWins") + 1)
    130 RETURN 0

    500 STORE("balance", LOAD("balance") + betAmount)
    505 RETURN 0

    700 STORE(ITOA(multiplier) + "xWins", LOAD(ITOA(multiplier) + "xWins") + 1)
    710 SEND_DERO_TO_ADDRESS(sendToAddr, betAmount)
    720 RETURN 0

    900 SEND_DERO_TO_ADDRESS(sendToAddr, betAmount)
    910 RETURN 0
End Function

Function RollDiceLow(multiplier Uint64) Uint64
    10  DIM rolledNum, targetNumber, payoutAmount, minWager, maxWager, minMultiplier, maxMultiplier, currentHeight, betAmount as Uint64
    11  DIM sendToAddr as String
    13  LET currentHeight = BLOCK_HEIGHT()
    14  LET betAmount = DEROVALUE()
    15  LET sendToAddr = SIGNER()
    16  IF ADDRESS_STRING(sendToAddr) == "" THEN GOTO 500

    40  LET minWager = LOAD("minWager")
    41  LET maxWager = LOAD("maxWager")
    42  LET minMultiplier = LOAD("minMultiplier")
    43  LET maxMultiplier = LOAD("maxMultiplier")
    45  IF betAmount < minWager THEN GOTO 900
    50  IF betAmount > maxWager THEN GOTO 900
    55  LET payoutAmount = LOAD("sc_giveback") * betAmount * multiplier / 10000
    
    60  IF EXISTS("Under-x" + ITOA(multiplier)) == 1 THEN GOTO 70 ELSE GOTO 900

    70  LET rolledNum = RANDOM(99)
    80  LET targetNumber = LOAD("Under-x" + ITOA(multiplier))
    85  STORE(ITOA(multiplier) + "xPlays", LOAD(ITOA(multiplier) + "xPlays") + 1)
    90  IF rolledNum <= targetNumber THEN GOTO 100 ELSE GOTO 500

    100 IF LOAD("balance") < payoutAmount THEN GOTO 700
    120 SEND_DERO_TO_ADDRESS(sendToAddr, payoutAmount)
    125 STORE("balance", LOAD("balance") + (betAmount - payoutAmount))
    126 STORE(ITOA(multiplier) + "xWins", LOAD(ITOA(multiplier) + "xWins") + 1)
    130 RETURN 0

    500 STORE("balance", LOAD("balance") + betAmount)
    505 RETURN 0

    700 STORE(ITOA(multiplier) + "xWins", LOAD(ITOA(multiplier) + "xWins") + 1)
    710 SEND_DERO_TO_ADDRESS(sendToAddr, betAmount)
    720 RETURN 0

    900 SEND_DERO_TO_ADDRESS(sendToAddr, betAmount)
    910 RETURN 0
End Function

Function TransferOwnership(newowner String) Uint64 
    10  IF LOAD("owner") == SIGNER() THEN GOTO 30 
    20  RETURN 1
    30  STORE("tmpowner",ADDRESS_RAW(newowner))
    40  RETURN 0
End Function

Function ClaimOwnership() Uint64 
    10  IF LOAD("tmpowner") == SIGNER() THEN GOTO 30 
    20  RETURN 1
    30  STORE("owner",SIGNER())
    40  RETURN 0
End Function

Function Withdraw(amount Uint64) Uint64
    10  IF LOAD("owner") == SIGNER() THEN GOTO 20 ELSE GOTO 50
    20  IF LOAD("balance") < amount THEN GOTO 50
    30  SEND_DERO_TO_ADDRESS(SIGNER(), amount)
    40  STORE("balance", LOAD("balance") - amount)
    50  RETURN 0
End Function