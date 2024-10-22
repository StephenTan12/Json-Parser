package main

type JsonObject map[string]interface{}
type JsonArray []interface{}
type Json interface{}

type Error struct {
	s string
}
