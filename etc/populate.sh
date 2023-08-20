#!/usr/bin/bash

names=("Ana" "Maria" "João")
stack=("Go" "Python" "Java")

add() {
    curl -X POST -H "Content-Type: application/json" \
      -d "{\"apelido\" : \"$1\", \"nome\" : \"$1\", \"nascimento\" : \"2000-10-01\",\"stack\" : [\"Java\",\"$2\"]}" \
       http://localhost:9999/pessoas
}

add Fulano Go
add Beltrano JS
add Ciclano Python
add Batata Python