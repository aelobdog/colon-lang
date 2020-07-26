# Colon Programming Language

## variable declaration

    v: var_name = value

## standard input

    i: var_name type

## standard output

    p: "{FOO} message {BAR}"
    p: "message"

## sleep / wait

    w: time units

## return

    r: value

## break

    b:

## continue

    c:

## switch / if

    s:
    (condition)
        <statement>
        <statement>
    :s

# loop

    l: (condition)
        <statement>
        <statement>
    :l

## function

    f: func_name (v: foo, v: bar)
        <statement>
        <statement>
    :f

OR

    v: name = f: (v: foobar)
    <statement>
    <statement>
    :f

## Data / Literal types

- Integer [-1, 0, 1, ...]
- Floating [-1.3, 0.0, 12.6, ...]
- Boolean [t, T, true, TRUE, f, F, false, FALSE]
- String ["hello", ...]
