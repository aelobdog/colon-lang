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

    i (condition) :
        <statement>
    :i

OR (with else)

    i (condition) :
        <statement>
    :i e:
        <statement>
    :e

# loop

    l (condition) :
        <statement>
    :l

## function

    v: name = f (foo, bar, baz) :
    <statement>
    <statement>
    :f

## function call

    add(1 2 add(12 3))
    mul(12 + 2, add(12, 2) + 1)
