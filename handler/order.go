package handler

import (
	"fmt"
	"net/http"
)

type Order struct{}

func (o *Order) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Criar ordens")
}

func (o *Order) List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Listar ordens")
}

func (o *Order) GetById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Buscar ordem por id")
}

func (o *Order) UpdateById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Atualizar ordem por id")
}

func (o *Order) DeleteById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Deletar ordem por id")
}
