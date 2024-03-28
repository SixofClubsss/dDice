// dDice7, adapted from newvcas8372/dDice
Function Donate() Uint64
1 IF pW() THEN GOTO 5
3 RETURN 1
5 STORE("bal", LOAD("bal")+DEROVALUE())
6 STORE("balD", LOAD("balD")+ASSETVALUE(LOAD("dReams")))
9 RETURN 0
End Function

Function Roll(bet Uint64) Uint64
1 IF LOAD("live") != 1 THEN GOTO 3
2 IF bet < 8 THEN GOTO 10
3 RETURN 1
10 DIM amt, roll, pay, die1, die2, t as Uint64
14 LET t = pW()
15 LET amt = value(t)
20 IF limit(t, amt, 1) THEN GOTO 3
25 DIM addr as String
30 LET addr = SIGNER()
40 LET die1 = 1+RANDOM(6)
41 LET die2 = 1+RANDOM(6)
42 LET roll = storeRoll(bet, die1, die2)
43 IF roll != 7 THEN GOTO 48
44 seven()
45 GOTO 50
48 placed(roll)
50 IF hit(bet, roll) THEN GOTO 100
52 storeBal(t, amt, 0, 0)
55 RETURN 0
100 LET pay = payRoll(bet, amt)
105 IF findBal(t) < pay THEN GOTO 200
120 send(t, addr, pay)
125 storeBal(t, amt, pay, 0)
130 RETURN 0
200 send(t, addr, amt)
220 RETURN 0
End Function

Function Place(p Uint64) Uint64
1 IF LOAD("live") != 1 THEN GOTO 4
2 IF LOAD("bets") > 29 THEN GOTO 4
3 IF p < 7 THEN GOTO 10
4 RETURN 1
10 DIM amt, b, t as Uint64
11 LET t = pW()
12 LET amt = value(t)
13 IF limit(t, amt, 1) THEN GOTO 4
15 LET b = b+1
16 IF b > 30 THEN GOTO 4
17 IF EXISTS(b) THEN GOTO 15
18 storeBet(t, b, p, amt)
40 RETURN 0
End Function

Function limit(t Uint64, amt Uint64, m Uint64) Uint64
1 IF t < 2 THEN GOTO 3
2 LET t = 300
3 IF amt < (LOAD("min")*m)*t THEN GOTO 7
4 IF amt > (LOAD("max")*m)*t THEN GOTO 7
5 RETURN 0
7 RETURN 1
End Function

Function Inside() Uint64
1 IF LOAD("live") != 1 THEN GOTO 4
2 IF LOAD("bets") < 27 THEN GOTO 10
4 RETURN 1
10 DIM amt, b, p, t as Uint64
11 LET t = pW()
12 LET amt = value(t)
13 IF limit(t, amt, 4) || amt%4 != 0 THEN GOTO 4
15 LET b = b+1
16 IF b > 30 THEN GOTO 4
17 IF EXISTS(b) THEN GOTO 15
18 LET p = p+1
19 storeBet(t, b, p, amt/4)
20 IF p < 4 THEN GOTO 15
40 RETURN 0
End Function

Function Outside() Uint64
1 IF LOAD("live") != 1 THEN GOTO 4
2 IF LOAD("bets") < 29 THEN GOTO 10
4 RETURN 1
10 DIM amt, b, p, t as Uint64
11 LET t = pW()
12 LET amt = value(t)
13 IF limit(t, amt, 2) || amt%2 != 0 THEN GOTO 4
15 LET b = b+1
16 IF b > 30 THEN GOTO 4
17 IF EXISTS(b) THEN GOTO 15
18 storeBet(t, b, p, amt/2)
19 LET p = p+5
20 IF p == 5 THEN GOTO 15
40 RETURN 0
End Function

