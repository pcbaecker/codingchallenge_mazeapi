package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaze_GetByteAddress(t *testing.T) {
	type fields struct {
		Id         string
		EntraceX   uint16
		EntraceY   uint16
		GridWidth  uint16
		GridHeight uint16
		Walls      []byte
	}
	type args struct {
		x uint16
		y uint16
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint32
	}{
		{
			name: "First byte first bit",
			fields: fields{
				GridWidth:  16,
				GridHeight: 16,
				Walls:      make([]byte, 16*16/8+1),
			},
			args: args{
				x: 0,
				y: 0,
			},
			want: 0,
		},
		{
			name: "First byte last bit",
			fields: fields{
				GridWidth:  16,
				GridHeight: 16,
				Walls:      make([]byte, 16*16/8+1),
			},
			args: args{
				x: 7,
				y: 0,
			},
			want: 0,
		},
		{
			name: "Second byte first bit",
			fields: fields{
				GridWidth:  16,
				GridHeight: 16,
				Walls:      make([]byte, 16*16/8+1),
			},
			args: args{
				x: 8,
				y: 0,
			},
			want: 1,
		},
		{
			name: "Last byte fist bit",
			fields: fields{
				GridWidth:  16,
				GridHeight: 16,
				Walls:      make([]byte, 16*16/8+1),
			},
			args: args{
				x: 15,
				y: 15,
			},
			want: 31,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Maze{
				Id:         tt.fields.Id,
				EntranceX:  tt.fields.EntraceX,
				EntranceY:  tt.fields.EntraceY,
				GridWidth:  tt.fields.GridWidth,
				GridHeight: tt.fields.GridHeight,
				Walls:      tt.fields.Walls,
			}
			if got := m.GetByteAddress(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("Maze.GetByteAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaze_InitWalls(t *testing.T) {
	type args struct {
		x uint16
		y uint16
	}
	tests := []struct {
		name   string
		args   args
		verify func(maze *Maze)
	}{
		{
			name: "3x3",
			args: args{
				x: 3,
				y: 3,
			},
			verify: func(maze *Maze) {
				if len(maze.Walls) != 2 {
					t.Errorf("len(Maze.Walls) = %v, want %v", len(maze.Walls), 2)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			maze := &Maze{}
			maze.InitWalls(tt.args.x, tt.args.y)
			tt.verify(maze)
		})
	}
}

func TestMaze_GetBitAddress(t *testing.T) {
	type fields struct {
		Id         string
		EntraceX   uint16
		EntraceY   uint16
		GridWidth  uint16
		GridHeight uint16
		Walls      []byte
	}
	type args struct {
		x uint16
		y uint16
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint8
	}{
		{
			name: "First byte first bit",
			fields: fields{
				GridWidth:  16,
				GridHeight: 16,
			},
			args: args{
				x: 0,
				y: 0,
			},
			want: 0,
		},
		{
			name: "First byte last bit",
			fields: fields{
				GridWidth:  16,
				GridHeight: 16,
			},
			args: args{
				x: 7,
				y: 0,
			},
			want: 7,
		},
		{
			name: "Last byte first bit",
			fields: fields{
				GridWidth:  16,
				GridHeight: 16,
			},
			args: args{
				x: 8,
				y: 15,
			},
			want: 0,
		},
		{
			name: "First byte last bit",
			fields: fields{
				GridWidth:  16,
				GridHeight: 16,
			},
			args: args{
				x: 15,
				y: 15,
			},
			want: 7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Maze{
				Id:         tt.fields.Id,
				EntranceX:  tt.fields.EntraceX,
				EntranceY:  tt.fields.EntraceY,
				GridWidth:  tt.fields.GridWidth,
				GridHeight: tt.fields.GridHeight,
				Walls:      tt.fields.Walls,
			}
			if got := m.GetBitAddress(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("Maze.GetBitAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readPosition(t *testing.T) {
	type args struct {
		position string
	}
	tests := []struct {
		name    string
		args    args
		want    uint16
		want1   uint16
		wantErr bool
	}{
		{
			name: "A1",
			args: args{
				position: "A1",
			},
			want:    0,
			want1:   0,
			wantErr: false,
		},
		{
			name: "B2",
			args: args{
				position: "B2",
			},
			want:    1,
			want1:   1,
			wantErr: false,
		},
		{
			name: "AA10",
			args: args{
				position: "AA10",
			},
			want:    26,
			want1:   9,
			wantErr: false,
		},
		{
			name: "invalid",
			args: args{
				position: "A",
			},
			want:    0,
			want1:   0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := readPosition(tt.args.position)
			if (err != nil) != tt.wantErr {
				t.Errorf("readPosition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("readPosition() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("readPosition() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_readGridSize(t *testing.T) {
	type args struct {
		gridsize string
	}
	tests := []struct {
		name    string
		args    args
		want    uint16
		want1   uint16
		wantErr bool
	}{
		{
			name: "16x16",
			args: args{
				gridsize: "16x16",
			},
			want:    16,
			want1:   16,
			wantErr: false,
		},
		{
			name: "invalid",
			args: args{
				gridsize: "16x",
			},
			want:    0,
			want1:   0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := readGridSize(tt.args.gridsize)
			if (err != nil) != tt.wantErr {
				t.Errorf("readGridSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("readGridSize() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("readGridSize() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_stringsToPath(t *testing.T) {
	type args struct {
		path []string
	}
	tests := []struct {
		name   string
		args   args
		verify func(*testing.T, *PathItem)
	}{
		{
			name: "Empty path",
			args: args{
				path: []string{},
			},
			verify: func(t *testing.T, got *PathItem) {
				assert.Nil(t, got)
			},
		},
		{
			name: "Valid path",
			args: args{
				path: []string{"A1", "A2", "A3"},
			},
			verify: func(t *testing.T, got *PathItem) {
				assert.Equal(t, uint16(0), got.X)
				assert.Equal(t, uint16(2), got.Y)
				got = got.Prev
				assert.Equal(t, uint16(0), got.X)
				assert.Equal(t, uint16(1), got.Y)
				got = got.Prev
				assert.Equal(t, uint16(0), got.X)
				assert.Equal(t, uint16(0), got.Y)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := stringsToPath(tt.args.path)
			if err != nil {
				t.Errorf("stringsToPath() error = %v", err)
				return
			}
			tt.verify(t, got)
		})
	}
}
