# Introduction

# Examples

See example scripts in *examples/*.

```
(import "string")

(defn take-second-column (lines)
    (map (fn (line) (item line 1)) lines))

(defn join-with-space (line)
    (join line " "))

(println ((. join-with-space take-second-column tail columnize @sh) "docker ps"))
```

# TODO:
 - [ ] keyword string syntax

```
(some-fn :keyword-argument-1 value-1 :keyword-argument-2 value-2 position-argument-3)

 - [x] variadic argument
```

(defn my-fn (arg-1 variadic...) (do-comething :with variadic))

 - [x] dictionary syntax

```
(set my-dict {
    :key-1 value-1
    :key-2 value-2
})
```

 - [.] error handling
 - [ ] syntax error handling
 - [x] comments
 - [x] list syntax

```
(set my-list [1 2 3 4])
```
