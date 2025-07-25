package model

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Packet struct {
	Name       string            `json:"name"`
	Version    string            `json:"ver"`
	Packets    []Dependency      `json:"packets,omitempty"`
	Targets    []Target          `json:"-"`
	RawTargets []json.RawMessage `json:"targets"`
}

type Target struct {
	Path    string `json:"path"`
	Exclude string `json:"exclude,omitempty"`
}

type Dependency struct {
	Name    string `json:"name"`
	Version string `json:"ver,omitempty"` // необязательное поле
}

func (p *Packet) UnmarshalJSON(data []byte) error {
	// временная структура, чтобы распарсить всё, кроме targets
	type Alias Packet
	aux := &struct {
		*Alias
		Targets []json.RawMessage `json:"targets"`
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	// обрабатываем каждый элемент targets
	for _, raw := range aux.Targets {
		// пробуем распарсить как строку
		var s string
		if err := json.Unmarshal(raw, &s); err == nil {
			// строка вида "path,exclude"
			parts := strings.SplitN(s, ",", 2)
			t := Target{Path: strings.TrimSpace(parts[0])}
			if len(parts) > 1 {
				t.Exclude = strings.TrimSpace(parts[1])
			}
			p.Targets = append(p.Targets, t)
			continue
		}

		// иначе — это объект
		var t Target
		if err := json.Unmarshal(raw, &t); err != nil {
			return fmt.Errorf("invalid target: %w", err)
		}
		p.Targets = append(p.Targets, t)
	}

	return nil
}
