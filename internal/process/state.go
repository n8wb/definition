/*
	Copyright 2019 Whiteblock Inc.
	This file is a part of the Definition.

	Definition is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	Definition is distributed in the hope that it will be useful,
	but dock ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package process

import (
	"fmt"

	"github.com/whiteblock/definition/internal/entity"
	"github.com/whiteblock/definition/internal/maker"
	"github.com/whiteblock/definition/internal/merger"
	"github.com/whiteblock/definition/internal/parser"
	"github.com/whiteblock/definition/schema"
)

//System is for diff calculations
type System interface {
	GetAlreadyExists(state *entity.State, systems []schema.SystemComponent) (
		[]schema.SystemComponent, []schema.SystemComponent, bool)

	UpdateChanged(state *entity.State, spec schema.RootSchema,
		systems []schema.SystemComponent) (toAdd []entity.Service,
		toRemove []entity.Service, err error)

	//  Add modifies State
	Add(state *entity.State, spec schema.RootSchema, systems []schema.SystemComponent) ([]entity.Service, error)
	//Remove modifies state
	Remove(state *entity.State, spec schema.RootSchema, systems []string) ([]entity.Service, error)

	Tasks(state *entity.State, spec schema.RootSchema, tasks []schema.Task) ([]entity.Service, error)
}

type system struct {
	namer  parser.Names
	maker  maker.Service
	merger merger.System
}

func NewSystem(
	namer parser.Names,
	maker maker.Service,
	merger merger.System) System {
	return &system{namer: namer, maker: maker, merger: merger}
}

func (sys system) UpdateChanged(state *entity.State, spec schema.RootSchema,
	systems []schema.SystemComponent) (toAdd []entity.Service,
	toRemove []entity.Service, err error) {

	for _, systemUpdate := range systems {
		name := sys.namer.SystemComponent(systemUpdate)
		old, exists := state.SystemState[name]
		if !exists {
			return nil, nil, fmt.Errorf("system \"%s\" not found", name)
		}
		system := sys.merger.MergeLeft(systemUpdate, old)

		serv, err := sys.maker.FromSystemDiff(spec, old, system)
		if err != nil {
			return nil, nil, err
		}

		if system.Count < old.Count {
			toRemove = append(toRemove, serv...)
		} else {
			toAdd = append(toAdd, serv...)
		}
		state.SystemState[name] = system
	}
	return
}

func (sys system) GetAlreadyExists(state *entity.State, systems []schema.SystemComponent) (
	exist []schema.SystemComponent, noExist []schema.SystemComponent, anyExist bool) {

	anyExist = false
	for _, s := range systems {
		name := sys.namer.SystemComponent(s)
		_, exists := state.SystemState[name]
		if exists {
			anyExist = true
			exist = append(exist, s)
		} else {
			noExist = append(noExist, s)
		}
	}
	return
}

//  Add modifies State
func (sys system) Add(state *entity.State, spec schema.RootSchema,
	systems []schema.SystemComponent) ([]entity.Service, error) {
	out := []entity.Service{}

	for _, system := range systems {
		name := sys.namer.SystemComponent(system)
		_, exists := state.SystemState[name]
		if exists {
			return nil, fmt.Errorf("already have a system with the name \"%s\"", name)
		}
		services, err := sys.maker.FromSystem(spec, system)
		if err != nil {
			return nil, err
		}
		out = append(out, services...)
	}

	for _, system := range systems {
		name := sys.namer.SystemComponent(system)
		state.SystemState[name] = system
	}

	return out, nil
}

//Remove modifies state
func (sys system) Remove(state *entity.State, spec schema.RootSchema, systems []string) ([]entity.Service, error) {
	out := []entity.Service{}
	for _, toRemove := range systems {
		system, exists := state.SystemState[toRemove]
		if !exists {
			return nil, fmt.Errorf("system not found")
		}
		services, err := sys.maker.FromSystem(spec, system)
		if err != nil {
			return nil, err
		}
		out = append(out, services...)
	}
	for _, toRemove := range systems {
		delete(state.SystemState, toRemove)
	}
	return out, nil
}

func (sys system) Tasks(state *entity.State, spec schema.RootSchema, tasks []schema.Task) ([]entity.Service, error) {
	out := []entity.Service{}
	for i, task := range tasks {
		service, err := sys.maker.FromTask(spec, task, i)
		if err != nil {
			return nil, err
		}
		out = append(out, service)
	}
	return out, nil
}
