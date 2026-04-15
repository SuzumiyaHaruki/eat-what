package main

import (
	"errors"
	"math/rand"
	"strings"
)

type FoodManager struct {
	options []string
	seen    map[string]struct{}
}

func NewFoodManager() *FoodManager {
	return &FoodManager{
		options: make([]string, 0),
		seen:    make(map[string]struct{}),
	}
}

func normalizeOption(s string) string {
	return strings.TrimSpace(s)
}

func (m *FoodManager) AddOption(s string) bool {
	s = normalizeOption(s)
	if s == "" {
		return false
	}
	if _, exists := m.seen[s]; exists {
		return false
	}
	m.options = append(m.options, s)
	m.seen[s] = struct{}{}
	return true
}

func (m *FoodManager) AddOptions(lines []string) int {
	count := 0
	for _, line := range lines {
		if m.AddOption(line) {
			count++
		}
	}
	return count
}

func (m *FoodManager) RemoveAt(index int) bool {
	if index < 0 || index >= len(m.options) {
		return false
	}
	item := m.options[index]
	delete(m.seen, item)
	m.options = append(m.options[:index], m.options[index+1:]...)
	return true
}

func (m *FoodManager) Clear() {
	m.options = make([]string, 0)
	m.seen = make(map[string]struct{})
}

func (m *FoodManager) ReplaceOptions(lines []string) int {
	m.Clear()
	return m.AddOptions(lines)
}

func (m *FoodManager) Count() int {
	return len(m.options)
}

func (m *FoodManager) Get(index int) string {
	if index < 0 || index >= len(m.options) {
		return ""
	}
	return m.options[index]
}

func (m *FoodManager) All() []string {
	out := make([]string, len(m.options))
	copy(out, m.options)
	return out
}

func (m *FoodManager) PickRandom() (string, error) {
	if len(m.options) == 0 {
		return "", errors.New("没有可供抽取的选项")
	}
	idx := rand.Intn(len(m.options))
	return m.options[idx], nil
}