Function Clear(b Uint64) Uint64
1 IF EXISTS("b_"+b) == 0 THEN GOTO 3
2 IF LOAD("b_"+b) == SIGNER() THEN GOTO 4
3 RETURN 1
4 DIM i, amt, t as Uint64
5 IF EXISTS(i) == 0 THEN GOTO 30
6 IF LOAD("b_"+i) != SIGNER() THEN GOTO 30
7 IF EXISTS("b_"+i+"t") THEN GOTO 10
8 LET t = 1
9 GOTO 11
10 LET t = 2
11 LET amt = LOAD("b_"+i+"amt")
12 IF t == 1 THEN GOTO 15
13 IF t == 2 THEN GOTO 17
14 GOTO 3
15 IF LOAD("bal") < amt THEN GOTO 30
16 GOTO 21
17 IF LOAD("balD") < amt THEN GOTO 30
21 send(t, LOAD("b_"+i), amt)
22 STORE("bets", LOAD("bets")-1)
23 storeBal(t, 0, amt, 1)
24 DELETE(i)
25 DELETE("b_"+i)
26 DELETE("b_"+i+"amt")
27 DELETE("b_"+i+"h")
28 DELETE("b_"+i+"t")
30 LET i = i+1
40 IF i <= 30 THEN GOTO 5
50 RETURN 0
End Function

Function pW() Uint64
1 IF DEROVALUE() > 0 THEN GOTO 5
2 IF ASSETVALUE(LOAD("dReams")) > 0 THEN GOTO 6
3 RETURN 0
5 RETURN 1
6 RETURN 2
End Function

Function send(t Uint64, a String, amt Uint64) Uint64
1 IF t == 1 THEN GOTO 5
2 IF t == 2 THEN GOTO 6
3 RETURN 0
5 RETURN SEND_DERO_TO_ADDRESS(a, amt)
6 RETURN SEND_ASSET_TO_ADDRESS(a, amt, LOAD("dReams"))
End Function

Function value(t Uint64) Uint64
1 IF t == 1 THEN GOTO 5
2 IF t == 2 THEN GOTO 6
3 RETURN 0
5 RETURN DEROVALUE()
6 RETURN ASSETVALUE(LOAD("dReams"))
End Function

Function storeBet(t Uint64, b Uint64, p Uint64, amt Uint64) Uint64
1 STORE("bets", LOAD("bets")+1)
2 STORE(b, p)
3 STORE("b_"+b, SIGNER())
4 STORE("b_"+b+"amt", amt)
5 STORE("b_"+b+"h", BLOCK_HEIGHT())
6 IF t < 2 THEN GOTO 20
9 STORE("b_"+b+"t", 1)
20 storeBal(t, amt, 0, 1)
30 RETURN 0
End Function

Function seven() Uint64
1 DIM i as Uint64
2 IF EXISTS("b_"+i+"h") == 0 THEN GOTO 10
3 IF BLOCK_HEIGHT() == LOAD("b_"+i+"h") THEN GOTO 10
4 DELETE(i)
5 DELETE("b_"+i)
6 DELETE("b_"+i+"amt")
7 DELETE("b_"+i+"h")
8 DELETE("b_"+i+"t")
10 LET i = i+1
11 IF i <= 30 THEN GOTO 2
20 STORE("bets", 0)
21 STORE("ot", 0)
22 STORE("otD", 0)
30 RETURN 0
End Function

Function findBal(t Uint64) Uint64
1 DIM b, p as Uint64
2 IF t == 1 THEN GOTO 5
3 IF t == 2 THEN GOTO 10
4 RETURN 0
5 LET b = LOAD("bal")
6 LET p = LOAD("ot")
7 GOTO 12
10 LET b = LOAD("balD")
11 LET p = LOAD("otD")
12 IF b > p THEN GOTO 20
15 RETURN 0
20 RETURN b-p
End Function

Function storeBal(t Uint64, amt Uint64, pay Uint64, p Uint64) Uint64
1 IF t == 1 THEN GOTO 4
2 IF t == 2 THEN GOTO 20
3 RETURN 0
4 STORE("bal", (LOAD("bal")+amt)-pay)
5 IF p == 0 THEN GOTO 7
6 STORE("ot", (LOAD("ot")+amt)-pay)
7 RETURN 0
20 STORE("balD", (LOAD("balD")+amt)-pay)
21 IF p == 0 THEN GOTO 30
22 STORE("otD", (LOAD("otD")+amt)-pay)
30 RETURN 0
End Function

