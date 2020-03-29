package main

type SessionUUID string

type sessions map[SessionUUID]*state.S
