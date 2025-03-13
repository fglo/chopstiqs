package chopstiqs

import (
	"image"
	"testing"

	"github.com/fglo/chopstiqs/component"
	"github.com/fglo/chopstiqs/option"
)

func TestGUI_alignRootContainerInBounds(t *testing.T) {
	type fields struct {
		horizontalAglignment option.HorizontalAlignment
		verticalAlignment    option.VerticalAlignment
		rootContainer        *component.Container
	}
	type args struct {
		bounds image.Rectangle
	}

	rootContainer := component.NewContainer(&component.ContainerOptions{
		Width:  option.Int(10),
		Height: option.Int(10),
	})

	bounds := image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{100, 100},
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		wantX  float64
		wantY  float64
	}{
		{
			name: "HorizontalAlignment Default, VerticalAlignment Default",
			fields: fields{
				rootContainer: rootContainer,
			},
			args: args{
				bounds: bounds,
			},
			wantX: 0,
			wantY: 0,
		},
		{
			name: "HorizontalAlignment Left, VerticalAlignment Default",
			fields: fields{
				horizontalAglignment: option.AlignmentLeft,
				rootContainer:        rootContainer,
			},
			args: args{
				bounds: bounds,
			},
			wantX: 0,
			wantY: 0,
		},
		{
			name: "HorizontalAlignment Center, VerticalAlignment Default",
			fields: fields{
				horizontalAglignment: option.AlignmentCenteredHorizontally,
				rootContainer:        rootContainer,
			},
			args: args{
				bounds: bounds,
			},
			wantX: 45,
			wantY: 0,
		},
		{
			name: "HorizontalAlignment Right, VerticalAlignment Default",
			fields: fields{
				horizontalAglignment: option.AlignmentRight,
				rootContainer:        rootContainer,
			},
			args: args{
				bounds: bounds,
			},
			wantX: 90,
			wantY: 0,
		},
		{
			name: "HorizontalAlignment Default, VerticalAlignment Top",
			fields: fields{
				verticalAlignment: option.AlignmentTop,
				rootContainer:     rootContainer,
			},
			args: args{
				bounds: bounds,
			},
			wantX: 0,
			wantY: 0,
		},
		{
			name: "HorizontalAlignment Default, VerticalAlignment Center",
			fields: fields{
				verticalAlignment: option.AlignmentCenteredVertically,
				rootContainer:     rootContainer,
			},
			args: args{
				bounds: bounds,
			},
			wantX: 0,
			wantY: 45,
		},
		{
			name: "HorizontalAlignment Default, VerticalAlignment Bottom",
			fields: fields{
				verticalAlignment: option.AlignmentBottom,
				rootContainer:     rootContainer,
			},
			args: args{
				bounds: bounds,
			},
			wantX: 0,
			wantY: 90,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gui := NewGUI(&GUIOptions{
				HorizontalAlignment: tt.fields.horizontalAglignment,
				VerticalAlignment:   tt.fields.verticalAlignment,
			})

			gui.SetRootContainer(tt.fields.rootContainer)

			gui.alignRootContainerInBounds(tt.args.bounds)

			gotX := tt.fields.rootContainer.PosX()
			if gotX != tt.wantX {
				t.Errorf("got X %f, want %f", gotX, tt.wantX)
			}

			gotY := tt.fields.rootContainer.PosY()
			if gotY != tt.wantY {
				t.Errorf("got Y %f, want %f", gotY, tt.wantY)
			}
		})
	}
}