Function placed(r Uint64) Uint64
1 DIM i, h, f, pay, t as Uint64
2 LET h = hitBox(r)
3 LET f = fieldBox(r)
4 IF EXISTS(i) == 0 THEN GOTO 80
5 IF BLOCK_HEIGHT() == LOAD("b_"+i+"h") THEN GOTO 80
6 IF EXISTS("b_"+i+"t") THEN GOTO 9
7 LET t = 1
8 GOTO 10
9 LET t = 2
10 IF LOAD(i) != 6 THEN GOTO 30
11 IF f THEN GOTO 60
12 DELETE(i)
13 DELETE("b_"+i)
14 IF t == 2 THEN GOTO 17
15 STORE("ot", LOAD("ot")-LOAD("b_"+i+"amt"))
16 GOTO 18
17 STORE("otD", LOAD("otD")-LOAD("b_"+i+"amt"))
18 DELETE("b_"+i+"amt")
19 DELETE("b_"+i+"h")
20 DELETE("b_"+i+"t")
21 STORE("bets", LOAD("bets")-1)
25 GOTO 80
30 IF LOAD(i) != h THEN GOTO 80
50 LET pay = payPlace(h, LOAD("b_"+i+"amt"))
55 GOTO 65
60 LET pay = (10000-LOAD("house"))*(LOAD("b_"+i+"amt")*f)/10000
65 IF findBal(t) < pay THEN GOTO 80
70 send(t, LOAD("b_"+i), pay)
75 storeBal(t, 0, pay, 0)
80 LET i = i+1
85 IF i <= 30 THEN GOTO 4
90 RETURN 0
End Function

Function storeRoll(b Uint64, d1 Uint64, d2 Uint64) Uint64 
1 DIM r as Uint64
2 LET r = d1+d2
3 STORE("rolls", LOAD("rolls")+1)
4 STORE("roll"+LOAD("rolls"), HEX(TXID())+"_"+b+"_"+d1+"_"+d2)
7 DELETE("roll"+(LOAD("rolls")-LOAD("display")))
8 IF EXISTS("out"+r) THEN GOTO 10
9 STORE("out"+r, 0)
10 STORE("out"+r, LOAD("out"+r)+1)
15 RETURN r
End Function

Function payRoll(b Uint64, amt Uint64) Uint64
1 IF b > 1 THEN GOTO 3
2 RETURN (10000-LOAD("house"))*(amt*2)/10000
3 IF b != 2 THEN GOTO 5
4 RETURN (10000-LOAD("house"))*(amt*5)/10000
5 IF b != 3 THEN GOTO 7
6 RETURN (10000-LOAD("house"))*(amt*8)/10000
7 IF b > 5 THEN GOTO 20
8 IF b < 6 THEN GOTO 10
9 RETURN 0
10 RETURN (10000-LOAD("house"))*(amt*16)/10000
20 RETURN (10000-LOAD("house"))*(amt*31)/10000
End Function

Function hitBox(r Uint64) Uint64
1 IF r == 4 THEN GOTO 10
2 IF r == 5 THEN GOTO 11
3 IF r == 6 THEN GOTO 12
4 IF r == 8 THEN GOTO 13
5 IF r == 9 THEN GOTO 14
6 IF r == 10 THEN GOTO 15
9 RETURN 999
10 RETURN 0
11 RETURN 1
12 RETURN 2
13 RETURN 3
14 RETURN 4
15 RETURN 5
End Function

Function payPlace(p Uint64, amt Uint64) Uint64
1 IF p == 0 || p == 5 THEN GOTO 10
2 IF p == 1 || p == 4 THEN GOTO 11
3 IF p == 2 || p == 3 THEN GOTO 12
5 RETURN 0 
10 RETURN (10000-LOAD("house"))*(9*amt/5)/10000
11 RETURN (10000-LOAD("house"))*(7*amt/5)/10000
12 RETURN (10000-LOAD("house"))*(7*amt/6)/10000
End Function

