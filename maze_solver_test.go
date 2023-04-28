package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_mazeSolverImpl_FindSolutions(t *testing.T) {
	type args struct {
		maze *Maze
	}
	tests := []struct {
		name    string
		m       *mazeSolverImpl
		args    args
		wantErr bool
		verify  func(t *testing.T, got []MazeSolution)
	}{
		{
			name: "Find one solution",
			m:    &mazeSolverImpl{},
			args: args{
				maze: &Maze{
					EntranceX:  0,
					EntranceY:  0,
					GridWidth:  8,
					GridHeight: 8,
					Walls: []byte{
						68, 85, 20, 118, 18, 218, 74, 2, 0,
					},
				},
			},
			wantErr: false,
			verify: func(t *testing.T, got []MazeSolution) {
				want := []MazeSolution{
					{
						Length: 10,
						Path: []string{
							"A1", "B1", "B2", "B3", "A3", "A4", "A5", "A6", "A7", "A8",
						},
						Exit: "A8",
					},
				}
				if !reflect.DeepEqual(got, want) {
					t.Errorf("mazeSolverImpl.FindSolutions() = %v, want %v", got, want)
				}
			},
		},
		{
			name: "Find two solutions",
			m:    &mazeSolverImpl{},
			args: args{
				maze: &Maze{
					EntranceX:  0,
					EntranceY:  0,
					GridWidth:  8,
					GridHeight: 8,
					Walls: []byte{
						68, 85, 20, 118, 18, 218, 64, 95, 0,
					},
				},
			},
			wantErr: false,
			verify: func(t *testing.T, got []MazeSolution) {
				want := []MazeSolution{
					{
						Length: 15,
						Path: []string{
							"A1", "B1", "B2", "B3", "A3", "A4", "A5", "A6", "A7", "B7", "C7", "D7", "E7", "F7", "F8",
						},
						Exit: "F8",
					},
					{
						Length: 31,
						Path: []string{
							"A1", "B1", "B2", "B3", "A3", "A4", "A5", "A6", "A7", "B7", "C7", "C6", "C5", "D5", "D4", "D3", "D2", "D1", "E1", "F1", "F2", "F3", "G3", "H3", "H4", "H5", "G5", "F5", "F6", "F7", "F8",
						},
						Exit: "F8",
					},
				}
				if !reflect.DeepEqual(got, want) {
					t.Errorf("mazeSolverImpl.FindSolutions() = %v, want %v", got, want)
				}
			},
		},
		{
			name: "Find three solutions with two exits",
			m:    &mazeSolverImpl{},
			args: args{
				maze: &Maze{
					EntranceX:  0,
					EntranceY:  0,
					GridWidth:  8,
					GridHeight: 8,
					Walls: []byte{
						68, 85, 20, 118, 18, 90, 64, 95, 0,
					},
				},
			},
			wantErr: false,
			verify: func(t *testing.T, got []MazeSolution) {
				assert.Equal(t, 4, len(got))
				assert.Equal(t, "F8", got[0].Exit)
				assert.Equal(t, "H8", got[1].Exit)
				assert.Equal(t, "H8", got[2].Exit)
				assert.Equal(t, "F8", got[3].Exit)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mazeSolverImpl{}
			got, err := m.FindSolutions(tt.args.maze)
			if (err != nil) != tt.wantErr {
				t.Errorf("mazeSolverImpl.FindSolutions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.verify(t, got)
		})
	}
}
