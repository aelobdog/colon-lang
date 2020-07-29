# Colon Programming Language

## variable declaration

    v: var_name = value

## standard input

    n: var_name type

## standard output

    p: FOO + "message" + BAR
    p: "message"

## sleep / wait

    w: time units

## return

    r: value

## break

    b:

## continue

    c:

## if condition

    i:  (condition)
        <statement>
    :i

OR (with else)

    i:  (condition)
        <statement>
    :i e:
        <statement>
    :e

# loop

    l: (condition)
        <statement>
    :l

## function

    f: func_name (v: foo, v: bar)
        <statement>
    :f

OR

    v: name = f: (v: foobar)
    <statement>
    <statement>
    :f