Function fieldBox(r Uint64) Uint64
1 IF r == 2 || r == 12 THEN GOTO 11
2 IF r == 3 || r == 4 || r == 9 || r == 10 || r == 11 THEN GOTO 10
5 RETURN 0
10 RETURN 1
11 RETURN 2
End Function

Function hit(b Uint64, r Uint64) Uint64
13 IF b == 0 THEN GOTO 43
14 IF b == 1 THEN GOTO 44
15 IF b == 2 THEN GOTO 45
16 IF b == 3 THEN GOTO 46
17 IF b == 4 THEN GOTO 47
18 IF b == 5 THEN GOTO 48
19 IF b == 6 THEN GOTO 49
20 IF b == 7 THEN GOTO 50
21 RETURN 0
43 IF r < 7 THEN GOTO 90 ELSE GOTO 80
44 IF r > 7 THEN GOTO 90 ELSE GOTO 80
45 IF r == 7 THEN GOTO 90 ELSE GOTO 80
46 IF r == 2 || r == 3 || r == 12 THEN GOTO 90 ELSE GOTO 80
47 IF r == 3 THEN GOTO 90 ELSE GOTO 80
48 IF r == 11 THEN GOTO 90 ELSE GOTO 80
49 IF r == 2 THEN GOTO 90 ELSE GOTO 80
50 IF r == 12 THEN GOTO 90
80 RETURN 0
90 RETURN 1
End Function

Function own() Uint64
1 IF SIGNER() == LOAD("owner") THEN GOTO 10
2 IF LOAD("sign") < 2 THEN GOTO 9
3 DIM i as Uint64
4 LET i = 1
5 IF EXISTS("si"+i) == 0 THEN GOTO 7
6 IF SIGNER() == LOAD("si"+i) THEN GOTO 10
7 LET i = i+1
8 IF i <= 9 THEN GOTO 5
9 RETURN 0
10 RETURN 1
End Function

Function AddS(new String) Uint64
1 IF IS_ADDRESS_VALID(ADDRESS_RAW(new)) == 0 THEN GOTO 3
2 IF SIGNER() == LOAD("owner") && LOAD("sign") < 10 THEN GOTO 4
3 RETURN 1
4 DIM i as Uint64
5 LET i = i+1
6 IF i > 9 THEN GOTO 3
7 IF EXISTS("si"+i) THEN GOTO 5
8 STORE("sign", LOAD("sign")+1)
9 STORE("si"+i, ADDRESS_RAW(new))
10 RETURN 0
End Function

Function RmvS(rm Uint64) Uint64 
1 IF SIGNER() == LOAD("owner") THEN GOTO 3
2 RETURN 1
3 IF EXISTS("si"+rm) == 0 THEN GOTO 6
4 STORE("sign", LOAD("sign")-1)
5 DELETE("si"+rm)
6 RETURN 0
End Function

Function Withdraw(t Uint64, amt Uint64) Uint64
1 IF findBal(t) < amt THEN GOTO 3
2 IF own() THEN GOTO 4
3 RETURN 1
4 send(t, SIGNER(), amt)
5 storeBal(t, 0, amt, 0)
6 RETURN 0
End Function

Function UpdateVar(min Uint64, max Uint64, live Uint64, h Uint64, d Uint64) Uint64
1 IF h > 1000 THEN GOTO 3
2 IF own() THEN GOTO 4
3 RETURN 1
4 STORE("min", min)
5 STORE("max", max)
6 STORE("live", live)
7 STORE("house", h)
8 STORE("display", d)
9 RETURN 0
End Function

Function UpdateCode(code String) Uint64
1 IF code == "" THEN GOTO 5
2 IF LOAD("ot") != 0 THEN GOTO 5
3 IF LOAD("otD") != 0 THEN GOTO 5
4 IF own() THEN GOTO 6
5 RETURN 1
6 UPDATE_SC_CODE(code)
7 STORE("v", LOAD("v")+1)
9 RETURN 0
End Function